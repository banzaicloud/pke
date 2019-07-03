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

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "container-runtime"
	short = "Container runtime installation"
)

var _ phases.Runnable = (*Runtime)(nil)

type Runtime struct {
	imageRepository string
}

func NewCommand() *cobra.Command {
	return phases.NewCommand(&Runtime{})
}

func (r *Runtime) Use() string {
	return use
}

func (r *Runtime) Short() string {
	return short
}

func (r *Runtime) RegisterFlags(flags *pflag.FlagSet) {
	// Image repository
	flags.String(constants.FlagImageRepository, "banzaicloud", "Prefix for image repository")
}

func (r *Runtime) Validate(cmd *cobra.Command) error {
	var err error
	r.imageRepository, err = cmd.Flags().GetString(constants.FlagImageRepository)
	if err != nil {
		return err
	}

	return validator.NotEmpty(map[string]interface{}{
		constants.FlagImageRepository: r.imageRepository,
	})
}

func (r *Runtime) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] running\n", r.Use())

	return r.installRuntime(out)
}
