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

package controlplane

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/.gen/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/node"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/flags"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/network"
	pipelineutil "github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/transport"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/goph/emperror"
	"github.com/lestrrat-go/backoff"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane installation"

	cmdKubeadm                    = "kubeadm"
	cmdKubectl                    = "kubectl"
	weaveNetUrl                   = "https://cloud.weave.works/k8s/net"
	kubeConfig                    = "/etc/kubernetes/admin.conf"
	kubeProxyConfig               = "/var/lib/kube-proxy/config.conf"
	kubeadmConfig                 = "/etc/kubernetes/kubeadm.conf"
	kubeadmAmazonConfig           = "/etc/kubernetes/aws.conf"
	kubeadmAzureConfig            = "/etc/kubernetes/azure.conf"
	kubeadmVsphereConfig          = "/etc/kubernetes/vsphere.conf"
	storageClassConfig            = "/etc/kubernetes/storage-class.yaml"
	kubernetesCASigningCert       = "/etc/kubernetes/pki/cm-signing-ca.crt"
	admissionConfig               = "/etc/kubernetes/admission-control.yaml"
	admissionEventRateLimitConfig = "/etc/kubernetes/admission-control/event-rate-limit.yaml"
	podSecurityPolicyConfig       = "/etc/kubernetes/admission-control/pod-security-policy.yaml"
	certificateAutoApprover       = "/etc/kubernetes/admission-control/deploy-auto-approver.yaml"
	cniDir                        = "/etc/cni/net.d"
	etcdDir                       = "/var/lib/etcd"
	auditPolicyFile               = "/etc/kubernetes/audit-policy-file.yaml"
	auditLogDir                   = "/var/log/audit/apiserver"
	encryptionSecretLength        = 32
	ciliumBpfMountSystemd         = "/etc/systemd/system/sys-fs-bpf.mount"
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	config config.Config

	kubernetesVersion                string
	containerRuntime                 string
	networkProvider                  string
	advertiseAddress                 string
	apiServerHostPort                string
	clusterName                      string
	serviceCIDR                      string
	podNetworkCIDR                   string
	mtu                              uint
	cloudProvider                    string
	nodepool                         string
	controllerManagerSigningCA       string
	clusterMode                      string
	joinControlPlane                 bool
	apiServerCertSANs                []string
	kubeletCertificateAuthority      string
	oidcIssuerURL                    string
	oidcClientID                     string
	imageRepository                  string
	withPluginPSP                    bool
	withoutPluginDenyEscalatingExec  bool
	useHyperKubeImage                bool
	withoutAuditLog                  bool
	node                             *node.Node
	azureTenantID                    string
	azureSubnetName                  string
	azureSecurityGroupName           string
	azureVNetName                    string
	azureVNetResourceGroup           string
	azureVMType                      string
	azureLoadBalancerSku             string
	azureRouteTableName              string
	azureStorageAccountType          string
	azureStorageKind                 string
	azureExcludeMasterFromStandardLB bool
	vsphereServer                    string
	vspherePort                      int
	vsphereFingerprint               string
	vsphereDatacenter                string
	vsphereDatastore                 string
	vsphereResourcePool              string
	vsphereFolder                    string
	vsphereUsername                  string
	vspherePassword                  string
	cidr                             string
	lbRange                          string
	disableDefaultStorageClass       bool
	taints                           []string
	labels                           []string
	etcdEndpoints                    []string
	etcdCAFile                       string
	etcdCertFile                     string
	etcdKeyFile                      string
	etcdPrefix                       string
	encryptionSecret                 string
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&ControlPlane{
		config: config,
		node:   &node.Node{},
	})
}

func NewDefault(kubernetesVersion, imageRepository string) *ControlPlane {
	return &ControlPlane{
		kubernetesVersion: kubernetesVersion,
		imageRepository:   imageRepository,
		node:              &node.Node{},
	}
}

func (c *ControlPlane) Use() string {
	return use
}

func (c *ControlPlane) Short() string {
	return short
}

func (c *ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, c.config.Kubernetes.Version, "Kubernetes version")
	// Kubernetes container runtime
	flags.String(constants.FlagContainerRuntime, c.config.ContainerRuntime.Type, "Kubernetes container runtime")
	// Kubernetes network
	flags.String(constants.FlagNetworkProvider, "calico", "Kubernetes network provider")
	flags.String(constants.FlagAdvertiseAddress, "", "Kubernetes API Server advertise address")
	flags.String(constants.FlagAPIServerHostPort, "", "Kubernetes API Server host port")
	flags.String(constants.FlagServiceCIDR, "10.10.0.0/16", "range of IP address for service VIPs")
	flags.String(constants.FlagPodNetworkCIDR, "10.20.0.0/16", "range of IP addresses for the pod network")
	flags.Uint(constants.FlagMTU, 0, "maximum transmission unit. 0 means default value of the Kubernetes network provider is used")
	// Kubernetes cluster name
	flags.String(constants.FlagClusterName, "pke", "Kubernetes cluster name")
	// Kubernetes certificates
	flags.StringSlice(constants.FlagAPIServerCertSANs, []string{}, "sets extra Subject Alternative Names for the API Server signing cert")
	flags.String(constants.FlagControllerManagerSigningCA, "", "Kubernetes Controller Manager signing cert")
	flags.String(constants.FlagKubeletCertificateAuthority, "/etc/kubernetes/pki/ca.crt", "Path to a cert file for the certificate authority. Used for kubelet server certificate verify.")
	// Kubernetes cluster mode
	flags.String(constants.FlagClusterMode, "default", "Kubernetes cluster mode")
	// Kubernetes cloud provider (optional)
	flags.String(constants.FlagCloudProvider, "", "cloud provider. example: aws")
	// Pipeline nodepool name (optional)
	flags.String(constants.FlagPipelineNodepool, "", "name of the nodepool the node belongs to")
	// OIDC authentication parameters (optional)
	flags.String(constants.FlagOIDCIssuerURL, "", "URL of the OIDC provider which allows the API server to discover public signing keys")
	flags.String(constants.FlagOIDCClientID, "", "A client ID that all OIDC tokens must be issued for")
	// Image repository
	flags.String(constants.FlagImageRepository, "banzaicloud", "Prefix for image repository")
	// PodSecurityPolicy admission plugin
	flags.Bool(constants.FlagAdmissionPluginPodSecurityPolicy, false, "Enable PodSecurityPolicy admission plugin")
	// DenyEscalatingExec admission plugin
	flags.Bool(constants.FlagNoAdmissionPluginDenyEscalatingExec, false, "Disable DenyEscalatingExec admission plugin")

	// AuditLog enable
	flags.Bool(constants.FlagAuditLog, false, "Disable apiserver audit log")
	// Azure cloud
	flags.String(constants.FlagAzureTenantID, "", "The AAD Tenant ID for the Subscription that the cluster is deployed in")
	flags.String(constants.FlagAzureSubnetName, "", "The name of the subnet that the cluster is deployed in")
	flags.String(constants.FlagAzureSecurityGroupName, "", "The name of the security group attached to the cluster's subnet")
	flags.String(constants.FlagAzureVNetName, "", "The name of the VNet that the cluster is deployed in")
	flags.String(constants.FlagAzureVNetResourceGroup, "", "The name of the resource group that the Vnet is deployed in")
	flags.String(constants.FlagAzureVMType, "standard", "The type of azure nodes. Candidate values are: vmss and standard")
	flags.String(constants.FlagAzureLoadBalancerSku, "basic", "Sku of Load Balancer and Public IP. Candidate values are: basic and standard")
	flags.String(constants.FlagAzureRouteTableName, "kubernetes-routes", "The name of the route table attached to the subnet that the cluster is deployed in")
	flags.String(constants.FlagAzureStorageAccountType, "Standard_LRS", "Azure storage account Sku tier")
	flags.String(constants.FlagAzureStorageKind, "dedicated", "Possible values are shared, dedicated, and managed")

	// VMware vSphere specific flags
	flags.String(constants.FlagVsphereServer, "", "The hostname or IP of vCenter to use")
	flags.Int(constants.FlagVspherePort, 443, "The TCP port where vCenter listens")
	flags.String(constants.FlagVsphereFingerprint, "", "The fingerprint of the server certificate of vCenter to use")
	flags.String(constants.FlagVsphereDatacenter, "", "The name of the datacenter to use to store persistent volumes (and deploy temporary VMs to create them)")
	flags.String(constants.FlagVsphereDatastore, "", "The name of the datastore that is in the given datacenter, and is available on all nodes")
	flags.String(constants.FlagVsphereResourcePool, "", `The path of the resource pool to create temporary VMs in during volume creation (for example "Cluster/Pool")`)
	flags.String(constants.FlagVsphereFolder, "", "The name of the folder (aka blue folder) to create temporary VMs in during volume creation, as well as all Kubernetes nodes are in")
	flags.String(constants.FlagVsphereUsername, "", "The name of vCenter SSO user to use for deploying persistent volumes (Should be avoided in favor of a K8S secret)")
	flags.String(constants.FlagVspherePassword, "", "The password of vCenter SSO user to use for deploying persistent volumes (should be avoided in favor of a K8S secret)")

	// Pipeline
	flags.StringP(constants.FlagPipelineAPIEndpoint, constants.FlagPipelineAPIEndpointShort, "", "Pipeline API server url")
	flags.StringP(constants.FlagPipelineAPIToken, constants.FlagPipelineAPITokenShort, "", "Token for accessing Pipeline API")
	flags.Bool(constants.FlagPipelineAPIInsecure, false, "If the Pipeline API should not verify the API's certificate")
	flags.Int32(constants.FlagPipelineOrganizationID, 0, "Organization ID to use with Pipeline API")
	flags.Int32(constants.FlagPipelineClusterID, 0, "Cluster ID to use with Pipeline API")
	flags.String(constants.FlagInfrastructureCIDR, "192.168.64.0/20", "network CIDR for the actual machine")
	// Storage class
	flags.Bool(constants.FlagDisableDefaultStorageClass, false, "Do not deploy a default storage class")
	flags.String(constants.FlagLbRange, "", "Advertise the specified IPv4 range via ARP and allocate addresses for LoadBalancer Services (non-cloud only, example: 192.168.0.100-192.168.0.110)")
	// Taints
	flags.StringSlice(constants.FlagTaints, []string{"node-role.kubernetes.io/master:NoSchedule"}, "Specifies the taints the Node should be registered with")
	// Labels
	flags.StringSlice(constants.FlagLabels, nil, "Specifies the labels the Node should be registered with")
	// External Etcd
	flags.StringSlice(constants.FlagExternalEtcdEndpoints, []string{}, "Endpoints of etcd members")
	flags.String(constants.FlagExternalEtcdCAFile, "", "An SSL Certificate Authority file used to secure etcd communication")
	flags.String(constants.FlagExternalEtcdCertFile, "", "An SSL certification file used to secure etcd communication")
	flags.String(constants.FlagExternalEtcdKeyFile, "", "An SSL key file used to secure etcd communication")
	flags.String(constants.FlagExternalEtcdPrefix, "", "The prefix to prepend to all resource paths in etcd")
	flags.String(constants.FlagEncryptionSecret, "", "Use this key to encrypt secrets (32 byte base64 encoded)")

	c.addHAControlPlaneFlags(flags)
}

func (c *ControlPlane) addHAControlPlaneFlags(flags *pflag.FlagSet) {
	var f = &pflag.FlagSet{}

	c.node.RegisterFlags(f)

	f.VisitAll(func(flag *pflag.Flag) {
		if flags.Lookup(flag.Name) == nil {
			flags.AddFlag(flag)
		}
	})

	flags.Bool(constants.FlagControlPlaneJoin, false, "Join an another control plane node")
}

func (c *ControlPlane) Validate(cmd *cobra.Command) error {
	if err := c.masterBootstrapParameters(cmd); err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: c.kubernetesVersion,
		constants.FlagContainerRuntime:  c.containerRuntime,
		constants.FlagNetworkProvider:   c.networkProvider,
		constants.FlagServiceCIDR:       c.serviceCIDR,
		constants.FlagPodNetworkCIDR:    c.podNetworkCIDR,
		constants.FlagClusterName:       c.clusterName,
		constants.FlagClusterMode:       c.clusterMode,
		constants.FlagImageRepository:   c.imageRepository,
	}); err != nil {
		return err
	}

	// Azure specific required flags
	if c.cloudProvider == constants.CloudProviderAzure {
		if err := validator.NotEmpty(map[string]interface{}{
			constants.FlagAzureTenantID:           c.azureTenantID,
			constants.FlagAzureSubnetName:         c.azureSubnetName,
			constants.FlagAzureSecurityGroupName:  c.azureSecurityGroupName,
			constants.FlagAzureVNetName:           c.azureVNetName,
			constants.FlagAzureVNetResourceGroup:  c.azureVNetResourceGroup,
			constants.FlagAzureVMType:             c.azureVMType,
			constants.FlagAzureLoadBalancerSku:    c.azureLoadBalancerSku,
			constants.FlagAzureRouteTableName:     c.azureRouteTableName,
			constants.FlagAzureStorageAccountType: c.azureStorageAccountType,
			constants.FlagAzureStorageKind:        c.azureStorageKind,
		}); err != nil {
			return err
		}
	}

	switch c.containerRuntime {
	case constants.ContainerRuntimeContainerd,
		constants.ContainerRuntimeDocker:
		// break
	default:
		return errors.Wrapf(constants.ErrUnsupportedContainerRuntime, "container runtime: %s", c.containerRuntime)
	}

	switch c.networkProvider {
	case constants.NetworkProviderWeave,
		constants.NetworkProviderCalico,
		constants.NetworkProviderNone:
		// break
	case constants.NetworkProviderCilium:
		if err := linux.KernelVersionConstraint(cmd.OutOrStdout(), ">=4.9.17-0"); err != nil {
			return err
		}
		// break
	default:
		return errors.Wrapf(constants.ErrUnsupportedNetworkProvider, "network provider: %s", c.networkProvider)
	}

	// Use Controller Manager Signing CA if present (pipeline-certificates step creates it).
	if c.controllerManagerSigningCA == "" {
		_, err := os.Stat(kubernetesCASigningCert)
		if err == nil {
			c.controllerManagerSigningCA = kubernetesCASigningCert
		}
	}

	switch c.clusterMode {
	case "single":
		c.azureExcludeMasterFromStandardLB = false
	case "default":
		// noop
	case "ha":
		if err := c.pipelineJoin(cmd); err != nil {
			return err
		}

		if c.joinControlPlane {
			return c.node.Validate(cmd)
		}

	default:
		return errors.New("Not supported --" + constants.FlagClusterMode + ". Possible values: single, default or ha.")
	}

	flags.PrintFlags(cmd.OutOrStdout(), c.Use(), cmd.Flags())

	return nil
}

func (c *ControlPlane) pipelineJoin(cmd *cobra.Command) error {
	if pipelineutil.Enabled(cmd) {
		// hostname
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		// ip
		ips, err := network.IPv4Addresses()
		if err != nil {
			return err
		}
		ip, err := network.ContainsFirst(c.cidr, ips)
		if err != nil {
			return err
		}

		// Pipeline client
		endpoint, token, insecure, orgID, clusterID, err := pipelineutil.CommandArgs(cmd)
		if err != nil {
			return err
		}

		p := pipelineutil.Client(os.Stdout, endpoint, token, insecure)

		// elect leader
		_, resp, err := p.ClustersApi.PostLeaderElection(context.Background(), orgID, clusterID, pipeline.PostLeaderElectionRequest{
			Hostname: hostname,
			Ip:       ip.String(),
		})

		if err != nil && resp == nil {
			return errors.Wrap(err, "failed to become leader")
		}
		if resp != nil && resp.StatusCode == http.StatusConflict {
			// check if leadership is ours or not
			var leader pipeline.GetLeaderElectionResponse
			leader, resp, err = p.ClustersApi.GetLeaderElection(context.Background(), orgID, clusterID)
			if err != nil {
				return errors.Wrap(err, "failed to get leader")
			}
			if resp.StatusCode == http.StatusNotFound {
				return errors.New("unexpected condition. leader not found and cloud not acquire leadership")
			}
			if resp.StatusCode == http.StatusOK {
				if hostname == leader.Hostname && ip.String() == leader.Ip {
					// we are the leaders, proceed with master installation
					return nil
				}
			}
			// somebody already took leadership

			c.joinControlPlane = true

			policy := backoff.NewExponential(
				backoff.WithInterval(time.Second),
				backoff.WithFactor(2),
				backoff.WithMaxElapsedTime(time.Hour),
				backoff.WithMaxInterval(30*time.Second),
				backoff.WithMaxRetries(0),
			)
			b, cancel := policy.Start(context.Background())
			defer cancel()

			for backoff.Continue(b) {
				// Wait for master to become ready
				var ready pipeline.PkeClusterReadinessResponse
				ready, resp, err = p.ClustersApi.GetReadyPKENode(context.Background(), orgID, clusterID)
				if resp != nil && resp.StatusCode == http.StatusOK && ready.Master.Ready {
					return nil
				}
			}
			// backoff timeout
			return errors.New("timeout exceeded. waiting for master to become ready failed")
		}
	}

	return nil
}

func (c *ControlPlane) appendAdvertiseAddressAsLoopback() error {
	addr := strings.Split(c.apiServerHostPort, ":")[0]

	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if _, err = fmt.Fprintf(f, "127.0.0.1 %s\n", addr); err != nil {
		return err
	}

	return nil
}

func (c *ControlPlane) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[%s] running\n", c.Use())

	if c.clusterMode == "ha" {
		// additional master node
		if c.joinControlPlane {
			// make sure api server stabilized operation and not restarting
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			if err := ensureAPIServerConnection(out, ctx, 5, c.apiServerHostPort); err != nil {
				return err
			}
			// install additional master node
			if err := writeMasterConfig(out, !c.withoutAuditLog, c.kubernetesVersion, c.encryptionSecret); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(out, "[%s] installing additional master node\n", c.Use())
			return c.node.Run(out)
		}

		// initial master node
		_, _ = fmt.Fprintf(out, "[%s] installing initial master node\n", c.Use())
		if err := c.appendAdvertiseAddressAsLoopback(); err != nil {
			return emperror.Wrap(err, "failed to write to /etc/hosts")
		}
	}

	if err := c.installMaster(out); err != nil {
		if rErr := kubeadm.Reset(out, c.containerRuntime); rErr != nil {
			_, _ = fmt.Fprintf(out, "%v\n", rErr)
		}
		return err
	}

	if err := linux.SystemctlEnableAndStart(out, "kubelet"); err != nil {
		return err
	}

	switch c.networkProvider {
	case constants.NetworkProviderWeave:
		if err := installWeave(out, c.cloudProvider, c.podNetworkCIDR, kubeConfig, c.mtu); err != nil {
			return err
		}
	case constants.NetworkProviderCalico:
		if err := installCalico(out, c.podNetworkCIDR, kubeConfig, c.mtu); err != nil {
			return err
		}
	case constants.NetworkProviderCilium:
		if err := installCilium(out, kubeConfig, c.mtu); err != nil {
			return err
		}
	}

	// install MetalLB if specified
	if err := applyLbRange(out, c.lbRange, c.cloudProvider); err != nil {
		return err
	}

	return taintRemoveNoSchedule(out, c.clusterMode, kubeConfig)
}

func ensureAPIServerConnection(out io.Writer, ctx context.Context, successTries int, apiServerHostPort string) error {
	host, port, err := kubeadm.SplitHostPort(apiServerHostPort, "6443")
	if err != nil {
		return err
	}
	apiServerHostPort = net.JoinHostPort(host, port)

	insecureTLS := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	tl := transport.NewLogger(out, insecureTLS)
	c := &http.Client{
		Transport: transport.NewRetryTransport(tl),
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	u := "https://" + apiServerHostPort + "/version"
	for {
		select {
		case <-ticker.C:
			resp, err := c.Get(u)
			if err != nil {
				return err
			}
			_ = resp.Body.Close()
			if resp.StatusCode/100 == 2 {
				successTries--
				if successTries == 0 {
					return nil
				}
			} else {
				successTries++
			}
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "api server connection cloud not be established")
		}
	}
}

// nolint: gocyclo
func (c *ControlPlane) masterBootstrapParameters(cmd *cobra.Command) (err error) {
	c.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	if err != nil {
		return
	}
	ver, err := semver.NewVersion(c.kubernetesVersion)
	if err != nil {
		return
	}
	c.kubernetesVersion = ver.String()

	c.containerRuntime, err = cmd.Flags().GetString(constants.FlagContainerRuntime)
	if err != nil {
		return
	}

	c.networkProvider, err = cmd.Flags().GetString(constants.FlagNetworkProvider)
	if err != nil {
		return
	}
	c.advertiseAddress, err = cmd.Flags().GetString(constants.FlagAdvertiseAddress)
	if err != nil {
		return
	}
	c.apiServerHostPort, err = cmd.Flags().GetString(constants.FlagAPIServerHostPort)
	if err != nil {
		return
	}
	c.serviceCIDR, err = cmd.Flags().GetString(constants.FlagServiceCIDR)
	if err != nil {
		return
	}
	c.podNetworkCIDR, err = cmd.Flags().GetString(constants.FlagPodNetworkCIDR)
	if err != nil {
		return
	}
	c.mtu, err = cmd.Flags().GetUint(constants.FlagMTU)
	if err != nil {
		return
	}
	c.cloudProvider, err = cmd.Flags().GetString(constants.FlagCloudProvider)
	if err != nil {
		return
	}
	c.nodepool, err = cmd.Flags().GetString(constants.FlagPipelineNodepool)
	if err != nil {
		return
	}
	c.controllerManagerSigningCA, err = cmd.Flags().GetString(constants.FlagControllerManagerSigningCA)
	if err != nil {
		return
	}
	c.clusterMode, err = cmd.Flags().GetString(constants.FlagClusterMode)
	if err != nil {
		return
	}
	c.apiServerCertSANs, err = cmd.Flags().GetStringSlice(constants.FlagAPIServerCertSANs)
	if err != nil {
		return
	}
	c.kubeletCertificateAuthority, err = cmd.Flags().GetString(constants.FlagKubeletCertificateAuthority)
	if err != nil {
		return
	}
	c.clusterName, err = cmd.Flags().GetString(constants.FlagClusterName)
	if err != nil {
		return
	}
	c.oidcIssuerURL, err = cmd.Flags().GetString(constants.FlagOIDCIssuerURL)
	if err != nil {
		return
	}
	c.oidcClientID, err = cmd.Flags().GetString(constants.FlagOIDCClientID)
	if err != nil {
		return
	}
	c.imageRepository, err = cmd.Flags().GetString(constants.FlagImageRepository)
	if err != nil {
		return
	}
	c.withPluginPSP, err = cmd.Flags().GetBool(constants.FlagAdmissionPluginPodSecurityPolicy)
	if err != nil {
		return
	}
	c.withoutPluginDenyEscalatingExec, err = cmd.Flags().GetBool(constants.FlagNoAdmissionPluginDenyEscalatingExec)
	if err != nil {
		return
	}
	c.withoutAuditLog, err = cmd.Flags().GetBool(constants.FlagAuditLog)
	if err != nil {
		return
	}
	c.joinControlPlane, err = cmd.Flags().GetBool(constants.FlagControlPlaneJoin)
	if err != nil {
		return
	}
	err = c.azureParameters(cmd)
	if err != nil {
		return
	}
	err = c.vsphereParameters(cmd)
	if err != nil {
		return
	}
	c.cidr, err = cmd.Flags().GetString(constants.FlagInfrastructureCIDR)
	if err != nil {
		return
	}
	c.disableDefaultStorageClass, err = cmd.Flags().GetBool(constants.FlagDisableDefaultStorageClass)
	if err != nil {
		return
	}
	c.lbRange, err = cmd.Flags().GetString(constants.FlagLbRange)
	if err != nil {
		return
	}
	c.taints, err = cmd.Flags().GetStringSlice(constants.FlagTaints)
	if err != nil {
		return
	}
	c.labels, err = cmd.Flags().GetStringSlice(constants.FlagLabels)
	if err != nil {
		return
	}

	return c.etcdParameters(cmd)
}

func (c *ControlPlane) azureParameters(cmd *cobra.Command) (err error) {
	c.azureTenantID, err = cmd.Flags().GetString(constants.FlagAzureTenantID)
	if err != nil {
		return
	}
	c.azureSubnetName, err = cmd.Flags().GetString(constants.FlagAzureSubnetName)
	if err != nil {
		return
	}
	c.azureSecurityGroupName, err = cmd.Flags().GetString(constants.FlagAzureSecurityGroupName)
	if err != nil {
		return
	}
	c.azureVNetName, err = cmd.Flags().GetString(constants.FlagAzureVNetName)
	if err != nil {
		return
	}
	c.azureVNetResourceGroup, err = cmd.Flags().GetString(constants.FlagAzureVNetResourceGroup)
	if err != nil {
		return
	}
	c.azureVMType, err = cmd.Flags().GetString(constants.FlagAzureVMType)
	if err != nil {
		return
	}
	c.azureLoadBalancerSku, err = cmd.Flags().GetString(constants.FlagAzureLoadBalancerSku)
	if err != nil {
		return
	}
	c.azureRouteTableName, err = cmd.Flags().GetString(constants.FlagAzureRouteTableName)
	if err != nil {
		return
	}
	c.azureStorageAccountType, err = cmd.Flags().GetString(constants.FlagAzureStorageAccountType)
	if err != nil {
		return
	}
	c.azureStorageKind, err = cmd.Flags().GetString(constants.FlagAzureStorageKind)

	return
}

func (c *ControlPlane) vsphereParameters(cmd *cobra.Command) (err error) {
	if c.vsphereServer, err = cmd.Flags().GetString(constants.FlagVsphereServer); err != nil {
		return
	}
	if c.vspherePort, err = cmd.Flags().GetInt(constants.FlagVspherePort); err != nil {
		return
	}
	if c.vsphereFingerprint, err = cmd.Flags().GetString(constants.FlagVsphereFingerprint); err != nil {
		return
	}
	if c.vsphereDatacenter, err = cmd.Flags().GetString(constants.FlagVsphereDatacenter); err != nil {
		return
	}
	if c.vsphereDatastore, err = cmd.Flags().GetString(constants.FlagVsphereDatastore); err != nil {
		return
	}
	if c.vsphereResourcePool, err = cmd.Flags().GetString(constants.FlagVsphereResourcePool); err != nil {
		return
	}
	if c.vsphereFolder, err = cmd.Flags().GetString(constants.FlagVsphereFolder); err != nil {
		return
	}
	if c.vsphereUsername, err = cmd.Flags().GetString(constants.FlagVsphereUsername); err != nil {
		return
	}
	if c.vspherePassword, err = cmd.Flags().GetString(constants.FlagVspherePassword); err != nil {
		return
	}
	return
}

func (c *ControlPlane) etcdParameters(cmd *cobra.Command) (err error) {
	c.etcdEndpoints, err = cmd.Flags().GetStringSlice(constants.FlagExternalEtcdEndpoints)
	if err != nil {
		return
	}
	c.etcdCAFile, err = cmd.Flags().GetString(constants.FlagExternalEtcdCAFile)
	if err != nil {
		return
	}
	c.etcdCertFile, err = cmd.Flags().GetString(constants.FlagExternalEtcdCertFile)
	if err != nil {
		return
	}
	c.etcdKeyFile, err = cmd.Flags().GetString(constants.FlagExternalEtcdKeyFile)
	if err != nil {
		return
	}
	c.etcdPrefix, err = cmd.Flags().GetString(constants.FlagExternalEtcdPrefix)
	if err != nil {
		return
	}
	c.encryptionSecret, err = cmd.Flags().GetString(constants.FlagEncryptionSecret)
	if err != nil {
		return
	}

	return validateEncryptionSecret(c.encryptionSecret)
}

func validateEncryptionSecret(encryptionSecret string) error {
	if encryptionSecret != "" {
		b, err := base64.StdEncoding.DecodeString(encryptionSecret)
		if err != nil {
			return errors.Wrapf(err, "Not valid --%s", constants.FlagEncryptionSecret)
		}

		if l := len(b); l != encryptionSecretLength {
			return errors.New(fmt.Sprintf(
				"Not valid --%s length. Expected length after base64 decode: %d, got: %d",
				constants.FlagEncryptionSecret,
				encryptionSecretLength,
				l,
			))
		}
	}

	return nil
}

func (c *ControlPlane) installMaster(out io.Writer) error {
	// create cni directory
	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, cniDir)
	err := os.MkdirAll(cniDir, 0755)
	if err != nil {
		return err
	}

	// create etcd directory
	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, etcdDir)
	err = os.MkdirAll(etcdDir, 0700)
	if err != nil {
		return err
	}

	// write kubeadm config
	err = c.WriteKubeadmConfig(out, kubeadmConfig)
	if err != nil {
		return err
	}

	// write master config
	if err := writeMasterConfig(out, !c.withoutAuditLog, c.kubernetesVersion, c.encryptionSecret); err != nil {
		return err
	}

	// write kubeadm aws.conf
	err = kubeadm.WriteKubeadmAmazonConfig(out, kubeadmAmazonConfig, c.cloudProvider)
	if err != nil {
		return err
	}

	// write kubeadm azure.conf
	err = kubeadm.WriteKubeadmAzureConfig(out, kubeadmAzureConfig, c.cloudProvider, c.azureTenantID, c.azureSubnetName, c.azureSecurityGroupName, c.azureVNetName, c.azureVNetResourceGroup, c.azureVMType, c.azureLoadBalancerSku, c.azureRouteTableName, c.azureExcludeMasterFromStandardLB)
	if err != nil {
		return err
	}

	// write vsphere.conf
	err = kubeadm.WriteKubeadmVsphereConfig(out, kubeadmVsphereConfig, c.cloudProvider, c.vsphereServer, c.vspherePort, c.vsphereFingerprint, c.vsphereDatacenter, c.vsphereDatastore, c.vsphereResourcePool, c.vsphereFolder, c.vsphereUsername, c.vspherePassword)
	if err != nil {
		return err
	}

	// kubeadm init --config=/etc/kubernetes/kubeadm.conf
	args := []string{
		"init",
		"--config=" + kubeadmConfig,
	}
	_, err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = waitForAPIServer(out)
	if err != nil {
		return err
	}

	// apply AutoApprover
	if err := writeCertificateAutoApprover(out); err != nil {
		return err
	}
	// apply PSP
	if c.withPluginPSP {
		if err := writePodSecurityPolicyConfig(out); err != nil {
			return err
		}
	}

	// if replica set started before default PSP is applied, replica set will hang. force re-create.
	err = deleteKubeDNSReplicaSet(out)
	if err != nil {
		_, _ = fmt.Fprintf(out, "[%s] kube-dns replica set is not started yet, skipping\n", use)
	}

	// apply default storage class
	if err := applyDefaultStorageClass(out, c.disableDefaultStorageClass, c.cloudProvider, c.azureStorageAccountType, c.azureStorageKind); err != nil {
		return err
	}

	return nil
}

//go:generate templify -t ${GOTMPL} -p controlplane -f calico calico.yaml.tmpl

func installCalico(out io.Writer, podNetworkCIDR, kubeConfig string, mtu uint) error {
	input := calicoTemplate()
	input = strings.ReplaceAll(input, "192.168.0.0/16", podNetworkCIDR)
	if mtu > 0 {
		input = strings.ReplaceAll(input, `veth_mtu: "1440"`, fmt.Sprintf(`veth_mtu: "%d"`, mtu))
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", "-")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	cmd.Stdin = strings.NewReader(input)
	_, err := cmd.CombinedOutputAsync()
	return err
}

func installWeave(out io.Writer, cloudProvider, podNetworkCIDR, kubeConfig string, mtu uint) error {
	// kubectl version
	cmd := runner.Cmd(out, cmdKubectl, "version")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// base64
	ver := base64.StdEncoding.EncodeToString(o)

	// weave network url
	u, err := url.Parse(weaveNetUrl)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("k8s-version", ver)
	if cloudProvider != constants.CloudProviderAzure {
		q.Set("env.IPALLOC_RANGE", podNetworkCIDR)
	}
	if mtu > 0 {
		q.Set("env.WEAVE_MTU", strconv.FormatUint(uint64(mtu), 10))
	}
	u.RawQuery = q.Encode()

	// kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')&env.IPALLOC_RANGE=10.200.0.0/16"
	cmd = runner.Cmd(out, cmdKubectl, "apply", "-f", u.String())
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	return err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f cilium cilium.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p controlplane -f ciliumSysFsBpf cilium_sys_fs_bpf.mount.tmpl

func installCilium(out io.Writer, kubeConfig string, mtu uint) error {
	// Mounting BPF filesystem
	if err := file.Overwrite(ciliumBpfMountSystemd, ciliumSysFsBpfTemplate()); err != nil {
		return err
	}
	if err := linux.SystemctlEnableAndStart(out, "sys-fs-bpf.mount"); err != nil {
		return err
	}

	// https://raw.githubusercontent.com/cilium/cilium/v1.6/install/kubernetes/quick-install.yaml
	input := ciliumTemplate()
	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", "-")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	cmd.Stdin = strings.NewReader(input)
	_, err := cmd.CombinedOutputAsync()
	return err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f admissionConfiguration admission_configuration.yaml.tmpl

func writeAdmissionConfiguration(out io.Writer, filename, rateLimitConfigFile string) error {
	tmpl, err := template.New("admission-config").Parse(admissionConfigurationTemplate())
	if err != nil {
		return err
	}

	type data struct {
		RateLimitConfigFile string
	}

	d := data{
		RateLimitConfigFile: rateLimitConfigFile,
	}

	return file.WriteTemplate(filename, tmpl, d)
}

//go:generate templify -t ${GOTMPL} -p controlplane -f kubeProxyConfig kube_proxy_config.yaml.tmpl

func writeKubeProxyConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	return file.Overwrite(filename, kubeProxyConfigTemplate())
}

//go:generate templify -t ${GOTMPL} -p controlplane -f eventRateLimit event_rate_limit.yaml.tmpl

func writeEventRateLimitConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	return file.Overwrite(filename, eventRateLimitTemplate())
}

func waitForAPIServer(out io.Writer) error {
	timeout := 30 * time.Second
	_, _ = fmt.Fprintf(out, "[%s] waiting for API Server to restart. this may take %s\n", use, timeout)

	tout := time.After(timeout)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// kubectl get cs. ensures kube-apiserver is restarted.
			cmd := runner.Cmd(out, cmdKubectl, "get", "cs")
			cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
			_, err := cmd.CombinedOutputAsync()
			if err == nil {
				return nil
			}
		case <-tout:
			return errors.New("wait timeout")
		}
	}
}

func taintRemoveNoSchedule(out io.Writer, clusterMode, kubeConfig string) error {
	if clusterMode != "single" {
		_, _ = fmt.Fprintf(out, "skipping NoSchedule taint removal\n")
		return nil
	}

	// kubectl taint node -l node-role.kubernetes.io/master node-role.kubernetes.io/master:NoSchedule-
	cmd := runner.Cmd(out, cmdKubectl, "taint", "node", "-l node-role.kubernetes.io/master", "node-role.kubernetes.io/master:NoSchedule-")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err := cmd.CombinedOutputAsync()
	return err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f certificateAutoApprover certificate_auto_approver.yaml.tmpl

func writeCertificateAutoApprover(out io.Writer) error {
	filename := certificateAutoApprover
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	err = file.Overwrite(filename, certificateAutoApproverTemplate())
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", filename)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	return err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f podSecurityPolicy pod_security_policy.yaml.tmpl

func writePodSecurityPolicyConfig(out io.Writer) error {
	filename := podSecurityPolicyConfig
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	err = file.Overwrite(filename, podSecurityPolicyTemplate())
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", filename)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	return err
}

func deleteKubeDNSReplicaSet(out io.Writer) error {
	cmd := runner.Cmd(out, cmdKubectl, "delete", "rs", "-n", "kube-system", "k8s-app=kube-dns")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err := cmd.CombinedOutputAsync()
	return err
}

func writeMasterConfig(out io.Writer, a bool, kubernetesVersion, encryptionSecret string) error {

	if a {
		if err := writeAuditPolicyFile(out); err != nil {
			return emperror.Wrap(err, "writing audit policy file failed")
		}
	}

	err := writeAdmissionConfiguration(out, admissionConfig, admissionEventRateLimitConfig)
	if err != nil {
		return emperror.Wrap(err, "writing admission config failed")
	}

	err = writeEventRateLimitConfig(out, admissionEventRateLimitConfig)
	if err != nil {
		return emperror.Wrap(err, "writing event limit config failed")
	}

	err = writeKubeProxyConfig(out, kubeProxyConfig)
	if err != nil {
		return emperror.Wrap(err, "writing kube proxy config failed")
	}

	err = kubeadm.WriteEncryptionProviderConfig(out, kubeadm.EncryptionProviderConfig, kubernetesVersion, encryptionSecret)
	if err != nil {
		return emperror.Wrap(err, "writing encryption provider config failed")
	}

	return nil
}
