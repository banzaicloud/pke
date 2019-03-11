// +build !linux

package kubernetes

import (
	"io"

	"github.com/pkg/errors"
)

func (r *Runtime) installRuntime(w io.Writer, kubernetesVersion string) error {
	return errors.Errorf("unsupported operating system")
}
