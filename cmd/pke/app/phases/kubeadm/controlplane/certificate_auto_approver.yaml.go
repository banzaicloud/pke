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
		"  name: auto-approver\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: auto-approver\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - certificates.k8s.io\n" +
		"  resources:\n" +
		"  - signers\n" +
		"  resourceNames:\n" +
		"  - \"kubernetes.io/legacy-unknown\"\n" +
		"  - \"kubernetes.io/kubelet-serving\"\n" +
		"  verbs:\n" +
		"  - approve\n" +
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
		"  - create\n" +
		"  - update\n" +
		"- apiGroups:\n" +
		"  - authorization.k8s.io\n" +
		"  resources:\n" +
		"  - subjectaccessreviews\n" +
		"  verbs:\n" +
		"  - create\n" +
		"---\n" +
		"kind: ClusterRoleBinding\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"metadata:\n" +
		"  name: auto-approver\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  namespace: kube-system\n" +
		"  name: auto-approver\n" +
		"roleRef:\n" +
		"  kind: ClusterRole\n" +
		"  name: auto-approver\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"---\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  name: auto-approver\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  replicas: 1\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      name: auto-approver\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      labels:\n" +
		"        name: auto-approver\n" +
		"    spec:\n" +
		"      serviceAccountName: auto-approver\n" +
		"      tolerations:\n" +
		"        - effect: NoSchedule\n" +
		"          operator: Exists\n" +
		"      nodeSelector:\n" +
		"        node-role.kubernetes.io/master: \"\"\n" +
		"      priorityClassName: system-cluster-critical\n" +
		"      containers:\n" +
		"        - name: auto-approver\n" +
		"          image: ghcr.io/banzaicloud/auto-approver:0.2.0\n" +
		"          args:\n" +
		"            - \"--v=2\"\n" +
		"          imagePullPolicy: Always\n" +
		"          env:\n" +
		"            - name: WATCH_NAMESPACE\n" +
		"              value: \"\"\n" +
		"            - name: POD_NAME\n" +
		"              valueFrom:\n" +
		"                fieldRef:\n" +
		"                  fieldPath: metadata.name\n" +
		"            - name: OPERATOR_NAME\n" +
		"              value: \"auto-approver\""
	return tmpl
}
