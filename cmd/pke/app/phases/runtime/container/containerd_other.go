// +build !linux

package container

import (
	"io"

	"github.com/pkg/errors"
)

func installRuntime(w io.Writer) error {
	return errors.Errorf("unsupported operating system")
}
