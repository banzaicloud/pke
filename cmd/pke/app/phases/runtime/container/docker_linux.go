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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
)

const (
	dockerdPath = "/usr/bin/dockerd"
)

func (r *Runtime) installDocker(out io.Writer) error {
	_ = linux.SystemctlDisableAndStop(out, "docker")

	// Check docker installed or not
	if err := installDocker(out); err != nil {
		return err
	}

	// Start docker
	if err := linux.SystemctlEnableAndStart(out, "docker"); err != nil {
		return err
	}

	_ = linux.SystemctlDisableAndStop(out, "kubelet")

	// systemctl daemon-reload
	return linux.SystemctlReload(out)
}

func installDocker(out io.Writer) error {
	// Check docker installed or not
	if _, err := os.Stat(dockerdPath); !os.IsNotExist(err) {
		_, _ = fmt.Fprintln(out, "docker already installed")
		return nil
	}

	return errors.New("docker installation is not supported")
}
