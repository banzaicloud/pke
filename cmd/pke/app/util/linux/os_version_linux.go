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

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

var (
	dash    = []byte("-")
	release = []byte("release")
	centos  = []byte("centos")
	redhat  = []byte("redhat")
)

// CentOSVersion extract CentOS operating system version
func CentOSVersion(w io.Writer) (string, error) {
	o, err := runner.Cmd(w, "rpm", "--query", "centos-release").Output()
	if err != nil {
		return "", err
	}

	s := bytes.SplitN(o, dash, 4)
	if len(s) >= 3 {
		if bytes.Equal(centos, s[0]) && bytes.Equal(release, s[1]) {
			return string(s[2]), nil
		}
	}

	return "", errors.Errorf("unhandled CentOS release output: %s", o)
}

// RedHatVersion extract RedHat operating system version
func RedHatVersion(w io.Writer) (string, error) {
	o, err := runner.Cmd(w, "rpm", "--query", "redhat-release").Output()
	if err != nil {
		return "", err
	}

	s := bytes.SplitN(o, dash, 4)
	if len(s) >= 3 {
		if bytes.Equal(redhat, s[0]) && bytes.Equal(release, s[1]) {
			return string(s[2]), nil
		}
	}

	return "", errors.Errorf("unhandled RedHat release output: %s", o)
}
