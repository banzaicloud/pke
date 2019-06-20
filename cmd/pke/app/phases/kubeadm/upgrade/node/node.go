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

package node

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
	use   = "kubernetes-node"
	short = "Kubernetes worker node upgrade"

	kubeConfig = "/etc/kubernetes/kubelet.conf"
	cmdKubeadm = "/bin/kubeadm"
)

var _ phases.Runnable = (*Node)(nil)

type Node struct {
	kubernetesVersion string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Node{})
}

func (*Node) Use() string {
	return use
}

func (*Node) Short() string {
	return short
}

func (*Node) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.14.3", "Kubernetes version")
}

func (n *Node) Validate(cmd *cobra.Command) error {
	var err error

	n.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: n.kubernetesVersion,
	}); err != nil {
		return err
	}

	return nil
}

func (n *Node) Run(out io.Writer) error {
	return upgrade.RunWithSkewCheck(out, use, n.kubernetesVersion, kubeConfig, n.upgradeMinor, n.upgradePatch)
}

func (n *Node) upgradeMinor(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] upgrading node from %s to %s\n", use, from, to)

	return n.upgradePatch(out, from, to)
}

func (n *Node) upgradePatch(out io.Writer, from, to *semver.Version) error {
	_, _ = fmt.Fprintf(out, "[%s] patching node from %s to %s\n", use, from, to)

	return n.upgrade(out, from, to)
}

func (n *Node) upgrade(out io.Writer, from, to *semver.Version) error {
	var pm linux.KubernetesPackages
	pm = linux.NewYumInstaller()
	err := pm.InstallKubeadmPackage(out, to.String())
	if err != nil {
		return errors.Wrapf(err, "failed to upgrade kubeadm to version %s", to)
	}

	args := []string{
		"upgrade",
		"node",
		"config",
		"--kubelet-version",
		to.String(),
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
