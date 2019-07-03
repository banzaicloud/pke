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

// podSecurityPolicyTemplate is a generated function returning the template as a string.
func podSecurityPolicyTemplate() string {
	var tmpl = "apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: RoleBinding\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:unprivileged-addon\n" +
		"  namespace: kube-system\n" +
		"  labels:\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: Role\n" +
		"  name: pke:podsecuritypolicy:unprivileged-addon\n" +
		"subjects:\n" +
		"- kind: Group\n" +
		"  # All service accounts in the kube-system namespace are allowed to use this.\n" +
		"  name: system:serviceaccounts:kube-system\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: RoleBinding\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:nodes\n" +
		"  namespace: kube-system\n" +
		"  annotations:\n" +
		"    kubernetes.io/description: 'Allow nodes to create privileged pods. Should\n" +
		"      be used in combination with the NodeRestriction admission plugin to limit\n" +
		"      nodes to mirror pods bound to themselves.'\n" +
		"  labels:\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"    kubernetes.io/cluster-service: 'true'\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: pke:podsecuritypolicy:privileged\n" +
		"subjects:\n" +
		"  - kind: Group\n" +
		"    apiGroup: rbac.authorization.k8s.io\n" +
		"    name: system:nodes\n" +
		"  - kind: User\n" +
		"    apiGroup: rbac.authorization.k8s.io\n" +
		"    # Legacy node ID\n" +
		"    name: kubelet\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"# The persistent volume binder creates recycler pods in the default namespace,\n" +
		"# but the addon manager only creates namespaced objects in the kube-system\n" +
		"# namespace, so this is a ClusterRoleBinding.\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:persistent-volume-binder\n" +
		"  labels:\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: pke:podsecuritypolicy:persistent-volume-binder\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: persistent-volume-binder\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"# The persistent volume binder creates recycler pods in the default namespace,\n" +
		"# but the addon manager only creates namespaced objects in the kube-system\n" +
		"# namespace, so this is a ClusterRole.\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:persistent-volume-binder\n" +
		"  namespace: default\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - policy\n" +
		"  resourceNames:\n" +
		"  - pke.persistent-volume-binder\n" +
		"  resources:\n" +
		"  - podsecuritypolicies\n" +
		"  verbs:\n" +
		"  - use\n" +
		"---\n" +
		"apiVersion: policy/v1beta1\n" +
		"kind: PodSecurityPolicy\n" +
		"metadata:\n" +
		"  name: pke.persistent-volume-binder\n" +
		"  annotations:\n" +
		"    kubernetes.io/description: 'Policy used by the persistent-volume-binder\n" +
		"      (a.k.a. persistentvolume-controller) to run recycler pods.'\n" +
		"    seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'\n" +
		"    seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default'\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: 'true'\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"spec:\n" +
		"  privileged: false\n" +
		"  volumes:\n" +
		"  - 'nfs'\n" +
		"  - 'secret'   # Required for service account credentials.\n" +
		"  - 'projected'\n" +
		"  hostNetwork: false\n" +
		"  hostIPC: false\n" +
		"  hostPID: false\n" +
		"  runAsUser:\n" +
		"    rule: 'RunAsAny'\n" +
		"  seLinux:\n" +
		"    rule: 'RunAsAny'\n" +
		"  supplementalGroups:\n" +
		"    rule: 'RunAsAny'\n" +
		"  fsGroup:\n" +
		"    rule: 'RunAsAny'\n" +
		"  readOnlyRootFilesystem: false\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: RoleBinding\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:privileged-binding\n" +
		"  namespace: kube-system\n" +
		"  labels:\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: pke:podsecuritypolicy:privileged\n" +
		"subjects:\n" +
		"  - kind: ServiceAccount\n" +
		"    name: kube-proxy\n" +
		"    namespace: kube-system\n" +
		"  - kind: ServiceAccount\n" +
		"    name: weave-net\n" +
		"    namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:privileged\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - policy\n" +
		"  resourceNames:\n" +
		"  - pke.privileged\n" +
		"  resources:\n" +
		"  - podsecuritypolicies\n" +
		"  verbs:\n" +
		"  - use\n" +
		"---\n" +
		"apiVersion: policy/v1beta1\n" +
		"kind: PodSecurityPolicy\n" +
		"metadata:\n" +
		"  name: pke.privileged\n" +
		"  annotations:\n" +
		"    kubernetes.io/description: 'privileged allows full unrestricted access to\n" +
		"      pod features, as if the PodSecurityPolicy controller was not enabled.'\n" +
		"    seccomp.security.alpha.kubernetes.io/allowedProfileNames: '*'\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"spec:\n" +
		"  privileged: true\n" +
		"  allowPrivilegeEscalation: true\n" +
		"  allowedCapabilities:\n" +
		"  - '*'\n" +
		"  volumes:\n" +
		"  - '*'\n" +
		"  hostNetwork: true\n" +
		"  hostPorts:\n" +
		"  - min: 0\n" +
		"    max: 65535\n" +
		"  hostIPC: true\n" +
		"  hostPID: true\n" +
		"  runAsUser:\n" +
		"    rule: 'RunAsAny'\n" +
		"  seLinux:\n" +
		"    rule: 'RunAsAny'\n" +
		"  supplementalGroups:\n" +
		"    rule: 'RunAsAny'\n" +
		"  fsGroup:\n" +
		"    rule: 'RunAsAny'\n" +
		"  readOnlyRootFilesystem: false\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: Role\n" +
		"metadata:\n" +
		"  name: pke:podsecuritypolicy:unprivileged-addon\n" +
		"  namespace: kube-system\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - policy\n" +
		"  resourceNames:\n" +
		"  - pke.unprivileged-addon\n" +
		"  resources:\n" +
		"  - podsecuritypolicies\n" +
		"  verbs:\n" +
		"  - use\n" +
		"---\n" +
		"apiVersion: policy/v1beta1\n" +
		"kind: PodSecurityPolicy\n" +
		"metadata:\n" +
		"  name: pke.unprivileged-addon\n" +
		"  annotations:\n" +
		"    kubernetes.io/description: 'This policy grants the minimum amount of\n" +
		"      privilege necessary to run non-privileged kube-system pods. This policy is\n" +
		"      not intended for use outside of kube-system, and may include further\n" +
		"      restrictions in the future.'\n" +
		"    seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'\n" +
		"    seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default'\n" +
		"  labels:\n" +
		"    kubernetes.io/cluster-service: 'true'\n" +
		"    addonmanager.kubernetes.io/mode: Reconcile\n" +
		"spec:\n" +
		"  privileged: false\n" +
		"  allowPrivilegeEscalation: false\n" +
		"  # The docker default set of capabilities\n" +
		"  allowedCapabilities:\n" +
		"  - SETPCAP\n" +
		"  - MKNOD\n" +
		"  - AUDIT_WRITE\n" +
		"  - CHOWN\n" +
		"  - NET_RAW\n" +
		"  - DAC_OVERRIDE\n" +
		"  - FOWNER\n" +
		"  - FSETID\n" +
		"  - KILL\n" +
		"  - SETGID\n" +
		"  - SETUID\n" +
		"  - NET_BIND_SERVICE\n" +
		"  - SYS_CHROOT\n" +
		"  - SETFCAP\n" +
		"  volumes:\n" +
		"  - 'emptyDir'\n" +
		"  - 'configMap'\n" +
		"  - 'secret'\n" +
		"  - 'projected'\n" +
		"  hostNetwork: false\n" +
		"  hostIPC: false\n" +
		"  hostPID: false\n" +
		"  # TODO: The addons using this profile should not run as root.\n" +
		"  runAsUser:\n" +
		"    rule: 'RunAsAny'\n" +
		"  seLinux:\n" +
		"    rule: 'RunAsAny'\n" +
		"  supplementalGroups:\n" +
		"    rule: 'RunAsAny'\n" +
		"  fsGroup:\n" +
		"    rule: 'RunAsAny'\n" +
		"  readOnlyRootFilesystem: false\n" +
		""
	return tmpl
}
