package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewPKECommand(in io.Reader, out io.Writer, gitVersion, gitCommit, gitTreeState, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "pke",
		Short:             "Bootstrap a secure Kubernetes cluster with Banzai Cloud Pipeline Kubernetes Engine (PKE)",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.ResetFlags()

	cmd.AddCommand(NewCmdInstall(out))
	cmd.AddCommand(NewCmdImage(out))
	cmd.AddCommand(NewCmdToken(out))
	cmd.AddCommand(NewCmdVersion(out, gitVersion, gitCommit, gitTreeState, buildDate))

	return cmd
}
