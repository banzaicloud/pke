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

package constants

import (
	"github.com/pkg/errors"
)

const (
	// Umask default umask
	Umask = 0022

	// FlagOutput output formatting.
	FlagOutput = "output"
	// FlagOutputShort output formatting.
	FlagOutputShort = "o"

	// FlagPipelineAPIEndpoint Pipeline API url.
	FlagPipelineAPIEndpoint = "pipeline-url"
	// FlagPipelineAPIEndpointShort Pipeline API url.
	FlagPipelineAPIEndpointShort = "u"

	// FlagPipelineAPIToken token for accessing Pipeline API.
	FlagPipelineAPIToken = "pipeline-token"
	// FlagPipelineAPITokenShort token for accessing Pipeline API.
	FlagPipelineAPITokenShort = "t"

	// FlagPipelineOrganizationID organization id in Pipeline.
	FlagPipelineOrganizationID = "pipeline-org-id"
	// FlagPipelineClusterID cluster id in Pipeline.
	FlagPipelineClusterID = "pipeline-cluster-id"

	// FlagPipelineNodepool name of the nodepool the node belongs to.
	FlagPipelineNodepool = "pipeline-nodepool"

	// FlagAPIServerHostPort Kubernetes advertise address API server and Etcd uses this.
	FlagAdvertiseAddress = "kubernetes-advertise-address"
	// FlagAPIServerHostPort Kubernetes API Server in host port format.
	FlagAPIServerHostPort = "kubernetes-api-server"
	// FlagKubeadmToken kubeadm token.
	FlagKubeadmToken = "kubernetes-node-token"
	// FlagCACertHash Kubernetes API Server CA Cert hash.
	FlagCACertHash = "kubernetes-api-server-ca-cert-hash"
	// FlagAPIServerCertSANs sets extra Subject Alternative Names for the API Server signing cert.
	FlagAPIServerCertSANs = "kubernetes-api-server-cert-sans"
	// FlagControllerManagerSigningCA Kubernetes Controller Manager needs a single signing cert.
	// This is needed when using Intermediate CAs.
	FlagControllerManagerSigningCA  = "kubernetes-controller-manager-signing-ca"
	FlagKubeletCertificateAuthority = "kubelet-certificate-authority"

	// FlagKubernetesVersion Kubernetes version.
	FlagKubernetesVersion = "kubernetes-version"

	// FlagNetworkProvider network provider for Kubernetes.
	FlagNetworkProvider = "kubernetes-network-provider"
	// FlagServiceCIDR range of IP address for service VIPs.
	FlagServiceCIDR = "kubernetes-service-cidr"
	// FlagPodNetworkCIDR range of IP addresses for the pod network.
	FlagPodNetworkCIDR = "kubernetes-pod-network-cidr"
	// FlagInfrastructureCIDR range of IP addresses from which the advertise address can be calculated using system's network interfaces.
	FlagInfrastructureCIDR = "kubernetes-infrastructure-cidr"

	// FlagCloudProvider cloud provider for kubeadm.
	FlagCloudProvider = "kubernetes-cloud-provider"

	// FlagClusterName cluster name
	FlagClusterName = "kubernetes-cluster-name"

	// FlagOIDCIssuerURL OIDC issuer URL
	FlagOIDCIssuerURL = "kubernetes-oidc-issuer-url"
	// FlagOIDCClientID OIDC client ID
	FlagOIDCClientID = "kubernetes-oidc-client-id"

	// FlagClusterMode possible values: single, default, ha.
	FlagClusterMode = "kubernetes-master-mode"
	// FlagControlPlaneJoin worker command should install control plane node.
	FlagControlPlaneJoin = "kubernetes-join-control-plane"

	// CloudProviderAmazon
	CloudProviderAmazon = "aws"
	// CloudProviderAzure
	CloudProviderAzure = "azure"

	// FlagImageRepository prefix for image repository.
	FlagImageRepository = "image-repository"

	// FlagAdmissionPluginPodSecurityPolicy enable admission plugin PodSecurityPolicy.
	FlagAdmissionPluginPodSecurityPolicy = "with-plugin-psp"

	// Azure specific flags
	// FlagAzureTenantID the AAD Tenant ID for the Subscription that the cluster is deployed in.
	FlagAzureTenantID = "azure-tenant-id"
	// FlagAzureSubnetName the name of the subnet that the cluster is deployed in.
	FlagAzureSubnetName = "azure-subnet-name"
	// FlagAzureSecurityGroupName the name of the security group attached to the cluster's subnet.
	FlagAzureSecurityGroupName = "azure-security-group-name"
	// FlagAzureVNetName the name of the VNet that the cluster is deployed in.
	FlagAzureVNetName = "azure-vnet-name"
	// FlagAzureVNetResourceGroup the name of the resource group that the Vnet is deployed in.
	FlagAzureVNetResourceGroup = "azure-vnet-resource-group"
	// FlagAzureVMType the type of azure nodes. Candidate values are: vmss and standard.
	FlagAzureVMType = "azure-vm-type"
	// FlagAzureLoadBalancerSku sku of Load Balancer and Public IP. Candidate values are: basic and standard.
	FlagAzureLoadBalancerSku = "azure-loadbalancer-sku"
	// FlagAzureRouteTableName the name of the route table attached to the subnet that the cluster is deployed in.
	FlagAzureRouteTableName = "azure-route-table-name"
)

var (
	ErrUnsupportedOS                = errors.New("unsupported operating system")
	ErrInvalidInput                 = errors.New("invalid input")
	ErrValidationFailed             = errors.New("validation failed")
	ErrUnsupportedNetworkProvider   = errors.New("unsupported network provider")
	ErrUnsupportedKubernetesVersion = errors.New("unsupported kubernetes version")
)
