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
	"net"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pkg/errors"
)

func (n Node) writeKubeadmConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	// API server advertisement
	bindPort := "6443"
	if n.advertiseAddress != "" {
		host, port, err := kubeadm.SplitHostPort(n.advertiseAddress, "6443")
		if err != nil {
			return err
		}
		n.advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if n.apiServerHostPort != "" {
		host, port, err := kubeadm.SplitHostPort(n.apiServerHostPort, "6443")
		if err != nil {
			return err
		}
		n.apiServerHostPort = net.JoinHostPort(host, port)
	}

	ver, err := semver.NewVersion(n.kubernetesVersion)
	if err != nil {
		return errors.Wrapf(err, "unable to parse Kubernetes version %q", n.kubernetesVersion)
	}

	var conf string
	switch ver.Minor() {
	case 12, 13:
		conf = kubeadmConfigV1Alpha3()
	case 14:
		conf = kubeadmConfigV1Beta1()
	default:
		return errors.Errorf("unsupported Kubernetes version %q for kubeadm", n.kubernetesVersion)
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

	taints, err := kubernetes.ParseTaints(n.taints)
	if err != nil {
		return err
	}

	type data struct {
		APIServerAdvertiseAddress string
		APIServerBindPort         string
		ControlPlaneEndpoint      string
		Token                     string
		CACertHash                string
		CloudProvider             string
		Nodepool                  string
		Taints                    []kubernetes.Taint
	}

	d := data{
		APIServerAdvertiseAddress: n.advertiseAddress,
		APIServerBindPort:         bindPort,
		ControlPlaneEndpoint:      n.apiServerHostPort,
		Token:                     n.kubeadmToken,
		CACertHash:                n.caCertHash,
		CloudProvider:             n.cloudProvider,
		Nodepool:                  n.nodepool,
		Taints:                    taints,
	}

	return tmpl.Execute(w, d)
}

func kubeadmConfigV1Beta1() string {
	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
	return `apiVersion: kubeadm.k8s.io/v1beta1
kind: JoinConfiguration
{{if and .APIServerAdvertiseAddress .APIServerBindPort }}
controlPlane:
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
discovery:
  bootstrapToken:
    apiServerEndpoint: "{{ .ControlPlaneEndpoint }}"
    token: {{ .Token }}
    caCertHashes:
      - {{ .CACertHash }}
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
serverTLSBootstrap: true
`
}

func kubeadmConfigV1Alpha3() string {
	// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
	return `apiVersion: kubeadm.k8s.io/v1alpha3
kind: JoinConfiguration
{{if and .APIServerAdvertiseAddress .APIServerBindPort }}controlPlane: true
apiEndpoint:
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
discoveryTokenAPIServers:
  - {{ .ControlPlaneEndpoint }}
token: {{ .Token }}
discoveryTokenCACertHashes:
  - {{ .CACertHash }}
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
serverTLSBootstrap: true
`
}
