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

// auditV1Beta1Template is a generated function returning the template as a string.
func auditV1Beta1Template() string {
	var tmpl = "apiVersion: audit.k8s.io/v1beta1\n" +
		"kind: Policy\n" +
		"rules:\n" +
		"  - level: None\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - endpoints\n" +
		"          - services\n" +
		"          - services/status\n" +
		"    users:\n" +
		"      - 'system:kube-proxy'\n" +
		"      - 'system:apiserver'\n" +
		"    verbs:\n" +
		"      - watch\n" +
		"\n" +
		"  - level: None\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - nodes\n" +
		"          - nodes/status\n" +
		"    userGroups:\n" +
		"      - 'system:nodes'\n" +
		"    verbs:\n" +
		"      - get\n" +
		"\n" +
		"  - level: None\n" +
		"    namespaces:\n" +
		"      - kube-system\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - endpoints\n" +
		"    users:\n" +
		"      - 'system:kube-controller-manager'\n" +
		"      - 'system:kube-scheduler'\n" +
		"      - 'system:serviceaccount:kube-system:endpoint-controller'\n" +
		"      - 'system:serviceaccount:kube-system:local-path-provisioner-service-account'\n" +
		"      - 'system:apiserver'\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - update\n" +
		"\n" +
		"  - level: None\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - namespaces\n" +
		"          - namespaces/status\n" +
		"          - namespaces/finalize\n" +
		"    users:\n" +
		"      - 'system:apiserver'\n" +
		"    verbs:\n" +
		"      - get\n" +
		"\n" +
		"  - level: None\n" +
		"    resources:\n" +
		"      - group: metrics.k8s.io\n" +
		"    users:\n" +
		"      - 'system:kube-controller-manager'\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"\n" +
		"  - level: None\n" +
		"    nonResourceURLs:\n" +
		"      - '/healthz*'\n" +
		"      - /version\n" +
		"      - '/swagger*'\n" +
		"\n" +
		"  - level: None\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - events\n" +
		"\n" +
		"  - level: Request\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - nodes/status\n" +
		"          - pods/status\n" +
		"    users:\n" +
		"      - kubelet\n" +
		"      - 'system:node-problem-detector'\n" +
		"      - 'system:serviceaccount:kube-system:node-problem-detector'\n" +
		"    verbs:\n" +
		"      - update\n" +
		"      - patch\n" +
		"\n" +
		"  - level: Request\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - nodes/status\n" +
		"          - pods/status\n" +
		"    userGroups:\n" +
		"      - 'system:nodes'\n" +
		"    verbs:\n" +
		"      - update\n" +
		"      - patch\n" +
		"\n" +
		"  - level: Request\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    users:\n" +
		"      - 'system:serviceaccount:kube-system:namespace-controller'\n" +
		"    verbs:\n" +
		"      - deletecollection\n" +
		"\n" +
		"  - level: Metadata\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"        resources:\n" +
		"          - secrets\n" +
		"          - configmaps\n" +
		"      - group: authentication.k8s.io\n" +
		"        resources:\n" +
		"          - tokenreviews\n" +
		"\n" +
		"  - level: Request\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"      - group: admissionregistration.k8s.io\n" +
		"      - group: apiextensions.k8s.io\n" +
		"      - group: apiregistration.k8s.io\n" +
		"      - group: apps\n" +
		"      - group: authentication.k8s.io\n" +
		"      - group: authorization.k8s.io\n" +
		"      - group: autoscaling\n" +
		"      - group: batch\n" +
		"      - group: certificates.k8s.io\n" +
		"      - group: extensions\n" +
		"      - group: metrics.k8s.io\n" +
		"      - group: networking.k8s.io\n" +
		"      - group: policy\n" +
		"      - group: rbac.authorization.k8s.io\n" +
		"      - group: scheduling.k8s.io\n" +
		"      - group: settings.k8s.io\n" +
		"      - group: storage.k8s.io\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"      - watch\n" +
		"\n" +
		"  - level: RequestResponse\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		"    resources:\n" +
		"      - group: \"\"\n" +
		"      - group: admissionregistration.k8s.io\n" +
		"      - group: apiextensions.k8s.io\n" +
		"      - group: apiregistration.k8s.io\n" +
		"      - group: apps\n" +
		"      - group: authentication.k8s.io\n" +
		"      - group: authorization.k8s.io\n" +
		"      - group: autoscaling\n" +
		"      - group: batch\n" +
		"      - group: certificates.k8s.io\n" +
		"      - group: extensions\n" +
		"      - group: metrics.k8s.io\n" +
		"      - group: networking.k8s.io\n" +
		"      - group: policy\n" +
		"      - group: rbac.authorization.k8s.io\n" +
		"      - group: scheduling.k8s.io\n" +
		"      - group: settings.k8s.io\n" +
		"      - group: storage.k8s.io\n" +
		"\n" +
		"  - level: Metadata\n" +
		"    omitStages:\n" +
		"      - RequestReceived\n" +
		""
	return tmpl
}
