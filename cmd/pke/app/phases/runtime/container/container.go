package container

import (
	"fmt"
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "container-runtime"
	short = "Container runtime installation"
)

var _ phases.Runnable = (*Runtime)(nil)

type Runtime struct{}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Runtime{})
}

func (r *Runtime) Use() string {
	return use
}

func (r *Runtime) Short() string {
	return short
}

func (r *Runtime) RegisterFlags(flags *pflag.FlagSet) {}

func (r *Runtime) Validate(cmd *cobra.Command) error {
	return nil
}

func (r *Runtime) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", r.Use())

	return installRuntime(out)
}
