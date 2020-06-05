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

package kubernetes

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-runtime"
	short = "Kubernetes runtime installation"
)

var _ phases.Runnable = (*Runtime)(nil)

type Runtime struct {
	config config.Config

	kubernetesVersion string
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
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, r.config.Kubernetes.Version, "Kubernetes version")
}

func (r *Runtime) Validate(cmd *cobra.Command) error {
	var err error
	r.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}
	ver, err := semver.NewVersion(r.kubernetesVersion)
	if err != nil {
		return err
	}
	r.kubernetesVersion = ver.String()

	return validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: r.kubernetesVersion,
	})
}

func (r *Runtime) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] running\n", r.Use())

	if r.config.Kubernetes.Installed {
		_, _ = fmt.Fprintf(out, "[%s] skipping installation (already installed)\n", r.Use())
	}

	return r.installRuntime(out)
}
