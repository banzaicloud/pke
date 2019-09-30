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

package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

const (
	cmdApt = "/usr/bin/apt-get"
)

var _ ContainerDPackages = (*AptInstaller)(nil)

type AptInstaller struct{}

func NewAptInstaller() *AptInstaller {
	return &AptInstaller{}
}

func (a *AptInstaller) InstallPrerequisites(out io.Writer, containerDVersion string) error {
	// apt-get install -y libseccomp
	if err := AptInstall(out, []string{"libseccomp2"}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}

func AptInstall(out io.Writer, packages []string) error {
	_, err := runner.Cmd(out, cmdApt, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
	return err
}
