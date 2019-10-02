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
	"bytes"
	"io"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

func KernelVersionConstraint(out io.Writer, constraint string) error {
	// uname -r
	version, err := runner.Cmd(out, "uname", "-r").CombinedOutput()
	if err != nil {
		return err
	}
	v := string(bytes.TrimSpace(version))
	ver, err := semver.NewVersion(v)
	if err != nil {
		return errors.Wrapf(err, "got version: %s", v)
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return err
	}
	if !c.Check(ver) {
		return errors.Wrapf(constants.ErrUnsupportedKernelVersion, "got: %q, expected: %q", v, constraint)
	}

	return nil
}
