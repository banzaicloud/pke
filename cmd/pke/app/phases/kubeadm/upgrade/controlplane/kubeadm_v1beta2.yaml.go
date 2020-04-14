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

// kubeadmConfigV1Beta2Template is a generated function returning the template as a string.
func kubeadmConfigV1Beta2Template() string {
	var tmpl = "apiVersion: kubeadm.k8s.io/v1beta2\n" +
		"kind: InitConfiguration\n" +
		"{{ if .APIServerAdvertiseAddress }}\n" +
		"localAPIEndpoint:\n" +
		"  advertiseAddress: \"{{ .APIServerAdvertiseAddress }}\"\n" +
		"  bindPort: {{ .APIServerBindPort }}{{end}}\n" +
		"---\n" +
		"apiVersion: kubeadm.k8s.io/v1beta2\n" +
		"kind: ClusterConfiguration\n" +
		"clusterName: \"{{ .KubeadmConfig.ClusterName }}\"\n" +
		"imageRepository: {{ .KubeadmConfig.ImageRepository }}\n" +
		"{{ if .KubeadmConfig.UseHyperKubeImage }}useHyperKubeImage: true{{end}}\n" +
		"kubernetesVersion: \"{{ .KubeadmConfig.KubernetesVersion }}\"\n" +
		"networking:\n" +
		"  serviceSubnet: \"{{ .KubeadmConfig.Networking.ServiceSubnet }}\"\n" +
		"  podSubnet: \"{{ .KubeadmConfig.Networking.PodSubnet }}\"\n" +
		"  dnsDomain: \"cluster.local\"\n" +
		"{{ if .KubeadmConfig.ControlPlaneEndpoint }}controlPlaneEndpoint: \"{{ .KubeadmConfig.ControlPlaneEndpoint }}\"{{end}}\n" +
		"certificatesDir: \"/etc/kubernetes/pki\"\n" +
		"apiServer:\n" +
		"  {{ if .KubeadmConfig.APIServer.CertSANs }}\n" +
		"  certSANs:\n" +
		"  {{range $k, $san := .KubeadmConfig.APIServer.CertSANs}}  - \"{{ $san }}\"\n" +
		"  {{end}}{{end}}\n" +
		"  extraArgs:\n" +
		"    profiling: \"false\"\n" +
		"    enable-admission-plugins: \"{{ .KubeadmConfig.APIServer.ExtraArgs.EnableAdmissionPlugins }}\"\n" +
		"    disable-admission-plugins: \"{{ .KubeadmConfig.APIServer.ExtraArgs.DisableAdmissionPlugins }}\"\n" +
		"    admission-control-config-file: \"{{ .KubeadmConfig.APIServer.ExtraArgs.AdmissionControlConfigFile }}\"\n" +
		"    audit-log-path: \"{{ .KubeadmConfig.APIServer.ExtraArgs.AuditLogPath }}\"\n" +
		"    audit-log-maxage: \"30\"\n" +
		"    audit-log-maxbackup: \"10\"\n" +
		"    audit-log-maxsize: \"100\"\n" +
		"    {{ if .KubeadmConfig.APIServer.ExtraArgs.AuditPolicyFile }}audit-policy-file: \"{{ .KubeadmConfig.APIServer.ExtraArgs.AuditPolicyFile }}\"{{ end }}\n" +
		"    {{ if .KubeadmConfig.APIServer.ExtraArgs.EtcdPrefix }}etcd-prefix: \"{{ .KubeadmConfig.APIServer.ExtraArgs.EtcdPrefix }}\"{{end}}\n" +
		"    service-account-lookup: \"true\"\n" +
		"    kubelet-certificate-authority: \"{{ .KubeadmConfig.APIServer.ExtraArgs.KubeletCertificateAuthority }}\"\n" +
		"    tls-cipher-suites: \"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256\"\n" +
		"    encryption-provider-config: \"/etc/kubernetes/admission-control/encryption-provider-config.yaml\"\n" +
		"    {{ if (and .KubeadmConfig.APIServer.ExtraArgs.OIDCIssuerURL .KubeadmConfig.APIServer.ExtraArgs.OIDCClientID) }}\n" +
		"    oidc-issuer-url: \"{{ .KubeadmConfig.APIServer.ExtraArgs.OIDCIssuerURL }}\"\n" +
		"    oidc-client-id: \"{{ .KubeadmConfig.APIServer.ExtraArgs.OIDCClientID }}\"\n" +
		"    oidc-username-claim: \"email\"\n" +
		"    oidc-username-prefix: \"oidc:\"\n" +
		"    oidc-groups-claim: \"groups\"{{end}}\n" +
		"    {{ if .KubeadmConfig.APIServer.ExtraArgs.CloudProvider }}cloud-provider: \"{{ .KubeadmConfig.APIServer.ExtraArgs.CloudProvider }}\"\n" +
		"    {{ if .KubeadmConfig.APIServer.ExtraArgs.CloudConfig }}cloud-config: {{ .KubeadmConfig.APIServer.ExtraArgs.CloudConfig }}{{end}}{{end}}\n" +
		"  extraVolumes:\n" +
		"  {{range $k, $volume := .KubeadmConfig.APIServer.ExtraVolumes }}\n" +
		"    - name: {{ $volume.Name }}\n" +
		"      hostPath: {{ $volume.HostPath }}\n" +
		"      mountPath: {{ $volume.MountPath }}\n" +
		"      pathType: {{ $volume.PathType }}\n" +
		"      readOnly: {{ $volume.ReadOnly }}{{end}}\n" +
		"scheduler:\n" +
		"  extraArgs:\n" +
		"    profiling: \"false\"\n" +
		"controllerManager:\n" +
		"  extraArgs:\n" +
		"    cluster-name: \"{{ .KubeadmConfig.ControllerManager.ExtraArgs.ClusterName }}\"\n" +
		"    profiling: \"false\"\n" +
		"    terminated-pod-gc-threshold: \"10\"\n" +
		"    feature-gates: \"RotateKubeletServerCertificate=true\"\n" +
		"    {{ if .KubeadmConfig.ControllerManager.ExtraArgs.ClusterSigningCertFile }}cluster-signing-cert-file: {{ .KubeadmConfig.ControllerManager.ExtraArgs.ClusterSigningCertFile }}{{end}}\n" +
		"    {{ if .KubeadmConfig.ControllerManager.ExtraArgs.CloudProvider }}cloud-provider: \"{{ .KubeadmConfig.ControllerManager.ExtraArgs.CloudProvider }}\n" +
		"    {{ if .KubeadmConfig.ControllerManager.ExtraArgs.CloudConfig }}cloud-config: \"{{ .KubeadmConfig.ControllerManager.ExtraArgs.CloudConfig }}\n" +
		"  extraVolumes:\n" +
		"  {{range $k, $volume := .KubeadmConfig.ControllerManager.ExtraVolumes }}\n" +
		"    - name: {{ $volume.Name }}\n" +
		"\t    hostPath: {{ $volume.HostPath }}\n" +
		"\t    mountPath: {{ $volume.MountPath }}\n" +
		"\t    pathType: {{ $volume.PathType }}\n" +
		"\t    readOnly: {{ $volume.ReadOnly }}{{end}}{{end}}{{end}}\n" +
		"etcd:\n" +
		"  {{ if .KubeadmConfig.Etcd.External.Endpoints }}\n" +
		"  external:\n" +
		"    endpoints:\n" +
		"    {{range $k, $endpoint := .KubeadmConfig.Etcd.External.Endpoints }}\n" +
		"      - caFile: {{ $endpoint.CAFile }}\n" +
		"        certFile: {{ $endpoint.CertFile }}\n" +
		"        keyFile: {{ $endpoint.KeyFile }}{{end}}\n" +
		"  {{else}}\n" +
		"  local:\n" +
		"    extraArgs:\n" +
		"      peer-auto-tls: \"false\"\n" +
		"  {{end}}\n" +
		""
	return tmpl
}
