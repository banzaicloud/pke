package kubeadm

import (
	"fmt"
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdKubeadm = "/bin/kubeadm"
)

func Reset(out io.Writer) error {
	// kubeadm reset --force
	_, _ = fmt.Fprintln(out, "Resetting kubeadm...")
	err := runner.Cmd(out, cmdKubeadm, "reset", "--force", "--cri-socket=unix:///run/containerd/containerd.sock").CombinedOutputAsync()
	if err != nil {
		return err
	}
	return nil
}
