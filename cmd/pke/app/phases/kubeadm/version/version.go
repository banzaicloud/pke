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

package version

import (
	"fmt"
	"io"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-version"
	short = "Check Kubernetes version is supported or not"

	constraint = "1.15.x-0 || 1.16.x-0 || 1.17.x-0 || 1.18.x-0"
)

var _ phases.Runnable = (*Version)(nil)

type Version struct {
	config config.Config

	kubernetesVersion string
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&Version{config: config})
}

func (v *Version) Use() string {
	return use
}

func (v *Version) Short() string {
	return short
}

func (v *Version) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, v.config.Kubernetes.Version, "Kubernetes version")
}

func (v *Version) Validate(cmd *cobra.Command) error {
	var err error
	v.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: v.kubernetesVersion,
	}); err != nil {
		return err
	}

	return validVersion(v.kubernetesVersion, constraint)
}

func (v *Version) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] Kubernetes version %q is supported\n", use, v.kubernetesVersion)
	return nil
}

func validVersion(version, constraint string) error {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return err
	}

	ver, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	if !c.Check(ver) {
		return errors.Wrapf(constants.ErrUnsupportedKubernetesVersion, "got: %q, expected: %q", version, constraint)
	}

	return nil
}
