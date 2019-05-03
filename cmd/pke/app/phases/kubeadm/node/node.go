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
	"text/template"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	pipelineutil "github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-node"
	short = "Kubernetes worker node installation"

	cmdKubeadm         = "/bin/kubeadm"
	kubeProxyConfig    = "/var/lib/kube-proxy/config.conf"
	kubeadmConfig      = "/etc/kubernetes/kubeadm.conf"
	kubeadmAzureConfig = "/etc/kubernetes/azure.conf"
	cniDir             = "/etc/cni/net.d"
	cniBridgeConfig    = "/etc/cni/net.d/10-bridge.conf"
	cniLoopbackConfig  = "/etc/cni/net.d/99-loopback.conf"
)

var _ phases.Runnable = (*Node)(nil)

type Node struct {
	kubernetesVersion      string
	advertiseAddress       string
	apiServerHostPort      string
	kubeadmToken           string
	caCertHash             string
	podNetworkCIDR         string
	cloudProvider          string
	nodepool               string
	azureTenantID          string
	azureSubnetName        string
	azureSecurityGroupName string
	azureVNetName          string
	azureVNetResourceGroup string
	azureVMType            string
	azureLoadBalancerSku   string
	azureRouteTableName    string
	taints                 []string
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
	// Kubernetes network
	flags.String(constants.FlagPodNetworkCIDR, "", "range of IP addresses for the pod network on the current node")
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
	// Azure cloud
	flags.String(constants.FlagAzureTenantID, "", "The AAD Tenant ID for the Subscription that the cluster is deployed in")
	flags.String(constants.FlagAzureSubnetName, "", "The name of the subnet that the cluster is deployed in")
	flags.String(constants.FlagAzureSecurityGroupName, "", "The name of the security group attached to the cluster's subnet")
	flags.String(constants.FlagAzureVNetName, "", "The name of the VNet that the cluster is deployed in")
	flags.String(constants.FlagAzureVNetResourceGroup, "", "The name of the resource group that the Vnet is deployed in")
	flags.String(constants.FlagAzureVMType, "standard", "The type of azure nodes. Candidate values are: vmss and standard")
	flags.String(constants.FlagAzureLoadBalancerSku, "basic", "Sku of Load Balancer and Public IP. Candidate values are: basic and standard")
	flags.String(constants.FlagAzureRouteTableName, "kubernetes-routes", "The name of the route table attached to the subnet that the cluster is deployed in")
	// Taints
	flags.StringSlice(constants.FlagTaints, nil, "Specifies the taints the Node should be registered with")
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

	// Azure specific required flags
	if n.cloudProvider == constants.CloudProviderAzure {
		if err := validator.NotEmpty(map[string]interface{}{
			constants.FlagAzureTenantID:          n.azureTenantID,
			constants.FlagAzureSubnetName:        n.azureSubnetName,
			constants.FlagAzureSecurityGroupName: n.azureSecurityGroupName,
			constants.FlagAzureVNetName:          n.azureVNetName,
			constants.FlagAzureVNetResourceGroup: n.azureVNetResourceGroup,
			constants.FlagAzureVMType:            n.azureVMType,
			constants.FlagAzureLoadBalancerSku:   n.azureLoadBalancerSku,
			constants.FlagAzureRouteTableName:    n.azureRouteTableName,
		}); err != nil {
			return err
		}
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
		n.apiServerHostPort, n.kubeadmToken, n.caCertHash, err = pipelineutil.NodeJoinArgs(os.Stdout, cmd)
		if err != nil {
			return
		}
	}

	n.podNetworkCIDR, err = cmd.Flags().GetString(constants.FlagPodNetworkCIDR)
	if err != nil {
		return
	}
	n.cloudProvider, err = cmd.Flags().GetString(constants.FlagCloudProvider)
	if err != nil {
		return
	}
	n.nodepool, err = cmd.Flags().GetString(constants.FlagPipelineNodepool)
	if err != nil {
		return
	}
	n.azureTenantID, err = cmd.Flags().GetString(constants.FlagAzureTenantID)
	if err != nil {
		return
	}
	n.azureSubnetName, err = cmd.Flags().GetString(constants.FlagAzureSubnetName)
	if err != nil {
		return
	}
	n.azureSecurityGroupName, err = cmd.Flags().GetString(constants.FlagAzureSecurityGroupName)
	if err != nil {
		return
	}
	n.azureVNetName, err = cmd.Flags().GetString(constants.FlagAzureVNetName)
	if err != nil {
		return
	}
	n.azureVNetResourceGroup, err = cmd.Flags().GetString(constants.FlagAzureVNetResourceGroup)
	if err != nil {
		return
	}
	n.azureVMType, err = cmd.Flags().GetString(constants.FlagAzureVMType)
	if err != nil {
		return
	}
	n.azureLoadBalancerSku, err = cmd.Flags().GetString(constants.FlagAzureLoadBalancerSku)
	if err != nil {
		return
	}
	n.azureRouteTableName, err = cmd.Flags().GetString(constants.FlagAzureRouteTableName)
	if err != nil {
		return
	}
	n.taints, err = cmd.Flags().GetStringSlice(constants.FlagTaints)

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

	// write kubeadm azure.conf
	err = kubeadm.WriteKubeadmAzureConfig(out, kubeadmAzureConfig, n.cloudProvider, n.azureTenantID, n.azureSubnetName, n.azureSecurityGroupName, n.azureVNetName, n.azureVNetResourceGroup, n.azureVMType, n.azureLoadBalancerSku, n.azureRouteTableName, true)
	if err != nil {
		return err
	}

	// create cni directory
	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, cniDir)

	if err := os.MkdirAll(cniDir, 0755); err != nil {
		return err
	}

	// CNI network bridge
	if err := writeCNIBridge(out, n.cloudProvider, n.podNetworkCIDR, cniBridgeConfig); err != nil {
		return err
	}

	// CNI network loopback
	if err := writeCNILoopback(out, n.cloudProvider, cniLoopbackConfig); err != nil {
		return err
	}

	// kubeadm join 10.240.0.11:6443 --token 0uk28q.e5i6ewi7xb0g8ye9 --discovery-token-ca-cert-hash sha256:a1a74c00ecccf947b69b49172390018096affbbae25447c4bd0c0906273c1482 --cri-socket=unix:///run/containerd/containerd.sock
	if err := runner.Cmd(out, cmdKubeadm, "join", "--config="+kubeadmConfig).CombinedOutputAsync(); err != nil {
		return err
	}

	return linux.SystemctlEnableAndStart(out, "kubelet")
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

func writeCNIBridge(out io.Writer, cloudProvider, podNetworkCIDR, filename string) error {
	if cloudProvider != constants.CloudProviderAzure || podNetworkCIDR == "" {
		return nil
	}

	conf := `{
    "cniVersion": "0.3.1",
    "name": "bridge",
    "type": "bridge",
    "bridge": "cnio0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "ranges": [
          [{"subnet": "{{ .PodNetworkCIDR }}"}]
        ],
        "routes": [{"dst": "0.0.0.0/0"}]
    }
}`

	tmpl, err := template.New("cni-bridge").Parse(conf)
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
		PodNetworkCIDR string
	}

	d := data{
		PodNetworkCIDR: podNetworkCIDR,
	}

	return tmpl.Execute(w, d)
}

func writeCNILoopback(out io.Writer, cloudProvider, filename string) error {
	if cloudProvider != constants.CloudProviderAzure {
		return nil
	}

	conf := `{
    "cniVersion": "0.3.1",
    "type": "loopback"
}`

	return file.Overwrite(filename, conf)
}
