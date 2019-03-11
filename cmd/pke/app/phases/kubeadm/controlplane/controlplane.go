package controlplane

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane installation"

	cmdKubeadm              = "/bin/kubeadm"
	cmdKubectl              = "/bin/kubectl"
	weaveNetUrl             = "https://cloud.weave.works/k8s/net"
	kubeConfig              = "/etc/kubernetes/admin.conf"
	kubeadmConfig           = "/etc/kubernetes/kubeadm.conf"
	kubeadmAmazonConfig     = "/etc/kubernetes/aws.conf"
	urlAWSAZ                = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
	kubernetesCASigningCert = "/etc/kubernetes/pki/cm-signing-ca.crt"
)

var _ phases.Runnable = (*ControlPlane)(nil)

type ControlPlane struct {
	kubernetesVersion          string
	networkProvider            string
	advertiseAddress           string
	apiServerHostPort          string
	clusterName                string
	serviceCIDR                string
	podNetworkCIDR             string
	cloudProvider              string
	nodepool                   string
	controllerManagerSigningCA string
	clusterMode                string
	apiServerCertSANs          []string
	oidcIssuerURL              string
	oidcClientID               string
	imageRepository            string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &ControlPlane{})
}

func (c *ControlPlane) Use() string {
	return use
}

func (c *ControlPlane) Short() string {
	return short
}

func (c *ControlPlane) RegisterFlags(flags *pflag.FlagSet) {
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, "1.13.3", "Kubernetes version")
	// Kubernetes network
	flags.String(constants.FlagNetworkProvider, "weave", "Kubernetes network provider")
	flags.String(constants.FlagAdvertiseAddress, "", "Kubernetes advertise address")
	flags.String(constants.FlagAPIServerHostPort, "", "Kubernetes API Server host port")
	flags.String(constants.FlagServiceCIDR, "10.10.0.0/16", "range of IP address for service VIPs")
	flags.String(constants.FlagPodNetworkCIDR, "10.20.0.0/16", "range of IP addresses for the pod network")
	// Kubernetes cluster name
	flags.String(constants.FlagClusterName, "pke", "Kubernetes cluster name")
	// Kubernetes certificates
	flags.StringArray(constants.FlagAPIServerCertSANs, []string{}, "sets extra Subject Alternative Names for the API Server signing cert")
	flags.String(constants.FlagControllerManagerSigningCA, "", "Kubernetes Controller Manager signing cert")
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
}

func (c *ControlPlane) Validate(cmd *cobra.Command) error {
	if err := c.masterBootstrapParameters(cmd); err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		constants.FlagKubernetesVersion: c.kubernetesVersion,
		constants.FlagNetworkProvider:   c.networkProvider,
		constants.FlagServiceCIDR:       c.serviceCIDR,
		constants.FlagPodNetworkCIDR:    c.podNetworkCIDR,
		constants.FlagClusterName:       c.clusterName,
		constants.FlagClusterMode:       c.clusterMode,
		constants.FlagImageRepository:   c.imageRepository,
	}); err != nil {
		return err
	}

	if c.networkProvider != "weave" {
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
	case "single", "default", "ha":
	default:
		return errors.New("Not supported --" + constants.FlagClusterMode + ". Possible values: single, default or ha.")
	}

	return nil
}

func (c *ControlPlane) Run(out io.Writer) error {
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", c.Use())

	if err := installMaster(out, c.kubernetesVersion, c.advertiseAddress, c.apiServerHostPort, c.clusterName, c.serviceCIDR, c.podNetworkCIDR, c.cloudProvider, c.nodepool, c.controllerManagerSigningCA, c.apiServerCertSANs, c.oidcIssuerURL, c.oidcClientID, c.imageRepository); err != nil {
		if rErr := kubeadm.Reset(out); rErr != nil {
			_, _ = fmt.Fprintf(out, "%v\n", rErr)
		}
		return err
	}

	if err := linux.SystemctlEnableAndStart(out, "kubelet"); err != nil {
		return err
	}

	if err := installPodNetwork(out, c.podNetworkCIDR, kubeConfig); err != nil {
		return err
	}

	if err := taintRemoveNoSchedule(out, c.clusterMode, kubeConfig); err != nil {
		return err
	}

	return nil
}

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
	c.apiServerCertSANs, err = cmd.Flags().GetStringArray(constants.FlagAPIServerCertSANs)
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

	return
}

func installMaster(out io.Writer, kubernetesVersion, advertiseAddress, apiServerHostPort, clusterName, serviceCIDR, podNetworkCIDR, cloudProvider, nodepool, controllerManagerSigningCA string, apiServerCertSANs []string, oidcIssuerURL, oidcClientID, imageRepository string) error {
	// write kubeadm config
	err := WriteKubeadmConfig(out, kubeadmConfig, advertiseAddress, apiServerHostPort, clusterName, "", kubernetesVersion, serviceCIDR, podNetworkCIDR, cloudProvider, nodepool, controllerManagerSigningCA, apiServerCertSANs, oidcIssuerURL, oidcClientID, imageRepository)
	if err != nil {
		return err
	}

	// write kubeadm aws.conf
	err = writeKubeadmAmazonConfig(out, kubeadmAmazonConfig, cloudProvider)

	// kubeadm init --config=/etc/kubernetes/kubeadm.conf
	args := []string{
		"init",
		"--config=" + kubeadmConfig,
	}
	err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}

func installPodNetwork(out io.Writer, podNetworkCIDR, kubeConfig string) error {
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
	q.Set("env.IPALLOC_RANGE", podNetworkCIDR)
	u.RawQuery = q.Encode()

	// kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')&env.IPALLOC_RANGE=10.200.0.0/16"
	cmd = runner.Cmd(out, cmdKubectl, "apply", "-f", u.String())
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	if err = cmd.CombinedOutputAsync(); err != nil {
		return err
	}

	return nil
}

func WriteKubeadmConfig(out io.Writer, filename, advertiseAddress, controlPlaneEndpoint, clusterName, fqdn, kubernetesVersion, serviceCIDR, podCIDR, cloudProvider, nodepool, controllerManagerSigningCA string, apiServerCertSANs []string, oidcIssuerURL, oidcClientID, imageRepository string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	// API server advertisement
	bindPort := "6443"
	if advertiseAddress != "" {
		host, port, err := splitHostPort(advertiseAddress, "6443")
		if err != nil {
			return err
		}
		advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if controlPlaneEndpoint != "" {
		host, port, err := splitHostPort(controlPlaneEndpoint, "6443")
		if err != nil {
			return err
		}
		controlPlaneEndpoint = net.JoinHostPort(host, port)
	}

	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
	conf := `apiVersion: kubeadm.k8s.io/v1alpha3
kind: InitConfiguration
{{ if .APIServerAdvertiseAddress}}apiEndpoint:
  advertiseAddress: "{{ .APIServerAdvertiseAddress }}"
  bindPort: {{ .APIServerBindPort }}{{end}}
nodeRegistration:
  criSocket: "unix:///run/containerd/containerd.sock"
  kubeletExtraArgs:
  {{if .Nodepool }}
    node-labels: "nodepool.banzaicloud.io/name={{ .Nodepool }}"{{end}}
  {{if .CloudProvider }}
    cloud-provider: {{ .CloudProvider }}{{end}}
---
apiVersion: kubeadm.k8s.io/v1alpha3
kind: ClusterConfiguration
clusterName: {{ .ClusterName }}
imageRepository: {{ .ImageRepository }}
unifiedControlPlaneImage: {{ .ImageRepository }}/hyperkube:v{{ .KubernetesVersion }}
networking:
  serviceSubnet: "{{ .ServiceCIDR }}"
  podSubnet: "{{ .PodCIDR }}"
  dnsDomain: "cluster.local"
kubernetesVersion: "v{{ .KubernetesVersion }}"
{{ if .ControlPlaneEndpoint }}controlPlaneEndpoint: "{{ .ControlPlaneEndpoint }}"{{end}}
certificatesDir: "/etc/kubernetes/pki"
{{if .APIServerCertSANs}}apiServerCertSANs:
{{range $k, $san := .APIServerCertSANs}}  - "{{ $san }}"
{{end}}{{end}}
apiServerExtraArgs:{{if (and .OIDCIssuerURL .OIDCClientID) }}
  oidc-issuer-url: "{{ .OIDCIssuerURL }}"
  oidc-client-id: "{{ .OIDCClientID }}"
  oidc-username-claim: "email"
  oidc-username-prefix: "oidc:"
  oidc-groups-claim: "groups"{{end}}
{{if eq .CloudProvider "aws" }}
  cloud-provider: aws
  cloud-config: /etc/kubernetes/aws.conf
apiServerExtraVolumes:
  - name: cloud-config
    hostPath: /etc/kubernetes/aws.conf
    mountPath: /etc/kubernetes/aws.conf
controllerManagerExtraVolumes:
  - name: cloud-config
    hostPath: /etc/kubernetes/aws.conf
    mountPath: /etc/kubernetes/aws.conf{{end}}
controllerManagerExtraArgs:{{if eq .CloudProvider "aws" }}
  cloud-provider: aws
  cloud-config: /etc/kubernetes/aws.conf{{end}}
  {{ if .ControllerManagerSigningCA }}cluster-signing-cert-file: {{ .ControllerManagerSigningCA }}{{end}}
`
	tmpl, err := template.New("kubeadm-config").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	type data struct {
		APIServerAdvertiseAddress  string
		APIServerBindPort          string
		ControlPlaneEndpoint       string
		APIServerCertSANs          []string
		ClusterName                string
		FQDN                       string
		KubernetesVersion          string
		ServiceCIDR                string
		PodCIDR                    string
		CloudProvider              string
		Nodepool                   string
		ControllerManagerSigningCA string
		OIDCIssuerURL              string
		OIDCClientID               string
		ImageRepository            string
	}

	d := data{
		APIServerAdvertiseAddress:  advertiseAddress,
		APIServerBindPort:          bindPort,
		ControlPlaneEndpoint:       controlPlaneEndpoint,
		APIServerCertSANs:          apiServerCertSANs,
		ClusterName:                clusterName,
		FQDN:                       fqdn,
		KubernetesVersion:          kubernetesVersion,
		ServiceCIDR:                serviceCIDR,
		PodCIDR:                    podCIDR,
		CloudProvider:              cloudProvider,
		Nodepool:                   nodepool,
		ControllerManagerSigningCA: controllerManagerSigningCA,
		OIDCIssuerURL:              oidcIssuerURL,
		OIDCClientID:               oidcClientID,
		ImageRepository:            imageRepository,
	}

	return tmpl.Execute(w, d)
}

func writeKubeadmAmazonConfig(out io.Writer, filename, cloudProvider string) error {
	if cloudProvider == constants.CloudProviderAmazon {
		if http.DefaultClient.Timeout < 10*time.Second {
			http.DefaultClient.Timeout = 10 * time.Second
		}

		// printf "[GLOBAL]\nZone="$(curl -q -s http://169.254.169.254/latest/meta-data/placement/availability-zone) > /etc/kubernetes/aws.conf
		resp, err := http.Get(urlAWSAZ)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("failed to get aws availability zone. http status code: %d", resp.StatusCode))
		}
		defer func() { _ = resp.Body.Close() }()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer func() { _ = w.Close() }()

		_, err = fmt.Fprintf(w, "[GLOBAL]\nZone=%s\n", b)
		return err
	}

	return nil
}

func taintRemoveNoSchedule(out io.Writer, clusterMode, kubeConfig string) error {
	if clusterMode != "single" {
		_, _ = fmt.Fprintf(out, "skipping NoSchedule taint removal\n")
		return nil
	}

	// kubectl taint node -l node-role.kubernetes.io/master node-role.kubernetes.io/master:NoSchedule-
	cmd := runner.Cmd(out, cmdKubectl, "taint", "node", "-l node-role.kubernetes.io/master", "node-role.kubernetes.io/master:NoSchedule-")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	return cmd.CombinedOutputAsync()
}

func splitHostPort(hostport, defaultPort string) (host, port string, err error) {
	host, port, err = net.SplitHostPort(hostport)
	if aerr, ok := err.(*net.AddrError); ok {
		if aerr.Err == "missing port in address" {
			hostport = net.JoinHostPort(hostport, defaultPort)
			host, port, err = net.SplitHostPort(hostport)
		}
	}
	return
}
