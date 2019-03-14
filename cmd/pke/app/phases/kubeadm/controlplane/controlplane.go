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
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	use   = "kubernetes-controlplane"
	short = "Kubernetes Control Plane installation"

	cmdKubeadm                    = "/bin/kubeadm"
	cmdKubectl                    = "/bin/kubectl"
	weaveNetUrl                   = "https://cloud.weave.works/k8s/net"
	kubeConfig                    = "/etc/kubernetes/admin.conf"
	kubeadmConfig                 = "/etc/kubernetes/kubeadm.conf"
	kubeadmAmazonConfig           = "/etc/kubernetes/aws.conf"
	urlAWSAZ                      = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
	kubernetesCASigningCert       = "/etc/kubernetes/pki/cm-signing-ca.crt"
	admissionConfig               = "/etc/kubernetes/admission-control.yaml"
	admissionEventRateLimitConfig = "/etc/kubernetes/admission-control/event-rate-limit.yaml"
	apiServerManifest             = "/etc/kubernetes/manifests/kube-apiserver.yaml"
	podSecurityPolicyConfig       = "/etc/kubernetes/admission-control/pod-security-policy.yaml"
	cniDir                        = "/etc/cni/net.d"
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
	// create cni directory
	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, cniDir)
	err := os.MkdirAll(cniDir, 0644)
	if err != nil {
		return err
	}

	// write kubeadm config
	err = WriteKubeadmConfig(out, kubeadmConfig, advertiseAddress, apiServerHostPort, admissionConfig, clusterName, "", kubernetesVersion, serviceCIDR, podNetworkCIDR, cloudProvider, nodepool, controllerManagerSigningCA, apiServerCertSANs, oidcIssuerURL, oidcClientID, imageRepository)
	if err != nil {
		return err
	}

	err = writeAdmissionConfiguration(out, admissionConfig, admissionEventRateLimitConfig)
	if err != nil {
		return err
	}

	err = writeEventRateLimitConfig(out, admissionEventRateLimitConfig)
	if err != nil {
		return err
	}

	// write kubeadm aws.conf
	err = writeKubeadmAmazonConfig(out, kubeadmAmazonConfig, cloudProvider)
	if err != nil {
		return err
	}

	// kubeadm init --config=/etc/kubernetes/kubeadm.conf
	args := []string{
		"init",
		"--config=" + kubeadmConfig,
	}
	err = runner.Cmd(out, cmdKubeadm, args...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	// --anonymous-auth=false implies this
	err = replaceLivelinessProbe(out, apiServerManifest)
	if err != nil {
		return err
	}

	err = waitForAPIServer(out)
	if err != nil {
		return err
	}

	// apply PSP
	if err := writePodSecurityPolicyConfig(out); err != nil {
		return err
	}

	// if replica set started before default PSP is applied, replica set will hang. force re-create.
	err = deleteKubeDNSReplicaSet(out)
	if err != nil {
		_, _ = fmt.Fprintf(out, "[%s] kube-dns replica set is not started yet, skipping\n", use)
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

func WriteKubeadmConfig(out io.Writer, filename, advertiseAddress, controlPlaneEndpoint, admissionConfig, clusterName, fqdn, kubernetesVersion, serviceCIDR, podCIDR, cloudProvider, nodepool, controllerManagerSigningCA string, apiServerCertSANs []string, oidcIssuerURL, oidcClientID, imageRepository string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0640)
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
apiServerExtraArgs:
  anonymous-auth: "false"
  profiling: "false"
  enable-admission-plugins: "AlwaysPullImages,DenyEscalatingExec,EventRateLimit,NamespaceLifecycle,NodeRestriction,PodSecurityPolicy,ServiceAccount"
  admission-control-config-file: "{{ .AdmissionConfig }}"
  audit-log-path: "/var/log/audit/apiserver.log"
  audit-log-maxage: "30"
  audit-log-maxbackup: "10"
  audit-log-maxsize: "100"
  service-account-lookup: "true"
  tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
{{if (and .OIDCIssuerURL .OIDCClientID) }}
  oidc-issuer-url: "{{ .OIDCIssuerURL }}"
  oidc-client-id: "{{ .OIDCClientID }}"
  oidc-username-claim: "email"
  oidc-username-prefix: "oidc:"
  oidc-groups-claim: "groups"{{end}}
{{if eq .CloudProvider "aws" }}
  cloud-provider: aws
  cloud-config: /etc/kubernetes/aws.conf{{end}}
schedulerExtraArgs:
  profiling: "false"
apiServerExtraVolumes:
  - name: admission-control-config-file
    hostPath: /etc/kubernetes/admission-control.yaml
    mountPath: /etc/kubernetes/admission-control.yaml
    writable: false
    pathType: File
  - name: admission-control-config-dir
    hostPath: /etc/kubernetes/admission-control/
    mountPath: /etc/kubernetes/admission-control/
    writable: false
    pathType: Directory{{if eq .CloudProvider "aws" }}
  - name: cloud-config
    hostPath: /etc/kubernetes/aws.conf
    mountPath: /etc/kubernetes/aws.conf
controllerManagerExtraVolumes:
  - name: cloud-config
    hostPath: /etc/kubernetes/aws.conf
    mountPath: /etc/kubernetes/aws.conf{{end}}
controllerManagerExtraArgs:
  profiling: "false"
  terminated-pod-gc-threshold: "10"
  feature-gates: "RotateKubeletServerCertificate=true"{{if eq .CloudProvider "aws" }}
  cloud-provider: aws
  cloud-config: /etc/kubernetes/aws.conf{{end}}
  {{ if .ControllerManagerSigningCA }}cluster-signing-cert-file: {{ .ControllerManagerSigningCA }}{{end}}
etcd:
  local:
    extraArgs:
      peer-auto-tls: "false"
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
		APIServerAdvertiseAddress  string
		APIServerBindPort          string
		ControlPlaneEndpoint       string
		APIServerCertSANs          []string
		AdmissionConfig            string
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
		AdmissionConfig:            admissionConfig,
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

		w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
		if err != nil {
			return err
		}
		defer func() { _ = w.Close() }()

		_, err = fmt.Fprintf(w, "[GLOBAL]\nZone=%s\n", b)
		return err
	}

	return nil
}

func writeAdmissionConfiguration(out io.Writer, filename, rateLimitConfigFile string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0640)
	if err != nil {
		return err
	}

	conf := `kind: AdmissionConfiguration
apiVersion: apiserver.k8s.io/v1alpha1
plugins:
- name: EventRateLimit
  path: {{ .RateLimitConfigFile }}
`

	tmpl, err := template.New("admission-config").Parse(conf)
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
		RateLimitConfigFile string
	}

	d := data{
		RateLimitConfigFile: rateLimitConfigFile,
	}

	return tmpl.Execute(w, d)
}

func writeEventRateLimitConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0640)
	if err != nil {
		return err
	}

	conf := `kind: Configuration
apiVersion: eventratelimit.admission.k8s.io/v1alpha1
limits:
- type: Namespace
  qps: 50
  burst: 100
  cacheSize: 2000
- type: User
  qps: 10
  burst: 50
`

	return file.Overwrite(filename, conf)
}

func replaceLivelinessProbe(out io.Writer, filename string) error {
	_, _ = fmt.Fprintf(out, "[%s] changing apiserver liveliness probe to tcp based: %q\n", use, filename)

	// read configuration
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// convert json to yaml
	j, err := yaml.YAMLToJSON(f)
	if err != nil {
		return err
	}
	sj := string(j)

	// get port value
	port := gjson.Get(sj, "spec.containers.0.livenessProbe.httpGet.port")

	// remove http get liveliness probe
	sj, err = sjson.Delete(sj, "spec.containers.0.livenessProbe.httpGet")
	if err != nil {
		return err
	}

	// add tcp probe
	sj, err = sjson.Set(sj, "spec.containers.0.livenessProbe.tcpSocket.port", port.Num)
	if err != nil {
		return err
	}

	// convert json to yaml
	j, err = yaml.JSONToYAML([]byte(sj))
	if err != nil {
		return err
	}

	// overwrite file
	return file.Overwrite(filename, string(j))
}

func waitForAPIServer(out io.Writer) error {
	timeout := 30 * time.Second
	_, _ = fmt.Fprintf(out, "[%s] waiting for API Server to restart. this may take %s\n", use, timeout)

	tout := time.After(timeout)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			// kubectl get cs. ensures kube-apiserver is restarted.
			cmd := runner.Cmd(out, cmdKubectl, "get", "cs")
			cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
			err := cmd.CombinedOutputAsync()
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
	return cmd.CombinedOutputAsync()
}

func writePodSecurityPolicyConfig(out io.Writer) error {
	filename := podSecurityPolicyConfig
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0640)
	if err != nil {
		return err
	}

	conf := `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:psp-binding
roleRef:
  kind: ClusterRole
  name: system:psp:privileged
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:serviceaccounts
- kind: Group
  apiGroup: rbac.authorization.k8s.io
  name: system:nodes
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:psp:privileged
rules:
- apiGroups:
  - extensions
  resourceNames:
  - pke-psp
  resources:
  - podsecuritypolicies
  verbs:
  - use
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: pke-psp
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'
spec:
  privileged: true
  allowPrivilegeEscalation: true
  allowedCapabilities:
  - '*'
  volumes:
  - '*'
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  hostIPC: true
  hostPID: true
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
    rule: 'RunAsAny'
`

	err = file.Overwrite(filename, conf)
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", filename)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	return cmd.CombinedOutputAsync()
}

func deleteKubeDNSReplicaSet(out io.Writer) error {
	cmd := runner.Cmd(out, cmdKubectl, "delete", "rs", "-n", "kube-system", "k8s-app=kube-dns")
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
