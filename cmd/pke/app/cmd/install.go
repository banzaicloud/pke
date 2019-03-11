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
