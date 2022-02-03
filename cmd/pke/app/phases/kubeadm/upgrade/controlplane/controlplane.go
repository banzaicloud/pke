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
	"os"
	"time"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
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

	kubeConfig            = "/etc/kubernetes/admin.conf"
	cmdKubeadm            = "kubeadm"
	kubeadmConfig         = "/etc/kubernetes/kubeadm.conf"
	kubeadmMigratedConfig = "/etc/kubernetes/kubeadm-migrated.conf"
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	config config.Config

	kubernetesVersion                string
	kubernetesAdditionalControlPlane bool
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&ControlPlane{config: config})
}

func (*ControlPlane) Use() string {
	return use
}

func (*ControlPlane) Short() string {
	return short
}

func (c *ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, c.config.Kubernetes.Version, "Kubernetes version")
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

		c, _ := semver.NewConstraint("<1.20")
		if c.Check(to) { // target version
			args = append(args, "--kubelet-version", to.String())
		}

	} else {
		fromVersion, _ := semver.NewConstraint("1.21.x")
		toVersion, _ := semver.NewConstraint("1.22.x")
		if fromVersion.Check(from) && toVersion.Check(to) {
			// migrate kubeadm config to v1beta3
			err = c.migrateKubeadmConfig(out, from, to)
			if err != nil {
				return err
			}
		}

		args = []string{
			"upgrade",
			"apply",
			"-f",
		}
		// target version
		args = append(args, to.String())
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

func (c *ControlPlane) uploadKubeadmConf(out io.Writer) error {
	args := []string{
		"init",
		"phase",
		"upload-config",
		"kubeadm",
		"--config",
		kubeadmConfig,
	}

	_, err := runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}

func (c *ControlPlane) migrateKubeadmConfig(out io.Writer, from, to *semver.Version) error {

	args := []string{
		"config",
		"migrate",
		"--old-config",
		kubeadmConfig,
		"--new-config",
		kubeadmMigratedConfig,
	}

	_, err := runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	if err := renameKubeadmConfigs(out); err != nil {
		return err
	}

	return c.uploadKubeadmConf(out)
}

func renameKubeadmConfigs(out io.Writer) error {
	timestampedOldConfig := "/etc/kubernetes/kubeadm-" + time.Now().Format(time.RFC3339) + ".conf"

	if err := os.Rename(kubeadmConfig, timestampedOldConfig); err != nil {
		return err
	}
	if err := os.Rename(kubeadmMigratedConfig, kubeadmConfig); err != nil {
		return err
	}

	return nil
}
