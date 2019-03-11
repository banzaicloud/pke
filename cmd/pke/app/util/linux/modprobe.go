package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdModprobe = "/sbin/modprobe"
)

func Modprobe(out io.Writer, module string) error {
	return runner.Cmd(out, cmdModprobe, module).Run()
}

func ModprobeOverlay(out io.Writer) error {
	return Modprobe(out, "overlay")
}

func ModprobeBFNetFilter(out io.Writer) error {
	return Modprobe(out, "br_netfilter")
}
