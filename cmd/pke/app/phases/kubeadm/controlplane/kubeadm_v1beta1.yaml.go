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

// kubeadmConfigV1Beta1Template is a generated function returning the template as a string.
func kubeadmConfigV1Beta1Template() string {
	var tmpl = "apiVersion: kubeadm.k8s.io/v1beta1\n" +
		"kind: InitConfiguration\n" +
		"{{ if .APIServerAdvertiseAddress}}\n" +
		"localAPIEndpoint:\n" +
		"  advertiseAddress: \"{{ .APIServerAdvertiseAddress }}\"\n" +
		"  bindPort: {{ .APIServerBindPort }}{{end}}\n" +
		"nodeRegistration:\n" +
		"  criSocket: \"{{ .CRISocket }}\"\n" +
		"  taints:{{ if not .Taints }} []{{end}}{{range .Taints}}\n" +
		"    - key: \"{{.Key}}\"\n" +
		"      value: \"{{.Value}}\"\n" +
		"      effect: \"{{.Effect}}\"{{end}}\n" +
		"  kubeletExtraArgs:\n" +
		"    {{ if .NodeLabels }}node-labels: \"{{ .NodeLabels }}\"{{end}}\n" +
		"    # pod-infra-container-image: {{ .ImageRepository }}/pause:3.1 # only needed by docker\n" +
		"    {{ if .CloudProvider }}cloud-provider: \"{{ .CloudProvider }}\"\n" +
		"    {{ if .KubeletCloudConfig }}cloud-config: \"/etc/kubernetes/{{ .CloudProvider }}.conf\"{{end}}{{end}}\n" +
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
		"---\n" +
		"apiVersion: kubeadm.k8s.io/v1beta1\n" +
		"kind: ClusterConfiguration\n" +
		"clusterName: \"{{ .ClusterName }}\"\n" +
		"imageRepository: {{ .ImageRepository }}\n" +
		"useHyperKubeImage:  {{ .UseHyperKubeImage }}\n" +
		"networking:\n" +
		"  serviceSubnet: \"{{ .ServiceCIDR }}\"\n" +
		"  podSubnet: \"{{ .PodCIDR }}\"\n" +
		"  dnsDomain: \"cluster.local\"\n" +
		"kubernetesVersion: \"v{{ .KubernetesVersion }}\"\n" +
		"{{ if .ControlPlaneEndpoint }}controlPlaneEndpoint: \"{{ .ControlPlaneEndpoint }}\"{{end}}\n" +
		"certificatesDir: \"/etc/kubernetes/pki\"\n" +
		"apiServer:\n" +
		"  {{ if .APIServerCertSANs }}\n" +
		"  certSANs:\n" +
		"  {{range $k, $san := .APIServerCertSANs}}  - \"{{ $san }}\"\n" +
		"  {{end}}{{end}}\n" +
		"  extraArgs:\n" +
		"    # anonymous-auth: \"false\"\n" +
		"    profiling: \"false\"\n" +
		"    enable-admission-plugins: \"AlwaysPullImages,{{ if not .WithoutPluginDenyEscalatingExec }}DenyEscalatingExec,{{end}}EventRateLimit,NodeRestriction,ServiceAccount{{ if .WithPluginPSP }},PodSecurityPolicy{{end}}\"\n" +
		"    disable-admission-plugins: \"\"\n" +
		"    admission-control-config-file: \"{{ .AdmissionConfig }}\"\n" +
		"    audit-log-path: \"{{ .AuditLogDir }}/apiserver.log\"\n" +
		"    audit-log-maxage: \"30\"\n" +
		"    audit-log-maxbackup: \"10\"\n" +
		"    audit-log-maxsize: \"100\"\n" +
		"    {{ if .WithAuditLog }}audit-policy-file: \"{{ .AuditPolicyFile }}\"{{ end }}\n" +
		"    {{ if .EtcdPrefix }}etcd-prefix: \"{{ .EtcdPrefix }}\"{{end}}\n" +
		"    service-account-lookup: \"true\"\n" +
		"    kubelet-certificate-authority: \"{{ .KubeletCertificateAuthority }}\"\n" +
		"    tls-cipher-suites: \"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256\"\n" +
		"    {{ .EncryptionProviderPrefix }}encryption-provider-config: \"/etc/kubernetes/admission-control/encryption-provider-config.yaml\"\n" +
		"    {{ if (and .OIDCIssuerURL .OIDCClientID) }}\n" +
		"    oidc-issuer-url: \"{{ .OIDCIssuerURL }}\"\n" +
		"    oidc-client-id: \"{{ .OIDCClientID }}\"\n" +
		"    oidc-username-claim: \"email\"\n" +
		"    oidc-username-prefix: \"oidc:\"\n" +
		"    oidc-groups-claim: \"groups\"{{end}}\n" +
		"    {{ if .CloudProvider }}cloud-provider: \"{{ .CloudProvider }}\"\n" +
		"    {{ if .CloudConfig }}cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}{{end}}\n" +
		"  extraVolumes:\n" +
		"    {{ if .WithAuditLog }}\n" +
		"    - name: audit-log-dir\n" +
		"      hostPath: {{ .AuditLogDir }}\n" +
		"      mountPath: {{ .AuditLogDir }}\n" +
		"      pathType: DirectoryOrCreate\n" +
		"    - name: audit-policy-file\n" +
		"      hostPath: {{ .AuditPolicyFile }}\n" +
		"      mountPath: {{ .AuditPolicyFile }}\n" +
		"      readOnly: true\n" +
		"      pathType: File{{ end }}\n" +
		"    - name: admission-control-config-file\n" +
		"      hostPath: {{ .AdmissionConfig }}\n" +
		"      mountPath: {{ .AdmissionConfig }}\n" +
		"      readOnly: true\n" +
		"      pathType: File\n" +
		"    - name: admission-control-config-dir\n" +
		"      hostPath: /etc/kubernetes/admission-control/\n" +
		"      mountPath: /etc/kubernetes/admission-control/\n" +
		"      readOnly: true\n" +
		"      pathType: Directory\n" +
		"    {{ if and .CloudProvider .CloudConfig }}\n" +
		"    - name: cloud-config\n" +
		"      hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf\n" +
		"      mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}\n" +
		"scheduler:\n" +
		"  extraArgs:\n" +
		"    profiling: \"false\"\n" +
		"controllerManager:\n" +
		"  extraArgs:\n" +
		"    cluster-name: \"{{ .ClusterName }}\"\n" +
		"    profiling: \"false\"\n" +
		"    terminated-pod-gc-threshold: \"10\"\n" +
		"    feature-gates: \"RotateKubeletServerCertificate=true\"\n" +
		"    {{ if .ControllerManagerSigningCA }}cluster-signing-cert-file: {{ .ControllerManagerSigningCA }}{{end}}\n" +
		"    {{ if .CloudProvider }}cloud-provider: \"{{ .CloudProvider }}\"\n" +
		"    {{ if .CloudConfig }}cloud-config: /etc/kubernetes/{{ .CloudProvider }}.conf\n" +
		"  extraVolumes:\n" +
		"    - name: cloud-config\n" +
		"      hostPath: /etc/kubernetes/{{ .CloudProvider }}.conf\n" +
		"      mountPath: /etc/kubernetes/{{ .CloudProvider }}.conf{{end}}{{end}}\n" +
		"etcd:\n" +
		"  {{ if .EtcdEndpoints }}\n" +
		"  external:\n" +
		"    endpoints:\n" +
		"    {{range $k, $endpoint := .EtcdEndpoints }}  - \"{{ $endpoint }}\"\n" +
		"    {{end}}\n" +
		"    caFile: {{ .EtcdCAFile }}\n" +
		"    certFile: {{ .EtcdCertFile }}\n" +
		"    keyFile: {{ .EtcdKeyFile }}\n" +
		"  {{else}}\n" +
		"  local:\n" +
		"    extraArgs:\n" +
		"      peer-auto-tls: \"false\"\n" +
		"  {{end}}\n" +
		"---\n" +
		"apiVersion: kubelet.config.k8s.io/v1beta1\n" +
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
