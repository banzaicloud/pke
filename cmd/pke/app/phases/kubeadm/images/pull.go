package images

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/controlplane"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "image-pull"
	short = "Pull images used bye PKE tool"

	cmdKubeadm    = "/bin/kubeadm"
	kubeadmConfig = "/etc/kubernetes/kubeadm.conf"
)

var _ phases.Runnable = (*Image)(nil)

type Image struct {
	kubernetesVersion string
	imageRepository   string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Image{})
}

func (i *Image) Use() string {
	return use
}

func (i *Image) Short() string {
	return short
}

func (i *Image) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.13.3", "Kubernetes version")
	// Image repository
	flags.String(constants.FlagImageRepository, "banzaicloud", "Prefix for image repository")
}

func (i *Image) Validate(cmd *cobra.Command) error {
	var err error
	i.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return err
	}
	ver, err := semver.NewVersion(i.kubernetesVersion)
	if err != nil {
		return err
	}
	i.kubernetesVersion = ver.String()

	i.imageRepository, err = cmd.Flags().GetString(constants.FlagImageRepository)
	if err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: i.kubernetesVersion,
		constants.FlagImageRepository:   i.imageRepository,
	}); err != nil {
		return err
	}

	return nil
}

func (i *Image) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", i.Use())

	err := controlplane.WriteKubeadmConfig(out, kubeadmConfig, "", "", "", "", i.kubernetesVersion, "", "", "", "", "", []string{}, "", "", i.imageRepository)
	if err != nil {
		return err
	}
	// kubeadm config images pull --kubernetes-version 1.13.3 --cri-socket unix:///run/containerd/containerd.sock
	err = runner.Cmd(
		out,
		cmdKubeadm,
		"config",
		"images",
		"pull",
		"--config="+kubeadmConfig,
	).CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}
