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
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"text/template"

	"github.com/banzaicloud/pipeline/client"
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

func (w *Node) Use() string {
	return use
}

func (w *Node) Short() string {
	return short
}

func (w *Node) RegisterFlags(flags *pflag.FlagSet) {
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

func (w *Node) Validate(cmd *cobra.Command) error {
	if err := w.workerBootstrapParameters(cmd); err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagAPIServerHostPort: w.apiServerHostPort,
		constants.FlagKubeadmToken:      w.kubeadmToken,
		constants.FlagCACertHash:        w.caCertHash,
	}); err != nil {
		return err
	}
	return nil
}

func (w *Node) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", w.Use())

	if err := install(out, w.advertiseAddress, w.apiServerHostPort, w.kubeadmToken, w.caCertHash, w.cloudProvider, w.nodepool); err != nil {
		if rErr := kubeadm.Reset(out); rErr != nil {
			_, _ = fmt.Fprintf(out, "%v\n", rErr)
		}
		return err
	}

	return nil
}

func (w *Node) workerBootstrapParameters(cmd *cobra.Command) (err error) {
	// Override values with flags
	w.advertiseAddress, err = cmd.Flags().GetString(constants.FlagAdvertiseAddress)
	if err != nil {
		return
	}
	w.apiServerHostPort, err = cmd.Flags().GetString(constants.FlagAPIServerHostPort)
	if err != nil {
		return
	}
	w.kubeadmToken, err = cmd.Flags().GetString(constants.FlagKubeadmToken)
	if err != nil {
		return
	}
	w.caCertHash, err = cmd.Flags().GetString(constants.FlagCACertHash)
	if err != nil {
		return
	}

	if w.apiServerHostPort == "" && w.kubeadmToken == "" && w.caCertHash == "" {
		w.apiServerHostPort, w.kubeadmToken, w.caCertHash, err = pipelineJoinArgs(cmd)
		if err != nil {
			return
		}
	}

	w.cloudProvider, err = cmd.Flags().GetString(constants.FlagCloudProvider)
	if err != nil {
		return
	}
	w.nodepool, err = cmd.Flags().GetString(constants.FlagPipelineNodepool)

	return
}

func pipelineJoinArgs(cmd *cobra.Command) (apiServerHostPort, kubeadmToken, caCertHash string, err error) {
	if !pipeline.Enabled(cmd) {
		return
	}
	endpoint, token, orgID, clusterID, err := pipeline.CommandArgs(cmd)
	if err != nil {
		return
	}

	// Pipeline client.
	c := pipeline.Client(os.Stdout, endpoint, token)

	var b client.GetClusterBootstrapResponse
	b, _, err = c.ClustersApi.GetClusterBootstrap(context.Background(), orgID, clusterID)
	if err != nil {
		return
	}
	apiServerHostPort = b.MasterAddress
	kubeadmToken = b.Token
	caCertHash = b.DiscoveryTokenCaCertHash
	return
}

func install(out io.Writer, advertiseAddress, apiServerHostPort, token, caCertHash, cloudProvider, nodepool string) error {
	// write kubeadm config
	if err := writeKubeadmConfig(out, kubeadmConfig, advertiseAddress, apiServerHostPort, token, caCertHash, cloudProvider, nodepool); err != nil {
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

func writeKubeadmConfig(out io.Writer, filename, advertiseAddress, controlPlaneEndpoint, token, caCertHash, cloudProvider, nodepool string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	// API server advertisement
	bindPort := "6443"
	if advertiseAddress != "" {
		host, port, err := kubeadm.SplitHostPort(advertiseAddress, "6443")
		if err != nil {
			return err
		}
		advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if controlPlaneEndpoint != "" {
		host, port, err := kubeadm.SplitHostPort(controlPlaneEndpoint, "6443")
		if err != nil {
			return err
		}
		controlPlaneEndpoint = net.JoinHostPort(host, port)
	}

	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
	conf := `apiVersion: kubeadm.k8s.io/v1alpha3
kind: JoinConfiguration
nodeRegistration:
  criSocket: "unix:///run/containerd/containerd.sock"
  kubeletExtraArgs:
  {{if .Nodepool }}
    node-labels: "nodepool.banzaicloud.io/name={{ .Nodepool }}"{{end}}
  {{if .CloudProvider }}
    cloud-provider: "{{ .CloudProvider }}"{{end}}
    read-only-port: "0"
    anonymous-auth: "false"
    streaming-connection-idle-timeout: "5m"
    protect-kernel-defaults: "true"
    event-qps: "0"
    client-ca-file: "/etc/kubernetes/pki/ca.crt"
    feature-gates: "RotateKubeletServerCertificate=true"
    rotate-certificates: "true"
    tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
    authorization-mode: "Webhook"
discoveryTokenAPIServers:
  - {{ .ControlPlaneEndpoint }}
token: {{ .Token }}
discoveryTokenCACertHashes:
  - {{ .CACertHash }}
{{if and .APIServerAdvertiseAddress .APIServerBindPort }}controlPlane: true
apiEndpoint:
  advertiseAddress: "{{ .APIServerAdvertiseAddress }}"
  bindPort: {{ .APIServerBindPort }}{{end}}
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
serverTLSBootstrap: true
`
	tmpl, err := template.New("kubeadm-config").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	type data struct {
		APIServerAdvertiseAddress string
		APIServerBindPort         string
		ControlPlaneEndpoint      string
		Token                     string
		CACertHash                string
		CloudProvider             string
		Nodepool                  string
	}

	d := data{
		APIServerAdvertiseAddress: advertiseAddress,
		APIServerBindPort:         bindPort,
		ControlPlaneEndpoint:      controlPlaneEndpoint,
		Token:                     token,
		CACertHash:                caCertHash,
		CloudProvider:             cloudProvider,
		Nodepool:                  nodepool,
	}

	return tmpl.Execute(w, d)
}
