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

// kubeadmConfigV1Beta2Template is a generated function returning the template as a string.
func kubeadmConfigV1Beta2Template() string {
	var tmpl = "apiVersion: kubeadm.k8s.io/v1beta2\n" +
		"kind: JoinConfiguration\n" +
		"{{ if and .APIServerAdvertiseAddress .APIServerBindPort }}\n" +
		"controlPlane:\n" +
		"  localAPIEndpoint:\n" +
		"    advertiseAddress: \"{{ .APIServerAdvertiseAddress }}\"\n" +
		"    bindPort: {{ .APIServerBindPort }}{{end}}\n" +
		"nodeRegistration:\n" +
		"  criSocket: \"{{ .CRISocket }}\"\n" +
		"  taints:{{ if not .Taints }} []{{end}}{{range .Taints}}\n" +
		"    - key: \"{{.Key}}\"\n" +
		"      value: \"{{.Value}}\"\n" +
		"      effect: \"{{.Effect}}\"{{end}}\n" +
		"  kubeletExtraArgs:\n" +
		"    {{ if .NodeLabels }}node-labels: \"{{ .NodeLabels }}\"{{end}}\n" +
		"    {{ if .CloudProvider }}cloud-provider: \"{{ .CloudProvider }}\"{{end}}\n" +
		"    {{if eq .CloudProvider \"azure\" }}cloud-config: \"/etc/kubernetes/{{ .CloudProvider }}.conf\"{{end}}\n" +
		"    read-only-port: \"0\"\n" +
		"    anonymous-auth: \"false\"\n" +
		"    streaming-connection-idle-timeout: \"5m\"\n" +
		"    event-qps: \"0\"\n" +
		"    client-ca-file: \"/etc/kubernetes/pki/ca.crt\"\n" +
		"    feature-gates: \"RotateKubeletServerCertificate=true\"\n" +
		"    rotate-certificates: \"true\"\n" +
		"    tls-cipher-suites: \"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256\"\n" +
		"    authorization-mode: \"Webhook\"\n" +
		"    experimental-kernel-memcg-notification: \"true\"\n" +
		"discovery:\n" +
		"  bootstrapToken:\n" +
		"    apiServerEndpoint: \"{{ .ControlPlaneEndpoint }}\"\n" +
		"    token: {{ .Token }}\n" +
		"    caCertHashes:\n" +
		"      - {{ .CACertHash }}\n" +
		"---\n" +
		"apiVersion: kubelet.config.k8s.io/v1beta2\n" +
		"kind: KubeletConfiguration\n" +
		"serverTLSBootstrap: true\n" +
		"systemReserved:\n" +
		"  cpu: 50m\n" +
		"  memory: 50Mi\n" +
		"  ephemeral-storage: 1Gi\n" +
		"kubeReserved:\n" +
		"  cpu: {{ .KubeReservedCPU }}\n" +
		"  memory: {{ .KubeReservedMemory }}\n" +
		"  ephemeral-storage: 1Gi\n" +
		"evictionHard:\n" +
		"  imagefs.available: 15%\n" +
		"  memory.available: 100Mi\n" +
		"  nodefs.available: 10%\n" +
		"  nodefs.inodesFree: 5%\n" +
		"protectKernelDefaults: true\n" +
		""
	return tmpl
}
