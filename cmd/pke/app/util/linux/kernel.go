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
	"io/ioutil"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

func KernelVersionConstraint(out io.Writer, constraint string) error {
	version, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		version, err = runner.Cmd(out, "uname", "-r").CombinedOutput()
	}
	if err != nil {
		return err
	}
	// Red Hat Linux uses underscore: 3.10.0-327.el7.x86_64
	version = bytes.ReplaceAll(version, []byte("_"), []byte(""))
	v := string(bytes.TrimSpace(version))
	ver, err := semver.NewVersion(v)
	if err != nil {
		return errors.Wrapf(err, "got kernel version: %s", v)
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
