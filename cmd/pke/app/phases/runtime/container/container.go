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
