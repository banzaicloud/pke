// Copyright © 2020 Banzai Cloud
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

package writeconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	config "github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

const (
	use   = "write-config"
	short = "Write configuration file"
)

var _ phases.Runnable = (*WriteConfig)(nil)

type WriteConfig struct {
	config config.Config

	kubernetesVersion string
	containerRuntime  string
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&WriteConfig{config: config})
}

func (w *WriteConfig) Use() string {
	return use
}

func (w *WriteConfig) Short() string {
	return short
}

func (w *WriteConfig) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, w.config.Kubernetes.Version, "Kubernetes version")

	// Kubernetes container runtime
	flags.String(constants.FlagContainerRuntime, w.config.ContainerRuntime.Type, "Kubernetes container runtime")
}

func (w *WriteConfig) Validate(cmd *cobra.Command) (err error) {
	w.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return
	}

	ver, err := semver.NewVersion(w.kubernetesVersion)
	if err != nil {
		return err
	}
	w.kubernetesVersion = ver.String()

	w.containerRuntime, err = cmd.Flags().GetString(constants.FlagContainerRuntime)
	if err != nil {
		return
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: w.kubernetesVersion,
		constants.FlagContainerRuntime:  w.containerRuntime,
	}); err != nil {
		return err
	}

	switch w.containerRuntime {
	case constants.ContainerRuntimeContainerd,
		constants.ContainerRuntimeDocker:
		// break
	default:
		return errors.Wrapf(constants.ErrUnsupportedContainerRuntime, "container runtime: %s", w.containerRuntime)
	}

	return nil
}

func (w *WriteConfig) Run(out io.Writer) (err error) {
	_, _ = fmt.Fprintf(out, "[%s] running\n", w.Use())

	return w.WriteConfig(out, "/etc/banzaicloud/pke.yaml")
}

func (w *WriteConfig) WriteConfig(_ io.Writer, fileName string) (err error) {
	c := config.Config{
		Kubernetes: config.KubernetesConfig{
			Version:   w.kubernetesVersion,
			Installed: true,
		},
		ContainerRuntime: config.ContainerRuntimeConfig{
			Type:      w.containerRuntime,
			Installed: true,
		},
	}

	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return errors.Wrap(err, "pke config")
	}
	defer func() {
		if err == nil {
			err = file.Close()
		}
	}()

	encoder := yaml.NewEncoder(file)
	defer func() {
		if err == nil {
			err = encoder.Close()
		}
	}()

	return encoder.Encode(c)
}
