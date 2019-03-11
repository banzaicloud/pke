// +build !linux

package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
)

func CentOSVersion(w io.Writer) (string, error) {
	return "", constants.ErrUnsupportedOS
}

func RedHatVersion(w io.Writer) (string, error) {
	return "", constants.ErrUnsupportedOS
}
