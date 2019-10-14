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
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
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
		// /etc/redhat-release
		b, err := ioutil.ReadFile("/etc/redhat-release")
		_, _ = fmt.Fprintf(w, "/etc/redhat-release: %q, err: %v\n", b, err)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`\d+(\.\d+)?`)
		ver := re.Find(b)
		if len(ver) == 0 {
			return "", errors.New("failed to parse version")
		}
		return string(ver), nil
	}

	s := bytes.SplitN(o, dash, 4)
	if len(s) >= 3 {
		if bytes.Equal(redhat, s[0]) && bytes.Equal(release, s[1]) {
			return string(s[2]), nil
		}
	}

	return "", errors.Errorf("unhandled RedHat release output: %s", o)
}

func LSBReleaseDistributorID(w io.Writer) (string, error) {
	o, err := runner.Cmd(w, "/usr/bin/lsb_release", "-si").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(o)), nil
}

func LSBReleaseReleaseNumber(w io.Writer) (string, error) {
	o, err := runner.Cmd(w, "/usr/bin/lsb_release", "-sr").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(o)), nil
}
