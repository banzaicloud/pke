package linux

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdYum = "/bin/yum"
)

func YumInstall(out io.Writer, packages []string) error {
	return runner.Cmd(out, cmdYum, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
}
