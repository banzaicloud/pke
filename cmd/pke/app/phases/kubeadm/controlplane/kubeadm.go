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
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pkg/errors"
)

func (c ControlPlane) WriteKubeadmConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	// API server advertisement
	bindPort := "6443"
	if c.advertiseAddress != "" {
		host, port, err := kubeadm.SplitHostPort(c.advertiseAddress, "6443")
		if err != nil {
			return err
		}
		c.advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if c.apiServerHostPort != "" {
		host, port, err := kubeadm.SplitHostPort(c.apiServerHostPort, "6443")
		if err != nil {
			return err
		}
		c.apiServerHostPort = net.JoinHostPort(host, port)
	}

	ver, err := semver.NewVersion(c.kubernetesVersion)
	if err != nil {
		return errors.Wrapf(err, "unable to parse Kubernetes version %q", c.kubernetesVersion)
	}

	encryptionProviderPrefix := ""
	if ver.LessThan(semver.MustParse("1.13.0")) {
		encryptionProviderPrefix = "experimental-"
	}

	var conf string
	switch ver.Minor() {
	case 12, 13:
		conf = kubeadmConfigV1Alpha3()
	case 14:
		conf = kubeadmConfigV1Beta1()
	default:
		return errors.New(fmt.Sprintf("unsupported Kubernetes version %q for kubeadm", c.kubernetesVersion))
	}

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

	taints, err := kubernetes.ParseTaints(c.taints)
	if err != nil {
		return err
	}

	type data struct {
		APIServerAdvertiseAddress   string
		APIServerBindPort           string
		ControlPlaneEndpoint        string
		APIServerCertSANs           []string
		KubeletCertificateAuthority string
		AdmissionConfig             string
		ClusterName                 string
		KubernetesVersion           string
		ServiceCIDR                 string
		PodCIDR                     string
		CloudProvider               string
		Nodepool                    string
		ControllerManagerSigningCA  string
		OIDCIssuerURL               string
		OIDCClientID                string
		ImageRepository             string
		EncryptionProviderPrefix    string
		WithPluginPSP               bool
		WithAuditLog                bool
		Taints                      []kubernetes.Taint
		AuditLogDir                 string
		AuditPolicyFile             string
	}

	d := data{
		APIServerAdvertiseAddress:   c.advertiseAddress,
		APIServerBindPort:           bindPort,
		ControlPlaneEndpoint:        c.apiServerHostPort,
		APIServerCertSANs:           c.apiServerCertSANs,
		KubeletCertificateAuthority: c.kubeletCertificateAuthority,
		AdmissionConfig:             admissionConfig,
		ClusterName:                 c.clusterName,
		KubernetesVersion:           c.kubernetesVersion,
		ServiceCIDR:                 c.serviceCIDR,
		PodCIDR:                     c.podNetworkCIDR,
		CloudProvider:               c.cloudProvider,
		Nodepool:                    c.nodepool,
		ControllerManagerSigningCA:  c.controllerManagerSigningCA,
		OIDCIssuerURL:               c.oidcIssuerURL,
		OIDCClientID:                c.oidcClientID,
		ImageRepository:             c.imageRepository,
		EncryptionProviderPrefix:    encryptionProviderPrefix,
		WithPluginPSP:               c.withPluginPSP,
		WithAuditLog:                c.withAuditLog,
		Taints:                      taints,
		AuditLogDir:                 auditLogDir,
		AuditPolicyFile:             auditPolicyFile,
	}

	return tmpl.Execute(w, d)

}

func kubeadmConfigV1Beta1() string {
	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
	return `apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
{{ if .APIServerAdvertiseAddress}}
localAPIEndpoint:
  advertiseAddress: "{{ .APIServerAdvertiseAddress }}"
  bindPort: {{ .APIServerBindPort }}{{end}}
nodeRegistration:
  criSocket: "unix:///run/containerd/containerd.sock"
  taints:{{if not .Taints}} []{{end}}{{range .Taints}}
    - key: "{{.Key}}"
      value: "{{.Value}}"
      effect: "{{.Effect}}"{{end}}
  kubeletExtraArgs:
  {{if .Nodepool }}
    node-labels: "nodepool.banzaicloud.io/name={{ .Nodepool }}"{{end}}
    # pod-infra-container-image: {{ .ImageRepository }}/pause:3.1 # only needed by docker
  {{if .CloudProvider }}
    cloud-provider: "{{ .CloudProvider }}"{{end}}
    {{if eq .CloudProvider "azure" }}cloud-config: "/etc/kubernetes/{{ .CloudProvider }}.conf"{{end}}
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
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
clusterName: "{{ .ClusterName }}"
imageRepository: {{ .ImageRepository }}
useHyperKubeImage: true
networking:
  serviceSubnet: "{{ .ServiceCIDR }}"
  podSubnet: "{{ .PodCIDR }}"
  dnsDomain: "cluster.local"
kubernetesVersion: "v{{ .KubernetesVersion }}"
{{ if .ControlPlaneEndpoint }}controlPlaneEndpoint: "{{ .ControlPlaneEndpoint }}"{{end}}
certificatesDir: "/etc/kubernetes/pki"
apiServer:
  {{if .APIServerCertSANs}}certSANs:
  {{range $k, $san := .APIServerCertSANs}}  - "{{ $san }}"
  {{end}}{{end}}
  extraArgs:
    # anonymous-auth: "false"
    profiling: "false"
    enable-admission-plugins: "AlwaysPullImages,DenyEscalatingExec,EventRateLimit,NodeRestriction,ServiceAccount{{ if .WithPluginPSP }},PodSecurityPolicy{{end}}"
    disable-admission-plugins: ""
    admission-control-config-file: "{{ .AdmissionConfig }}"
    audit-log-path: "{{ .AuditLogDir }}/apiserver.log"
    audit-log-maxage: "30"
    audit-log-maxbackup: "10"
    audit-log-maxsize: "100"
    {{ if .WithAuditLog }}audit-policy-file: "{{ .AuditPolicyFile }}"{{ end }}
    service-account-lookup: "true"
    kubelet-certificate-authority: "{{ .KubeletCertificateAuthority }}"
    tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
    {{ .EncryptionProviderPrefix }}encryption-provider-config: "/etc/kubernetes/admission-control/encryption-provider-config.yaml"
{{if (and .OIDCIssuerURL .OIDCClientID) }}
    oidc-issuer-url: "{{ .OIDCIssuerURL }}"
    oidc-client-id: "{{ .OIDCClientID }}"
    oidc-username-claim: "email"
    oidc-username-prefix: "oidc:"
    oidc-groups-claim: "groups"{{end}}
{{ if .CloudProvider }}
    cloud-provider: "{{ .CloudProvider }}"
    cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
  extraVolumes:
{{ if .WithAuditLog }}
    - name: audit-log-dir
      hostPath: {{ .AuditLogDir }}
      mountPath: {{ .AuditLogDir }}
      pathType: DirectoryOrCreate
    - name: audit-policy-file
      hostPath: {{ .AuditPolicyFile }}
      mountPath: {{ .AuditPolicyFile }}
      readOnly: true
      pathType: File{{ end }}
    - name: admission-control-config-file
      hostPath: {{ .AdmissionConfig }}
      mountPath: {{ .AdmissionConfig }}
      readOnly: true
      pathType: File
    - name: admission-control-config-dir
      hostPath: /etc/kubernetes/admission-control/
      mountPath: /etc/kubernetes/admission-control/
      readOnly: true
      pathType: Directory
{{ if .CloudProvider }}
    - name: cloud-config
      hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf
      mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
scheduler:
  extraArgs:
    profiling: "false"
controllerManager:
  extraArgs:
    cluster-name: "{{ .ClusterName }}"
    profiling: "false"
    terminated-pod-gc-threshold: "10"
    feature-gates: "RotateKubeletServerCertificate=true"
    {{ if .ControllerManagerSigningCA }}cluster-signing-cert-file: {{ .ControllerManagerSigningCA }}{{end}}
{{ if .CloudProvider }}
    cloud-provider: "{{ .CloudProvider }}"
    cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf
  extraVolumes:
    - name: cloud-config
      hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf
      mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
etcd:
  local:
    extraArgs:
      peer-auto-tls: "false"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
serverTLSBootstrap: true
`
}

func kubeadmConfigV1Alpha3() string {
	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
	return `apiVersion: kubeadm.k8s.io/v1alpha3
kind: InitConfiguration
{{ if .APIServerAdvertiseAddress}}apiEndpoint:
  advertiseAddress: "{{ .APIServerAdvertiseAddress }}"
  bindPort: {{ .APIServerBindPort }}{{end}}
nodeRegistration:
  criSocket: "unix:///run/containerd/containerd.sock"
  taints:{{if not .Taints}} []{{end}}{{range .Taints}}
    - key: "{{.Key}}"
      value: "{{.Value}}"
      effect: "{{.Effect}}"{{end}}
  kubeletExtraArgs:
{{if .Nodepool }}
    node-labels: "nodepool.banzaicloud.io/name={{ .Nodepool }}"{{end}}
    # pod-infra-container-image: {{ .ImageRepository }}/pause:3.1 # only needed by docker
{{if .CloudProvider }}
    cloud-provider: "{{ .CloudProvider }}"{{end}}
    {{if eq .CloudProvider "azure" }}cloud-config: "/etc/kubernetes/{{ .CloudProvider }}.conf"{{end}}
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
---
apiVersion: kubeadm.k8s.io/v1alpha3
kind: ClusterConfiguration
clusterName: "{{ .ClusterName }}"
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
  # anonymous-auth: "false"
  profiling: "false"
  enable-admission-plugins: "AlwaysPullImages,DenyEscalatingExec,EventRateLimit,NodeRestriction,ServiceAccount{{ if .WithPluginPSP }},PodSecurityPolicy{{end}}"
  disable-admission-plugins: ""
  admission-control-config-file: "{{ .AdmissionConfig }}"
  audit-log-path: "{{ .AuditLogDir }}/apiserver.log"
  audit-log-maxage: "30"
  audit-log-maxbackup: "10"
  audit-log-maxsize: "100"
  {{ if .WithAuditLog }}audit-policy-file: "{{ .AuditPolicyFile }}"{{ end }}
  service-account-lookup: "true"
  kubelet-certificate-authority: "{{ .KubeletCertificateAuthority }}"
  tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
  {{ .EncryptionProviderPrefix }}encryption-provider-config: "/etc/kubernetes/admission-control/encryption-provider-config.yaml"
{{if (and .OIDCIssuerURL .OIDCClientID) }}
  oidc-issuer-url: "{{ .OIDCIssuerURL }}"
  oidc-client-id: "{{ .OIDCClientID }}"
  oidc-username-claim: "email"
  oidc-username-prefix: "oidc:"
  oidc-groups-claim: "groups"{{end}}
{{ if .CloudProvider }}
  cloud-provider: "{{ .CloudProvider }}"
  cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
schedulerExtraArgs:
  profiling: "false"
apiServerExtraVolumes:
{{ if .WithAuditLog }}
  - name: audit-log-dir
    hostPath: {{ .AuditLogDir }}
    mountPath: {{ .AuditLogDir }}
    pathType: DirectoryOrCreate
  - name: audit-policy-file
    hostPath: {{ .AuditPolicyFile }}
    mountPath: {{ .AuditPolicyFile }}
    readOnly: true
    pathType: FileOrCreate{{ end }}
  - name: admission-control-config-file
    hostPath: {{ .AdmissionConfig }}
    mountPath: {{ .AdmissionConfig }}
    writable: false
    pathType: File
  - name: admission-control-config-dir
    hostPath: /etc/kubernetes/admission-control/
    mountPath: /etc/kubernetes/admission-control/
    writable: false
    pathType: Directory
{{ if .CloudProvider }}
  - name: cloud-config
    hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf
    mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
controllerManagerExtraArgs:
  cluster-name: "{{ .ClusterName }}"
  profiling: "false"
  terminated-pod-gc-threshold: "10"
  feature-gates: "RotateKubeletServerCertificate=true"
  {{ if .ControllerManagerSigningCA }}cluster-signing-cert-file: {{ .ControllerManagerSigningCA }}{{end}}
{{ if .CloudProvider }}
  cloud-provider: "{{ .CloudProvider }}"
  cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf
controllerManagerExtraVolumes:
  - name: cloud-config
    hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf
    mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}
etcd:
  local:
    extraArgs:
      peer-auto-tls: "false"
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
serverTLSBootstrap: true
`
}

func writeAuditPolicyFile(out io.Writer) error {
	filename := auditPolicyFile
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	conf := `apiVersion: audit.k8s.io/v1beta1
kind: Policy
rules:
  - level: None
    resources:
      - group: ""
        resources:
          - endpoints
          - services
          - services/status
    users:
      - 'system:kube-proxy'
      - 'system:apiserver'
    verbs:
      - watch

  - level: None
    resources:
      - group: ""
        resources:
          - nodes
          - nodes/status
    userGroups:
      - 'system:nodes'
    verbs:
      - get

  - level: None
    namespaces:
      - kube-system
    resources:
      - group: ""
        resources:
          - endpoints
    users:
      - 'system:kube-controller-manager'
      - 'system:kube-scheduler'
      - 'system:serviceaccount:kube-system:endpoint-controller'
      - 'system:serviceaccount:kube-system:local-path-provisioner-service-account'
      - 'system:apiserver'
    verbs:
      - get
      - update

  - level: None
    resources:
      - group: ""
        resources:
          - namespaces
          - namespaces/status
          - namespaces/finalize
    users:
      - 'system:apiserver'
    verbs:
      - get

  - level: None
    resources:
      - group: metrics.k8s.io
    users:
      - 'system:kube-controller-manager'
    verbs:
      - get
      - list

  - level: None
    nonResourceURLs:
      - '/healthz*'
      - /version
      - '/swagger*'

  - level: None
    resources:
      - group: ""
        resources:
          - events

  - level: Request
    omitStages:
      - RequestReceived
    resources:
      - group: ""
        resources:
          - nodes/status
          - pods/status
    users:
      - kubelet
      - 'system:node-problem-detector'
      - 'system:serviceaccount:kube-system:node-problem-detector'
    verbs:
      - update
      - patch

  - level: Request
    omitStages:
      - RequestReceived
    resources:
      - group: ""
        resources:
          - nodes/status
          - pods/status
    userGroups:
      - 'system:nodes'
    verbs:
      - update
      - patch

  - level: Request
    omitStages:
      - RequestReceived
    users:
      - 'system:serviceaccount:kube-system:namespace-controller'
    verbs:
      - deletecollection

  - level: Metadata
    omitStages:
      - RequestReceived
    resources:
      - group: ""
        resources:
          - secrets
          - configmaps
      - group: authentication.k8s.io
        resources:
          - tokenreviews

  - level: Request
    omitStages:
      - RequestReceived
    resources:
      - group: ""
      - group: admissionregistration.k8s.io
      - group: apiextensions.k8s.io
      - group: apiregistration.k8s.io
      - group: apps
      - group: authentication.k8s.io
      - group: authorization.k8s.io
      - group: autoscaling
      - group: batch
      - group: certificates.k8s.io
      - group: extensions
      - group: metrics.k8s.io
      - group: networking.k8s.io
      - group: policy
      - group: rbac.authorization.k8s.io
      - group: scheduling.k8s.io
      - group: settings.k8s.io
      - group: storage.k8s.io
    verbs:
      - get
      - list
      - watch

  - level: RequestResponse
    omitStages:
      - RequestReceived
    resources:
      - group: ""
      - group: admissionregistration.k8s.io
      - group: apiextensions.k8s.io
      - group: apiregistration.k8s.io
      - group: apps
      - group: authentication.k8s.io
      - group: authorization.k8s.io
      - group: autoscaling
      - group: batch
      - group: certificates.k8s.io
      - group: extensions
      - group: metrics.k8s.io
      - group: networking.k8s.io
      - group: policy
      - group: rbac.authorization.k8s.io
      - group: scheduling.k8s.io
      - group: settings.k8s.io
      - group: storage.k8s.io

  - level: Metadata
    omitStages:
      - RequestReceived
`

	err = file.Overwrite(filename, conf)
	if err != nil {
		return err
	}
	return nil
}
