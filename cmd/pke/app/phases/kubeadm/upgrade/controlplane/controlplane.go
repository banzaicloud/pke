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
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/upgrade"
	"github.com/banzaicloud/pke/cmd/pke/app/util/flags"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane upgrade"

	kubeConfig    = "/etc/kubernetes/admin.conf"
	cmdKubeadm    = "kubeadm"
	kubeadmConfig = "/etc/kubernetes/kubeadm.conf"
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	kubernetesVersion                string
	kubernetesAdditionalControlPlane bool
}

func NewCommand() *cobra.Command {
	return phases.NewCommand(&ControlPlane{})
}

func (*ControlPlane) Use() string {
	return use
}

func (*ControlPlane) Short() string {
	return short
}

func (*ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.17.0", "Kubernetes version")
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

	flags.PrintFlags(cmd.OutOrStdout(), c.Use(), cmd.Flags())

	return err
}

func (c *ControlPlane) Run(out io.Writer) error {
	return upgrade.RunWithSkewCheck(out, use, c.kubernetesVersion, kubeConfig, c.upgradeMinor, c.upgradePatch)
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
	pm, err := linux.KubernetesPackagesImpl(out)
	if err != nil {
		return err
	}
	err = pm.InstallKubeadmPackage(out, to.String())
	if err != nil {
		return errors.Wrapf(err, "failed to upgrade kubeadm to version %s", to)
	}

	var args []string
	if c.kubernetesAdditionalControlPlane {
		args = []string{
			"upgrade",
			"node",
		}
		c, _ := semver.NewConstraint("<1.15")
		if c.Check(to) {
			args = append(args, "experimental-control-plane")
		} else {
			args = append(args, "--kubelet-version")
		}

	} else {
		// TODO migrate here
		err := c.migrate(out, from, to)
		if err != nil {
			// TODO revert to the previous version of kubeadm if migration failed
			return errors.Wrapf(err, "failed to migrate kubeadm to version %s", to)
		}
		// TODO generate new kubeadm config
		// TODO missig drain node

		args = []string{
			"upgrade",
			"apply",
			"-f",
		}
		c, _ := semver.NewConstraint("1.16.x")
		if c.Check(to) {
			args = append(args, "--ignore-preflight-errors=CoreDNSUnsupportedPlugins")
		}
	}
	// target version
	args = append(args, to.String())

	_, err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
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

func (c *ControlPlane) migrate(out io.Writer, from, to *semver.Version) error {
	oldKubeadmConfig := "/etc/kubernetes/kubeadm.conf-" + from.String() + time.Now().Format(time.RFC3339)
	input, err := ioutil.ReadFile(kubeadmConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(oldKubeadmConfig, input, 0644)
	if err != nil {
		return err
	}

	args := []string{
		"config",
		"migrate",
		"--old-config",
		oldKubeadmConfig,
		"--new-config",
		kubeadmConfig,
	}

	_, err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = updateKubeadmConfig(from, to)
	if err != nil {
		return err
	}

	err = uploadKubeadmConf(out)
	if err != nil {
		return err
	}

	return nil
}

func updateKubeadmConfig(from, to *semver.Version) error {
	input, err := ioutil.ReadFile(kubeadmConfig)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	c, _ := semver.NewConstraint(">1.17")

	for i, line := range lines {
		if strings.Contains(line, "kubernetesVersion: v"+from.String()) {
			lines[i] = "kubernetesVersion: v" + to.String()
		}
		if strings.Contains(line, "useHyperKubeImage:") && c.Check(to) {
			lines[i] = ""
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(kubeadmConfig, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

func uploadKubeadmConf(out io.Writer) error {
	args := []string{
		"init",
		"phase",
		"upload-config",
		"all",
		"--config",
		kubeadmConfig,
	}

	_, err := runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}
	return nil
}
