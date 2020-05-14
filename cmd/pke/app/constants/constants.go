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
	"emperror.dev/errors"
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

	// FlagPipelineAPIInsecure if the Pipeline API should not verify the API's certificate.
	FlagPipelineAPIInsecure = "pipeline-insecure"

	// FlagPipelineOrganizationID organization id in Pipeline.
	FlagPipelineOrganizationID = "pipeline-org-id"
	// FlagPipelineClusterID cluster id in Pipeline.
	FlagPipelineClusterID = "pipeline-cluster-id"

	// FlagPipelineNodepool name of the nodepool the node belongs to.
	FlagPipelineNodepool = "pipeline-nodepool"

	// FlagAdvertiseAddress Kubernetes advertise address API server and Etcd uses this.
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

	// FlagContainerRuntime Kuberneter container runtime.
	FlagContainerRuntime = "kubernetes-container-runtime"

	ContainerRuntimeContainerd = "containerd"
	ContainerRuntimeDocker     = "docker"

	// FlagNetworkProvider network provider for Kubernetes.
	FlagNetworkProvider = "kubernetes-network-provider"
	// FlagServiceCIDR range of IP address for service VIPs.
	FlagServiceCIDR = "kubernetes-service-cidr"
	// FlagPodNetworkCIDR range of IP addresses for the pod network.
	FlagPodNetworkCIDR = "kubernetes-pod-network-cidr"
	// FlagInfrastructureCIDR range of IP addresses from which the advertise address can be calculated using system's network interfaces.
	FlagInfrastructureCIDR = "kubernetes-infrastructure-cidr"
	// FlagMTU maximum transmission unit. 0 means default value of the Kubernetes network provider is used.
	FlagMTU = "kubernetes-mtu"

	NetworkProviderNone   = "none"
	NetworkProviderWeave  = "weave"
	NetworkProviderCalico = "calico"
	NetworkProviderCilium = "cilium"

	// FlagCloudProvider cloud provider for kubeadm.
	FlagCloudProvider = "kubernetes-cloud-provider"

	// CloudProviderAmazon Amazon Web Services
	CloudProviderAmazon = "aws"
	// CloudProviderAzure Azure Cloud Services
	CloudProviderAzure = "azure"
	// CloudProviderVsphere VMware vSphere platform
	CloudProviderVsphere = "vsphere"
	// CloudProviderExternal External cloud provider
	CloudProviderExternal = "external"

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
	// FlagAdditionalControlPlane upgrade additional control plane node.
	FlagAdditionalControlPlane = "kubernetes-additional-control-plane"

	// FlagImageRepository prefix for image repository.
	FlagImageRepository = "image-repository"

	// FlagAdmissionPluginPodSecurityPolicy enable admission plugin PodSecurityPolicy.
	FlagAdmissionPluginPodSecurityPolicy = "with-plugin-psp"

	// FlagNoAdmissionPluginDenyEscalatingExec disable admission plugin DenyEscalatingExec.
	FlagNoAdmissionPluginDenyEscalatingExec = "without-plugin-deny-escalating-exec"

	// FlagAuditLog enable audit log.
	FlagAuditLog = "without-audit-log"

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
	// FlagAzureStorageAccountType Azure storage account Sku tier.
	FlagAzureStorageAccountType = "azure-storage-account-type"
	// FlagAzureStorageKind possible values are shared, dedicated, and managed (default).
	FlagAzureStorageKind = "azure-storage-kind"

	// Vsphere specific flags

	// FlagVsphereServer is the hostname or IP of vCenter to use.
	FlagVsphereServer = "vsphere-server"
	// FlagVspherePort is the TCP port where vCenter listens.
	FlagVspherePort = "vsphere-port"
	// FlagVsphereFingerprint is the fingerprint of the server certificate of vCenter to use.
	FlagVsphereFingerprint = "vsphere-fingerprint"
	// FlagVsphereDatacenter is the name of the datacenter to use to store persistent volumes (and deploy temporary VMs to create them).
	FlagVsphereDatacenter = "vsphere-datacenter"
	// FlagVsphereDatastore is the name of the datastore that is in the given datacenter, and is available on all nodes.
	FlagVsphereDatastore = "vsphere-datastore"
	// FlagVsphereResourcePool is the path of the resource pool to create temporary VMs in during volume creation (for example "Cluster/Pool").
	FlagVsphereResourcePool = "vsphere-resourcepool"
	// FlagVsphereFolder is the name of the folder (aka blue folder) to create temporary VMs in during volume creation as well as all Kubernetes nodes are there.
	FlagVsphereFolder = "vsphere-folder"
	// FlagVsphereUsername is the name of vCenter SSO user to use for deploying persistent volumes. (Should be avoided in favor of a K8S secret.)
	FlagVsphereUsername = "vsphere-username"
	// FlagVspherePassword is the password of vCenter SSO user to use for deploying persistent volumes. (Should be avoided in favor of a K8S secret.)
	FlagVspherePassword = "vsphere-password"

	// FlagLbRange is the IPv4 range advertised via ARP and allocated for LoadBalancer Services.
	FlagLbRange = "lb-range"

	// FlagDisableDefaultStorageClass adds default storage class.
	FlagDisableDefaultStorageClass = "disable-default-storage-class"

	// FlagTaints specifies the taints the Node should be registered with.
	FlagTaints = "taints"

	// FlagLabels specifies the labels the Node should be registered with.
	FlagLabels = "kubernetes-node-labels"

	// Etcd specific flags
	// FlagExternalEtcdEndpoints endpoints of etcd members.
	FlagExternalEtcdEndpoints = "etcd-endpoints"
	// FlagExternalEtcdCAFile is an SSL Certificate Authority file used to secure etcd communication.
	FlagExternalEtcdCAFile = "etcd-ca-file"
	// FlagExternalEtcdCertFile is an SSL certification file used to secure etcd communication.
	FlagExternalEtcdCertFile = "etcd-cert-file"
	// FlagExternalEtcdKeyFile is an SSL key file used to secure etcd communication.
	FlagExternalEtcdKeyFile = "etcd-key-file"
	// FlagExternalEtcdPrefix the prefix to prepend to all resource paths in etcd.
	FlagExternalEtcdPrefix = "etcd-prefix"
	// FlagEncryptionSecret use this key to encrypt secrets.
	FlagEncryptionSecret = "encryption-secret"
)

var (
	ErrUnsupportedOS                = errors.New("unsupported operating system")
	ErrInvalidInput                 = errors.New("invalid input")
	ErrValidationFailed             = errors.New("validation failed")
	ErrUnsupportedContainerRuntime  = errors.New("unsupported container runtime")
	ErrUnsupportedNetworkProvider   = errors.New("unsupported network provider")
	ErrUnsupportedKubernetesVersion = errors.New("unsupported kubernetes version")
	ErrUnsupportedKernelVersion     = errors.New("unsupported kernel version")
)
