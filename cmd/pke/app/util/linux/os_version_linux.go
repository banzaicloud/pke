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
