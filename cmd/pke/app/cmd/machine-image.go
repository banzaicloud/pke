package cmd

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/images"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/version"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/runtime/container"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/runtime/kubernetes"
	"github.com/spf13/cobra"
)

// NewCmdImage .
func NewCmdImage(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machine-image",
		Short: "Machine image build helper for Banzai Cloud Pipeline Kubernetes Engine (PKE)",
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(version.NewCommand(out))
	cmd.AddCommand(container.NewCommand(out))
	cmd.AddCommand(kubernetes.NewCommand(out))
	cmd.AddCommand(images.NewCommand(out))

	phases.MakeRunnable(cmd)

	return cmd
}
