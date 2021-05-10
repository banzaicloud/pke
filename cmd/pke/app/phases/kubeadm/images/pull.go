// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package images

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
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
	short = "Pull images used by PKE tool"

	cmdKubeadm    = "kubeadm"
	kubeadmConfig = "/etc/kubernetes/kubeadm.conf"
)

var _ phases.Runnable = (*Image)(nil)

type Image struct {
	config config.Config

	kubernetesVersion       string
	imageRepository         string
	useImageRepositoryToK8s bool
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&Image{config: config})
}

func (i *Image) Use() string {
	return use
}

func (i *Image) Short() string {
	return short
}

func (i *Image) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, i.config.Kubernetes.Version, "Kubernetes version")
	// Image repository
	flags.String(constants.FlagImageRepository, "banzaicloud", "Prefix for image repository")
	// Use defined image repository for K8s images as well
	flags.Bool(constants.FlagUseImageRepositoryToK8s, false, "Use defined image repository for K8s Images as well")
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
	i.useImageRepositoryToK8s, err = cmd.Flags().GetBool(constants.FlagUseImageRepositoryToK8s)
	if err != nil {
		return err
	}
	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion:       i.kubernetesVersion,
		constants.FlagImageRepository:         i.imageRepository,
		constants.FlagUseImageRepositoryToK8s: i.useImageRepositoryToK8s,
	}); err != nil {
		return err
	}

	return nil
}

func (i *Image) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] running\n", i.Use())

	imageRepository := "k8s.gcr.io"
	if i.useImageRepositoryToK8s {
		imageRepository = i.imageRepository
	}
	c := controlplane.NewDefault(i.kubernetesVersion, imageRepository)

	err := c.WriteKubeadmConfig(out, kubeadmConfig)
	if err != nil {
		return err
	}
	// kubeadm config images pull --kubernetes-version 1.14.0 --cri-socket unix:///run/containerd/containerd.sock
	_, err = runner.Cmd(
		out,
		cmdKubeadm,
		"config",
		"images",
		"pull",
		"--config="+kubeadmConfig,
	).CombinedOutputAsync()
	return err
}
