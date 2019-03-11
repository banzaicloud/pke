package cmd

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/token/create"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/token/list"
	"github.com/spf13/cobra"
)

// NewCmdVersion provides the version information of banzai-cloud-pke.
func NewCmdToken(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage Kubernetes bootstrap tokens",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(create.NewCommand(out))
	cmd.AddCommand(list.NewCommand(out))

	return cmd
}
