// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane installation"

	cmdKubeadm = "/bin/kubeadm"
	cmdKubectl = "/bin/kubectl"
	kubeConfig = "/etc/kubernetes/admin.conf"

	MaximumAllowedMinorVersionUpgradeSkew = 1
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	kubernetesVersion                string
	kubernetesAdditionalControlPlane bool
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &ControlPlane{})
}

func (*ControlPlane) Use() string {
	return use
}

func (*ControlPlane) Short() string {
	return short
}

func (*ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.14.0", "Kubernetes version")
	// Additional Control Plane
	flags.Bool(constants.FlagAdditionalControlPlane, false, "Treat node as additional control plane")
}

func (c *ControlPlane) Validate(cmd *cobra.Command) error {
	var err error

	c.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: c.kubernetesVersion,
	}); err != nil {
		return err
	}

	c.kubernetesAdditionalControlPlane, err = cmd.Flags().GetBool(constants.FlagAdditionalControlPlane)

	return err
}

func (c *ControlPlane) Run(out io.Writer) error {
	// query current kubernetes version
	clientVersion, serverVersion, err := kubectlVersion(out)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "[%s] client version: %s\n", use, clientVersion)
	_, _ = fmt.Fprintf(out, "[%s] server version: %s\n", use, serverVersion)

	cmVer, err := configmapVersion(out)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "[%s] configmap version: %s\n", use, cmVer)

	// minor or patch update
	srvVer, err := semver.NewVersion(serverVersion)
	if err != nil {
		return err
	}
	ver, err := semver.NewVersion(c.kubernetesVersion)
	if err != nil {
		return err
	}
	if srvVer.Major() != ver.Major() {
		return errors.New(fmt.Sprintf(
			"major version upgrade not supported. trying to upgrade from %d.x to %d.x",
			srvVer.Major(),
			ver.Major(),
		))
	}

	if srvVer.Minor() != ver.Minor() {
		if srvVer.Minor() > ver.Minor() {
			return errors.New(fmt.Sprintf(
				"downgrade not supported. trying to upgrade from %s to %s",
				srvVer,
				ver,
			))
		}
		if srvVer.Minor()+MaximumAllowedMinorVersionUpgradeSkew < ver.Minor() {
			return errors.New(fmt.Sprintf(
				"only %d minor version can be updated at a time. trying to upgrade from %d.%d to %d.%d",
				MaximumAllowedMinorVersionUpgradeSkew,
				srvVer.Major(),
				srvVer.Minor(),
				ver.Major(),
				ver.Minor(),
			))
		}
		// Minor version bump
		return c.upgradeMinor(out, srvVer, ver)
	}

	if srvVer.Patch() > ver.Patch() {
		return errors.New(fmt.Sprintf(
			"downgrade not supported. trying to upgrade from %s to %s",
			srvVer,
			ver,
		))
	}

	return c.upgradePatch(out, srvVer, ver)
}

func (c *ControlPlane) upgradeMinor(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] upgrading node from %s to %s\n", use, from, to)

	return c.upgradePatch(out, from, to)
}

func (c *ControlPlane) upgradePatch(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] patching node from %s to %s\n", use, from, to)

	return c.upgrade(out, from, to)
}

func (c *ControlPlane) upgrade(out io.Writer, from, to *semver.Version) error {
	var pm linux.KubernetesPackages
	pm = linux.NewYumInstaller()
	err := pm.InstallKubeadmPackage(out, to.String())
	if err != nil {
		return errors.Wrapf(err, "failed to upgrade kubeadm to version %s", to)
	}

	var args []string
	if c.kubernetesAdditionalControlPlane {
		args = []string{
			"upgrade",
			"node",
			"experimental-control-plane",
			to.String(),
		}
	} else {
		args = []string{
			"upgrade",
			"apply",
			"-f",
			to.String(),
		}
	}
	err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = pm.InstallKubernetesPackages(out, to.String())
	if err != nil {
		return err
	}

	err = linux.SystemctlReload(out)
	if err != nil {
		return err
	}

	err = linux.SystemctlStop(out, "kubelet")
	if err != nil {
		return err
	}

	err = linux.SystemctlStart(out, "kubelet")
	if err != nil {
		return err
	}

	return nil
}

func kubectlVersion(out io.Writer) (clientVersion, serverVersion string, err error) {
	cmd := runner.Cmd(ioutil.Discard, cmdKubectl, "version", "-o", "json")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	ver, err := cmd.Output()
	if err != nil {
		return
	}
	serverVer := struct {
		ClientVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"clientVersion"`
		ServerVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"serverVersion"`
	}{}
	err = json.Unmarshal(ver, &serverVer)
	if err != nil {
		return
	}

	clientVersion = serverVer.ClientVersion.GitVersion
	serverVersion = serverVer.ServerVersion.GitVersion
	return
}

func configmapVersion(out io.Writer) (version string, err error) {
	o, err := runner.Cmd(ioutil.Discard, cmdKubectl, "-n", "kube-system", "get", "cm", "kubeadm-config", "-ojsonpath={.data.ClusterConfiguration}").Output()
	if err != nil {
		return
	}

	cmVer := struct {
		KubernetesVersion string `yaml:"kubernetesVersion"`
	}{}
	err = yaml.Unmarshal(o, &cmVer)
	if err != nil {
		return
	}

	return cmVer.KubernetesVersion, nil
}
