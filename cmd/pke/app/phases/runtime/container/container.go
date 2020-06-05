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

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
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
	config config.Config

	containerRuntime string
	imageRepository  string
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&Runtime{config: config})
}

func (r *Runtime) Use() string {
	return use
}

func (r *Runtime) Short() string {
	return short
}

func (r *Runtime) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes container runtime
	flags.String(constants.FlagContainerRuntime, r.config.ContainerRuntime.Type, "Kubernetes container runtime")

	// Image repository
	flags.String(constants.FlagImageRepository, "banzaicloud", "Prefix for image repository")
}

func (r *Runtime) Validate(cmd *cobra.Command) (err error) {
	r.containerRuntime, err = cmd.Flags().GetString(constants.FlagContainerRuntime)
	if err != nil {
		return
	}

	r.imageRepository, err = cmd.Flags().GetString(constants.FlagImageRepository)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagContainerRuntime: r.containerRuntime,
		constants.FlagImageRepository:  r.imageRepository,
	}); err != nil {
		return err
	}

	switch r.containerRuntime {
	case constants.ContainerRuntimeContainerd,
		constants.ContainerRuntimeDocker:
		// break
	default:
		return errors.Wrapf(constants.ErrUnsupportedContainerRuntime, "container runtime: %s", r.containerRuntime)
	}

	return nil
}

func (r *Runtime) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] running\n", r.Use())

	if r.config.ContainerRuntime.Installed {
		_, _ = fmt.Fprintf(out, "[%s] skipping installation (already installed)\n", r.Use())
	}

	switch r.containerRuntime {
	case constants.ContainerRuntimeContainerd:
		return r.installContainerd(out)

	case constants.ContainerRuntimeDocker:
		return r.installDocker(out)

	default:
		return errors.Wrapf(constants.ErrUnsupportedContainerRuntime, "container runtime: %s", r.containerRuntime)
	}
}
