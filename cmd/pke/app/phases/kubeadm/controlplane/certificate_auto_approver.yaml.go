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

// certificateAutoApproverTemplate is a generated function returning the template as a string.
func certificateAutoApproverTemplate() string {
	var tmpl = "apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: kubelet-csr-approver\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: kubelet-csr-approver\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - certificates.k8s.io\n" +
		"  resources:\n" +
		"  - certificatesigningrequests\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - certificates.k8s.io\n" +
		"  resources:\n" +
		"  - certificatesigningrequests/approval\n" +
		"  verbs:\n" +
		"  - update\n" +
		"- apiGroups:\n" +
		"  - certificates.k8s.io\n" +
		"  resourceNames:\n" +
		"  - kubernetes.io/kubelet-serving\n" +
		"  resources:\n" +
		"  - signers\n" +
		"  verbs:\n" +
		"  - approve\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: kubelet-csr-approver\n" +
		"  namespace: kube-system\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: kubelet-csr-approver\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: kubelet-csr-approver\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  name: kubelet-csr-approver\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      app: kubelet-csr-approver\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      annotations:\n" +
		"        prometheus.io/port: '8080'\n" +
		"        prometheus.io/scrape: 'true'\n" +
		"      labels:\n" +
		"        app: kubelet-csr-approver\n" +
		"    spec:\n" +
		"      serviceAccountName: kubelet-csr-approver\n" +
		"      priorityClassName: system-cluster-critical\n" +
		"      containers:\n" +
		"        - name: kubelet-csr-approver\n" +
		"          {{ if ne .ImageRepository \"banzaicloud\" }}\n" +
		"          image: \"{{ .ImageRepository }}/kubelet-csr-approver:v0.1.2\"\n" +
		"          {{ else }}\n" +
		"          image: \"postfinance/kubelet-csr-approver:v0.1.2\"\n" +
		"          {{ end }}\n" +
		"          resources:\n" +
		"            limits:\n" +
		"              memory: \"128Mi\"\n" +
		"              cpu: \"500m\"\n" +
		"          args:\n" +
		"            - -metrics-bind-address\n" +
		"            - \":8080\"\n" +
		"            - -health-probe-bind-address\n" +
		"            - \":8081\"\n" +
		"          livenessProbe:\n" +
		"            httpGet:\n" +
		"              path: /healthz\n" +
		"              port: 8081\n" +
		"          env:\n" +
		"            - name: PROVIDER_REGEX\n" +
		"              value: \\w*\n" +
		"            - name: MAX_EXPIRATION_SECONDS\n" +
		"              value: '31622400' # 366 days\n" +
		"            - name: BYPASS_DNS_RESOLUTION\n" +
		"              value: 'true'\n" +
		"      tolerations:\n" +
		"        - effect: NoSchedule\n" +
		"          key: node-role.kubernetes.io/master\n" +
		"          operator: Equal\n" +
		"        - effect: NoSchedule\n" +
		"          key: node-role.kubernetes.io/control-plane\n" +
		"          operator: Equal"
	return tmpl
}
