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

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/upgrade"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane upgrade"

	kubeConfig = "/etc/kubernetes/admin.conf"
	cmdKubeadm = "/bin/kubeadm"
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
	flags.String(constants.FlagKubernetesVersion, "1.14.3", "Kubernetes version")
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
