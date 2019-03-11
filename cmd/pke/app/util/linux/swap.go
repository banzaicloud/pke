package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	swapoff = "/sbin/swapoff"
	sed     = "/bin/sed"
	fstab   = "/etc/fstab"
)

// SwapOff disables Linux swap.
func SwapOff(out io.Writer) error {
	// swapoff -a
	if err := runner.Cmd(out, swapoff, "-a").Run(); err != nil {
		return err
	}

	// sed -i '/swap/s/^/#/' /etc/fstab
	if err := runner.Cmd(out, sed, "-i", "/swap/s/^/#/", fstab).Run(); err != nil {
		return err
	}

	return nil
}
