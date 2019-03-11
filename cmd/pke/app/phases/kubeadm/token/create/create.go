package create

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "create"
	short = "Create Kubernetes bootstrap token"
)

var _ phases.Runnable = (*Create)(nil)

type Create struct {
	kubernetesVersion string
	imageRepository   string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Create{})
}

func (*Create) Use() string {
	return use
}

func (*Create) Short() string {
	return short
}

func (*Create) RegisterFlags(flags *pflag.FlagSet) {}

func (*Create) Validate(cmd *cobra.Command) error {
	return nil
}

func (*Create) Run(out io.Writer) error {
	return nil
}
