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

package node

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-node"
	short = "Kubernetes worker node installation"

	cmdKubeadm      = "/bin/kubeadm"
	kubeProxyConfig = "/var/lib/kube-proxy/config.conf"
	kubeadmConfig   = "/etc/kubernetes/kubeadm.conf"
)

var _ phases.Runnable = (*Node)(nil)

type Node struct {
	kubernetesVersion string
	advertiseAddress  string
	apiServerHostPort string
	kubeadmToken      string
	caCertHash        string
	cloudProvider     string
	nodepool          string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Node{})
}

func (n *Node) Use() string {
	return use
}

func (n *Node) Short() string {
	return short
}

func (n *Node) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.14.0", "Kubernetes version")
	// Pipeline
	flags.StringP(constants.FlagPipelineAPIEndpoint, constants.FlagPipelineAPIEndpointShort, "", "Pipeline API server url")
	flags.StringP(constants.FlagPipelineAPIToken, constants.FlagPipelineAPITokenShort, "", "Token for accessing Pipeline API")
	flags.Int32(constants.FlagPipelineOrganizationID, 0, "Organization ID to use with Pipeline API")
	flags.Int32(constants.FlagPipelineClusterID, 0, "Cluster ID to use with Pipeline API")
	// Kubernetes cloud provider (optional)
	flags.String(constants.FlagCloudProvider, "", "cloud provider. example: aws")
	// Control Plane
	flags.String(constants.FlagAdvertiseAddress, "", "Kubernetes API Server advertise address")
	_ = flags.MarkHidden(constants.FlagAdvertiseAddress)
	// Kubernetes cluster join parameters
	flags.String(constants.FlagAPIServerHostPort, "", "Kubernetes API Server host port")
	flags.String(constants.FlagKubeadmToken, "", "PKE join token")
	flags.String(constants.FlagCACertHash, "", "CA cert hash")
	// Pipeline nodepool name (optional)
	flags.String(constants.FlagPipelineNodepool, "", "name of the nodepool the node belongs to")
}

func (n *Node) Validate(cmd *cobra.Command) error {
	if err := n.workerBootstrapParameters(cmd); err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: n.kubernetesVersion,
		constants.FlagAPIServerHostPort: n.apiServerHostPort,
		constants.FlagKubeadmToken:      n.kubeadmToken,
		constants.FlagCACertHash:        n.caCertHash,
	}); err != nil {
		return err
	}
	return nil
}

func (n *Node) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", n.Use())

	if err := n.install(out); err != nil {
		if rErr := kubeadm.Reset(out); rErr != nil {
			_, _ = fmt.Fprintf(out, "%v\n", rErr)
		}
		return err
	}

	return nil
}

func (n *Node) workerBootstrapParameters(cmd *cobra.Command) (err error) {
	n.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return
	}
	// Override values with flags
	n.advertiseAddress, err = cmd.Flags().GetString(constants.FlagAdvertiseAddress)
	if err != nil {
		return
	}
	n.apiServerHostPort, err = cmd.Flags().GetString(constants.FlagAPIServerHostPort)
	if err != nil {
		return
	}
	n.kubeadmToken, err = cmd.Flags().GetString(constants.FlagKubeadmToken)
	if err != nil {
		return
	}
	n.caCertHash, err = cmd.Flags().GetString(constants.FlagCACertHash)
	if err != nil {
		return
	}

	if n.kubeadmToken == "" && n.caCertHash == "" {
		n.apiServerHostPort, n.kubeadmToken, n.caCertHash, err = pipeline.NodeJoinArgs(os.Stdout, cmd)
		if err != nil {
			return
		}
	}

	n.cloudProvider, err = cmd.Flags().GetString(constants.FlagCloudProvider)
	if err != nil {
		return
	}
	n.nodepool, err = cmd.Flags().GetString(constants.FlagPipelineNodepool)

	return
}

func (n *Node) install(out io.Writer) error {
	// write kubeadm config
	if err := n.writeKubeadmConfig(out, kubeadmConfig); err != nil {
		return err
	}

	err := writeKubeProxyConfig(out, kubeProxyConfig)
	if err != nil {
		return err
	}

	// kubeadm join 10.240.0.11:6443 --token 0uk28q.e5i6ewi7xb0g8ye9 --discovery-token-ca-cert-hash sha256:a1a74c00ecccf947b69b49172390018096affbbae25447c4bd0c0906273c1482 --cri-socket=unix:///run/containerd/containerd.sock
	if err := runner.Cmd(out, cmdKubeadm, "join", "--config="+kubeadmConfig).CombinedOutputAsync(); err != nil {
		return err
	}

	if err := linux.SystemctlEnableAndStart(out, "kubelet"); err != nil {
		return err
	}

	return nil
}

func writeKubeProxyConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	conf := `apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
`

	return file.Overwrite(filename, conf)
}
