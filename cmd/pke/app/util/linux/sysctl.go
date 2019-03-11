package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdSysctl = "/sbin/sysctl"
)

func SysctlLoadAllFiles(out io.Writer) error {
	// sysctl --system
	return runner.Cmd(out, cmdSysctl, "--system").Run()
}
