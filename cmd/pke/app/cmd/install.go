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

package cmd

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/controlplane"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/node"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/version"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/pipeline/certificates"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/pipeline/ready"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/runtime/container"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/runtime/kubernetes"
	"github.com/spf13/cobra"
)

// NewCmdVersion provides the version information of banzai-cloud-pke.
func NewCmdInstall(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a single Banzai Cloud Pipeline Kubernetes Engine (PKE) machine",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(single(out))
	cmd.AddCommand(master(out))
	cmd.AddCommand(worker(out))

	return cmd
}

func master(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Master node",
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(version.NewCommand(out))
	cmd.AddCommand(container.NewCommand(out))
	cmd.AddCommand(kubernetes.NewCommand(out))
	cmd.AddCommand(certificates.NewCommand(out))
	cmd.AddCommand(controlplane.NewCommand(out))
	cmd.AddCommand(ready.NewCommand(out, ready.RoleMaster))

	phases.MakeRunnable(cmd)

	return cmd
}

func single(out io.Writer) *cobra.Command {
	m := master(out)
	m.Use = "single"
	m.Short = "Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) on a single machine"
	_ = m.Flags().Set(constants.FlagClusterMode, "single")

	return m
}

func worker(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Worker node",
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(version.NewCommand(out))
	cmd.AddCommand(container.NewCommand(out))
	cmd.AddCommand(kubernetes.NewCommand(out))
	cmd.AddCommand(node.NewCommand(out))
	cmd.AddCommand(ready.NewCommand(out, ready.RoleWorker))

	phases.MakeRunnable(cmd)

	return cmd
}
