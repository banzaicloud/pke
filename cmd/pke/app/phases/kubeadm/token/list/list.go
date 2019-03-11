package list

import (
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

const (
	use   = "list"
	short = "List Kubernetes bootstrap token(s)"
)

var _ phases.Runnable = (*List)(nil)

type List struct {
	kubernetesVersion string
	imageRepository   string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &List{})
}

func (*List) Use() string {
	return use
}

func (*List) Short() string {
	return short
}

func (*List) RegisterFlags(flags *pflag.FlagSet) {}

func (*List) Validate(cmd *cobra.Command) error {
	return nil
}

func (*List) Run(out io.Writer) error {
	return nil
}
