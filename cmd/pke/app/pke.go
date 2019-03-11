package app

import (
	"flag"
	"os"

	"github.com/banzaicloud/pke/cmd/pke/app/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Run(gitVersion, gitCommit, gitTreeState, buildDate string) error {
	cobra.EnableCommandSorting = false
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	c := cmd.NewPKECommand(os.Stdin, os.Stdout, gitVersion, gitCommit, gitTreeState, buildDate)
	return c.Execute()
}
