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

// calicoTemplate is a generated function returning the template as a string.
func calicoTemplate() string {
	var tmpl = "---\n" +
		"# Source: calico/templates/calico-config.yaml\n" +
		"# This ConfigMap is used to configure a self-hosted Calico installation.\n" +
		"kind: ConfigMap\n" +
		"apiVersion: v1\n" +
		"metadata:\n" +
		"  name: calico-config\n" +
		"  namespace: kube-system\n" +
		"data:\n" +
		"  # Typha is disabled.\n" +
		"  typha_service_name: \"none\"\n" +
		"  # Configure the backend to use.\n" +
		"  calico_backend: \"bird\"\n" +
		"\n" +
		"  # Configure the MTU to use\n" +
		"  veth_mtu: \"1440\"\n" +
		"\n" +
		"  # The CNI network configuration to install on each node.  The special\n" +
		"  # values in this config will be automatically populated.\n" +
		"  cni_network_config: |-\n" +
		"    {\n" +
		"      \"name\": \"k8s-pod-network\",\n" +
		"      \"cniVersion\": \"0.3.1\",\n" +
		"      \"plugins\": [\n" +
		"        {\n" +
		"          \"type\": \"calico\",\n" +
		"          \"log_level\": \"info\",\n" +
		"          \"datastore_type\": \"kubernetes\",\n" +
		"          \"nodename\": \"__KUBERNETES_NODE_NAME__\",\n" +
		"          \"mtu\": __CNI_MTU__,\n" +
		"          \"ipam\": {\n" +
		"              \"type\": \"calico-ipam\"\n" +
		"          },\n" +
		"          \"policy\": {\n" +
		"              \"type\": \"k8s\"\n" +
		"          },\n" +
		"          \"kubernetes\": {\n" +
		"              \"kubeconfig\": \"__KUBECONFIG_FILEPATH__\"\n" +
		"          }\n" +
		"        },\n" +
		"        {\n" +
		"          \"type\": \"portmap\",\n" +
		"          \"snat\": true,\n" +
		"          \"capabilities\": {\"portMappings\": true}\n" +
		"        }\n" +
		"      ]\n" +
		"    }\n" +
		"\n" +
		"---\n" +
		"# Source: calico/templates/kdd-crds.yaml\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"   name: felixconfigurations.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: FelixConfiguration\n" +
		"    plural: felixconfigurations\n" +
		"    singular: felixconfiguration\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: ipamblocks.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: IPAMBlock\n" +
		"    plural: ipamblocks\n" +
		"    singular: ipamblock\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: blockaffinities.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: BlockAffinity\n" +
		"    plural: blockaffinities\n" +
		"    singular: blockaffinity\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: ipamhandles.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: IPAMHandle\n" +
		"    plural: ipamhandles\n" +
		"    singular: ipamhandle\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: ipamconfigs.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: IPAMConfig\n" +
		"    plural: ipamconfigs\n" +
		"    singular: ipamconfig\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: bgppeers.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: BGPPeer\n" +
		"    plural: bgppeers\n" +
		"    singular: bgppeer\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: bgpconfigurations.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: BGPConfiguration\n" +
		"    plural: bgpconfigurations\n" +
		"    singular: bgpconfiguration\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: ippools.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: IPPool\n" +
		"    plural: ippools\n" +
		"    singular: ippool\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: hostendpoints.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: HostEndpoint\n" +
		"    plural: hostendpoints\n" +
		"    singular: hostendpoint\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: clusterinformations.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: ClusterInformation\n" +
		"    plural: clusterinformations\n" +
		"    singular: clusterinformation\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: globalnetworkpolicies.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: GlobalNetworkPolicy\n" +
		"    plural: globalnetworkpolicies\n" +
		"    singular: globalnetworkpolicy\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: globalnetworksets.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Cluster\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: GlobalNetworkSet\n" +
		"    plural: globalnetworksets\n" +
		"    singular: globalnetworkset\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: networkpolicies.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Namespaced\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: NetworkPolicy\n" +
		"    plural: networkpolicies\n" +
		"    singular: networkpolicy\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: apiextensions.k8s.io/v1beta1\n" +
		"kind: CustomResourceDefinition\n" +
		"metadata:\n" +
		"  name: networksets.crd.projectcalico.org\n" +
		"spec:\n" +
		"  scope: Namespaced\n" +
		"  group: crd.projectcalico.org\n" +
		"  version: v1\n" +
		"  names:\n" +
		"    kind: NetworkSet\n" +
		"    plural: networksets\n" +
		"    singular: networkset\n" +
		"---\n" +
		"# Source: calico/templates/rbac.yaml\n" +
		"\n" +
		"# Include a clusterrole for the kube-controllers component,\n" +
		"# and bind it to the calico-kube-controllers serviceaccount.\n" +
		"kind: ClusterRole\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"metadata:\n" +
		"  name: calico-kube-controllers\n" +
		"rules:\n" +
		"  # Nodes are watched to monitor for deletions.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - nodes\n" +
		"    verbs:\n" +
		"      - watch\n" +
		"      - list\n" +
		"      - get\n" +
		"  # Pods are queried to check for existence.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - pods\n" +
		"    verbs:\n" +
		"      - get\n" +
		"  # IPAM resources are manipulated when nodes are deleted.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - ippools\n" +
		"    verbs:\n" +
		"      - list\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - blockaffinities\n" +
		"      - ipamblocks\n" +
		"      - ipamhandles\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"      - create\n" +
		"      - update\n" +
		"      - delete\n" +
		"  # Needs access to update clusterinformations.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - clusterinformations\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - create\n" +
		"      - update\n" +
		"---\n" +
		"kind: ClusterRoleBinding\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"metadata:\n" +
		"  name: calico-kube-controllers\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: calico-kube-controllers\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: calico-kube-controllers\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"# Include a clusterrole for the calico-node DaemonSet,\n" +
		"# and bind it to the calico-node serviceaccount.\n" +
		"kind: ClusterRole\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"metadata:\n" +
		"  name: calico-node\n" +
		"rules:\n" +
		"  # The CNI plugin needs to get pods, nodes, and namespaces.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - pods\n" +
		"      - nodes\n" +
		"      - namespaces\n" +
		"    verbs:\n" +
		"      - get\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - endpoints\n" +
		"      - services\n" +
		"    verbs:\n" +
		"      # Used to discover service IPs for advertisement.\n" +
		"      - watch\n" +
		"      - list\n" +
		"      # Used to discover Typhas.\n" +
		"      - get\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - nodes/status\n" +
		"    verbs:\n" +
		"      # Needed for clearing NodeNetworkUnavailable flag.\n" +
		"      - patch\n" +
		"      # Calico stores some configuration information in node annotations.\n" +
		"      - update\n" +
		"  # Watch for changes to Kubernetes NetworkPolicies.\n" +
		"  - apiGroups: [\"networking.k8s.io\"]\n" +
		"    resources:\n" +
		"      - networkpolicies\n" +
		"    verbs:\n" +
		"      - watch\n" +
		"      - list\n" +
		"  # Used by Calico for policy information.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - pods\n" +
		"      - namespaces\n" +
		"      - serviceaccounts\n" +
		"    verbs:\n" +
		"      - list\n" +
		"      - watch\n" +
		"  # The CNI plugin patches pods/status.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - pods/status\n" +
		"    verbs:\n" +
		"      - patch\n" +
		"  # Calico monitors various CRDs for config.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - globalfelixconfigs\n" +
		"      - felixconfigurations\n" +
		"      - bgppeers\n" +
		"      - globalbgpconfigs\n" +
		"      - bgpconfigurations\n" +
		"      - ippools\n" +
		"      - ipamblocks\n" +
		"      - globalnetworkpolicies\n" +
		"      - globalnetworksets\n" +
		"      - networkpolicies\n" +
		"      - networksets\n" +
		"      - clusterinformations\n" +
		"      - hostendpoints\n" +
		"      - blockaffinities\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"      - watch\n" +
		"  # Calico must create and update some CRDs on startup.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - ippools\n" +
		"      - felixconfigurations\n" +
		"      - clusterinformations\n" +
		"    verbs:\n" +
		"      - create\n" +
		"      - update\n" +
		"  # Calico stores some configuration information on the node.\n" +
		"  - apiGroups: [\"\"]\n" +
		"    resources:\n" +
		"      - nodes\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"      - watch\n" +
		"  # These permissions are only requried for upgrade from v2.6, and can\n" +
		"  # be removed after upgrade or on fresh installations.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - bgpconfigurations\n" +
		"      - bgppeers\n" +
		"    verbs:\n" +
		"      - create\n" +
		"      - update\n" +
		"  # These permissions are required for Calico CNI to perform IPAM allocations.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - blockaffinities\n" +
		"      - ipamblocks\n" +
		"      - ipamhandles\n" +
		"    verbs:\n" +
		"      - get\n" +
		"      - list\n" +
		"      - create\n" +
		"      - update\n" +
		"      - delete\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - ipamconfigs\n" +
		"    verbs:\n" +
		"      - get\n" +
		"  # Block affinities must also be watchable by confd for route aggregation.\n" +
		"  - apiGroups: [\"crd.projectcalico.org\"]\n" +
		"    resources:\n" +
		"      - blockaffinities\n" +
		"    verbs:\n" +
		"      - watch\n" +
		"  # The Calico IPAM migration needs to get daemonsets. These permissions can be\n" +
		"  # removed if not upgrading from an installation using host-local IPAM.\n" +
		"  - apiGroups: [\"apps\"]\n" +
		"    resources:\n" +
		"      - daemonsets\n" +
		"    verbs:\n" +
		"      - get\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: calico-node\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: calico-node\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: calico-node\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: calico/templates/calico-node.yaml\n" +
		"# This manifest installs the calico-node container, as well\n" +
		"# as the CNI plugins and network config on\n" +
		"# each master and worker node in a Kubernetes cluster.\n" +
		"kind: DaemonSet\n" +
		"apiVersion: apps/v1\n" +
		"metadata:\n" +
		"  name: calico-node\n" +
		"  namespace: kube-system\n" +
		"  labels:\n" +
		"    k8s-app: calico-node\n" +
		"spec:\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      k8s-app: calico-node\n" +
		"  updateStrategy:\n" +
		"    type: RollingUpdate\n" +
		"    rollingUpdate:\n" +
		"      maxUnavailable: 1\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      labels:\n" +
		"        k8s-app: calico-node\n" +
		"      annotations:\n" +
		"        # This, along with the CriticalAddonsOnly toleration below,\n" +
		"        # marks the pod as a critical add-on, ensuring it gets\n" +
		"        # priority scheduling and that its resources are reserved\n" +
		"        # if it ever gets evicted.\n" +
		"        scheduler.alpha.kubernetes.io/critical-pod: ''\n" +
		"    spec:\n" +
		"      nodeSelector:\n" +
		"        beta.kubernetes.io/os: linux\n" +
		"      hostNetwork: true\n" +
		"      tolerations:\n" +
		"        # Make sure calico-node gets scheduled on all nodes.\n" +
		"        - effect: NoSchedule\n" +
		"          operator: Exists\n" +
		"        # Mark the pod as a critical add-on for rescheduling.\n" +
		"        - key: CriticalAddonsOnly\n" +
		"          operator: Exists\n" +
		"        - effect: NoExecute\n" +
		"          operator: Exists\n" +
		"      serviceAccountName: calico-node\n" +
		"      # Minimize downtime during a rolling upgrade or deletion; tell Kubernetes to do a \"force\n" +
		"      # deletion\": https://kubernetes.io/docs/concepts/workloads/pods/pod/#termination-of-pods.\n" +
		"      terminationGracePeriodSeconds: 0\n" +
		"      priorityClassName: system-node-critical\n" +
		"      initContainers:\n" +
		"        # This container performs upgrade from host-local IPAM to calico-ipam.\n" +
		"        # It can be deleted if this is a fresh installation, or if you have already\n" +
		"        # upgraded to use calico-ipam.\n" +
		"        - name: upgrade-ipam\n" +
		"          image: calico/cni:v3.10.1\n" +
		"          command: [\"/opt/cni/bin/calico-ipam\", \"-upgrade\"]\n" +
		"          env:\n" +
		"            - name: KUBERNETES_NODE_NAME\n" +
		"              valueFrom:\n" +
		"                fieldRef:\n" +
		"                  fieldPath: spec.nodeName\n" +
		"            - name: CALICO_NETWORKING_BACKEND\n" +
		"              valueFrom:\n" +
		"                configMapKeyRef:\n" +
		"                  name: calico-config\n" +
		"                  key: calico_backend\n" +
		"          volumeMounts:\n" +
		"            - mountPath: /var/lib/cni/networks\n" +
		"              name: host-local-net-dir\n" +
		"            - mountPath: /host/opt/cni/bin\n" +
		"              name: cni-bin-dir\n" +
		"        # This container installs the CNI binaries\n" +
		"        # and CNI network config file on each node.\n" +
		"        - name: install-cni\n" +
		"          image: calico/cni:v3.10.1\n" +
		"          command: [\"/install-cni.sh\"]\n" +
		"          env:\n" +
		"            # Name of the CNI config file to create.\n" +
		"            - name: CNI_CONF_NAME\n" +
		"              value: \"10-calico.conflist\"\n" +
		"            # The CNI network config to install on each node.\n" +
		"            - name: CNI_NETWORK_CONFIG\n" +
		"              valueFrom:\n" +
		"                configMapKeyRef:\n" +
		"                  name: calico-config\n" +
		"                  key: cni_network_config\n" +
		"            # Set the hostname based on the k8s node name.\n" +
		"            - name: KUBERNETES_NODE_NAME\n" +
		"              valueFrom:\n" +
		"                fieldRef:\n" +
		"                  fieldPath: spec.nodeName\n" +
		"            # CNI MTU Config variable\n" +
		"            - name: CNI_MTU\n" +
		"              valueFrom:\n" +
		"                configMapKeyRef:\n" +
		"                  name: calico-config\n" +
		"                  key: veth_mtu\n" +
		"            # Prevents the container from sleeping forever.\n" +
		"            - name: SLEEP\n" +
		"              value: \"false\"\n" +
		"          volumeMounts:\n" +
		"            - mountPath: /host/opt/cni/bin\n" +
		"              name: cni-bin-dir\n" +
		"            - mountPath: /host/etc/cni/net.d\n" +
		"              name: cni-net-dir\n" +
		"        # Adds a Flex Volume Driver that creates a per-pod Unix Domain Socket to allow Dikastes\n" +
		"        # to communicate with Felix over the Policy Sync API.\n" +
		"        - name: flexvol-driver\n" +
		"          image: calico/pod2daemon-flexvol:v3.10.1\n" +
		"          volumeMounts:\n" +
		"          - name: flexvol-driver-host\n" +
		"            mountPath: /host/driver\n" +
		"      containers:\n" +
		"        # Runs calico-node container on each Kubernetes node.  This\n" +
		"        # container programs network policy and routes on each\n" +
		"        # host.\n" +
		"        - name: calico-node\n" +
		"          image: calico/node:v3.10.1\n" +
		"          env:\n" +
		"            # Use Kubernetes API as the backing datastore.\n" +
		"            - name: DATASTORE_TYPE\n" +
		"              value: \"kubernetes\"\n" +
		"            # Wait for the datastore.\n" +
		"            - name: WAIT_FOR_DATASTORE\n" +
		"              value: \"true\"\n" +
		"            # Set based on the k8s node name.\n" +
		"            - name: NODENAME\n" +
		"              valueFrom:\n" +
		"                fieldRef:\n" +
		"                  fieldPath: spec.nodeName\n" +
		"            # Choose the backend to use.\n" +
		"            - name: CALICO_NETWORKING_BACKEND\n" +
		"              valueFrom:\n" +
		"                configMapKeyRef:\n" +
		"                  name: calico-config\n" +
		"                  key: calico_backend\n" +
		"            # Cluster type to identify the deployment type\n" +
		"            - name: CLUSTER_TYPE\n" +
		"              value: \"k8s,bgp\"\n" +
		"            # Auto-detect the BGP IP address.\n" +
		"            - name: IP\n" +
		"              value: \"autodetect\"\n" +
		"            # Enable IPIP\n" +
		"            - name: CALICO_IPV4POOL_IPIP\n" +
		"              value: \"Always\"\n" +
		"            # Set MTU for tunnel device used if ipip is enabled\n" +
		"            - name: FELIX_IPINIPMTU\n" +
		"              valueFrom:\n" +
		"                configMapKeyRef:\n" +
		"                  name: calico-config\n" +
		"                  key: veth_mtu\n" +
		"            # The default IPv4 pool to create on startup if none exists. Pod IPs will be\n" +
		"            # chosen from this range. Changing this value after installation will have\n" +
		"            # no effect. This should fall within `--cluster-cidr`.\n" +
		"            - name: CALICO_IPV4POOL_CIDR\n" +
		"              value: \"192.168.0.0/16\"\n" +
		"            # Disable file logging so `kubectl logs` works.\n" +
		"            - name: CALICO_DISABLE_FILE_LOGGING\n" +
		"              value: \"true\"\n" +
		"            # Set Felix endpoint to host default action to ACCEPT.\n" +
		"            - name: FELIX_DEFAULTENDPOINTTOHOSTACTION\n" +
		"              value: \"ACCEPT\"\n" +
		"            # Disable IPv6 on Kubernetes.\n" +
		"            - name: FELIX_IPV6SUPPORT\n" +
		"              value: \"false\"\n" +
		"            # Set Felix logging to \"info\"\n" +
		"            - name: FELIX_LOGSEVERITYSCREEN\n" +
		"              value: \"info\"\n" +
		"            - name: FELIX_HEALTHENABLED\n" +
		"              value: \"true\"\n" +
		"          securityContext:\n" +
		"            privileged: true\n" +
		"          resources:\n" +
		"            requests:\n" +
		"              cpu: 250m\n" +
		"          livenessProbe:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /bin/calico-node\n" +
		"              - -felix-live\n" +
		"            periodSeconds: 10\n" +
		"            initialDelaySeconds: 10\n" +
		"            failureThreshold: 6\n" +
		"          readinessProbe:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /bin/calico-node\n" +
		"              - -felix-ready\n" +
		"              - -bird-ready\n" +
		"            periodSeconds: 10\n" +
		"          volumeMounts:\n" +
		"            - mountPath: /lib/modules\n" +
		"              name: lib-modules\n" +
		"              readOnly: true\n" +
		"            - mountPath: /run/xtables.lock\n" +
		"              name: xtables-lock\n" +
		"              readOnly: false\n" +
		"            - mountPath: /var/run/calico\n" +
		"              name: var-run-calico\n" +
		"              readOnly: false\n" +
		"            - mountPath: /var/lib/calico\n" +
		"              name: var-lib-calico\n" +
		"              readOnly: false\n" +
		"            - name: policysync\n" +
		"              mountPath: /var/run/nodeagent\n" +
		"      volumes:\n" +
		"        # Used by calico-node.\n" +
		"        - name: lib-modules\n" +
		"          hostPath:\n" +
		"            path: /lib/modules\n" +
		"        - name: var-run-calico\n" +
		"          hostPath:\n" +
		"            path: /var/run/calico\n" +
		"        - name: var-lib-calico\n" +
		"          hostPath:\n" +
		"            path: /var/lib/calico\n" +
		"        - name: xtables-lock\n" +
		"          hostPath:\n" +
		"            path: /run/xtables.lock\n" +
		"            type: FileOrCreate\n" +
		"        # Used to install CNI.\n" +
		"        - name: cni-bin-dir\n" +
		"          hostPath:\n" +
		"            path: /opt/cni/bin\n" +
		"        - name: cni-net-dir\n" +
		"          hostPath:\n" +
		"            path: /etc/cni/net.d\n" +
		"        # Mount in the directory for host-local IPAM allocations. This is\n" +
		"        # used when upgrading from host-local to calico-ipam, and can be removed\n" +
		"        # if not using the upgrade-ipam init container.\n" +
		"        - name: host-local-net-dir\n" +
		"          hostPath:\n" +
		"            path: /var/lib/cni/networks\n" +
		"        # Used to create per-pod Unix Domain Sockets\n" +
		"        - name: policysync\n" +
		"          hostPath:\n" +
		"            type: DirectoryOrCreate\n" +
		"            path: /var/run/nodeagent\n" +
		"        # Used to install Flex Volume Driver\n" +
		"        - name: flexvol-driver-host\n" +
		"          hostPath:\n" +
		"            type: DirectoryOrCreate\n" +
		"            path: /usr/libexec/kubernetes/kubelet-plugins/volume/exec/nodeagent~uds\n" +
		"---\n" +
		"\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: calico-node\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: calico/templates/calico-kube-controllers.yaml\n" +
		"\n" +
		"# See https://github.com/projectcalico/kube-controllers\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  name: calico-kube-controllers\n" +
		"  namespace: kube-system\n" +
		"  labels:\n" +
		"    k8s-app: calico-kube-controllers\n" +
		"spec:\n" +
		"  # The controllers can only have a single active instance.\n" +
		"  replicas: 1\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      k8s-app: calico-kube-controllers\n" +
		"  strategy:\n" +
		"    type: Recreate\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      name: calico-kube-controllers\n" +
		"      namespace: kube-system\n" +
		"      labels:\n" +
		"        k8s-app: calico-kube-controllers\n" +
		"      annotations:\n" +
		"        scheduler.alpha.kubernetes.io/critical-pod: ''\n" +
		"    spec:\n" +
		"      nodeSelector:\n" +
		"        beta.kubernetes.io/os: linux\n" +
		"      tolerations:\n" +
		"        # Mark the pod as a critical add-on for rescheduling.\n" +
		"        - key: CriticalAddonsOnly\n" +
		"          operator: Exists\n" +
		"        - key: node-role.kubernetes.io/master\n" +
		"          effect: NoSchedule\n" +
		"      serviceAccountName: calico-kube-controllers\n" +
		"      priorityClassName: system-cluster-critical\n" +
		"      containers:\n" +
		"        - name: calico-kube-controllers\n" +
		"          image: calico/kube-controllers:v3.10.1\n" +
		"          env:\n" +
		"            # Choose which controllers to run.\n" +
		"            - name: ENABLED_CONTROLLERS\n" +
		"              value: node\n" +
		"            - name: DATASTORE_TYPE\n" +
		"              value: kubernetes\n" +
		"          readinessProbe:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /usr/bin/check-status\n" +
		"              - -r\n" +
		"\n" +
		"---\n" +
		"\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: calico-kube-controllers\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"# Source: calico/templates/calico-etcd-secrets.yaml\n" +
		"\n" +
		"---\n" +
		"# Source: calico/templates/calico-typha.yaml\n" +
		"\n" +
		"---\n" +
		"# Source: calico/templates/configure-canal.yaml\n" +
		""
	return tmpl
}
