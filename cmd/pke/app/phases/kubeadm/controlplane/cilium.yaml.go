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

// ciliumTemplate is a generated function returning the template as a string.
func ciliumTemplate() string {
	var tmpl = "---\n" +
		"# Source: cilium/charts/config/templates/configmap.yaml\n" +
		"apiVersion: v1\n" +
		"kind: ConfigMap\n" +
		"metadata:\n" +
		"  name: cilium-config\n" +
		"  namespace: kube-system\n" +
		"data:\n" +
		"\n" +
		"  # Identity allocation mode selects how identities are shared between cilium\n" +
		"  # nodes by setting how they are stored. The options are \"crd\" or \"kvstore\".\n" +
		"  # - \"crd\" stores identities in kubernetes as CRDs (custom resource definition).\n" +
		"  #   These can be queried with:\n" +
		"  #     kubectl get ciliumid\n" +
		"  # - \"kvstore\" stores identities in a kvstore, etcd or consul, that is\n" +
		"  #   configured below. Cilium versions before 1.6 supported only the kvstore\n" +
		"  #   backend. Upgrades from these older cilium versions should continue using\n" +
		"  #   the kvstore by commenting out the identity-allocation-mode below, or\n" +
		"  #   setting it to \"kvstore\".\n" +
		"  identity-allocation-mode: crd\n" +
		"\n" +
		"  # If you want to run cilium in debug mode change this value to true\n" +
		"  debug: \"false\"\n" +
		"\n" +
		"  # Enable IPv4 addressing. If enabled, all endpoints are allocated an IPv4\n" +
		"  # address.\n" +
		"  enable-ipv4: \"true\"\n" +
		"\n" +
		"  # Enable IPv6 addressing. If enabled, all endpoints are allocated an IPv6\n" +
		"  # address.\n" +
		"  enable-ipv6: \"false\"\n" +
		"\n" +
		"  # If you want cilium monitor to aggregate tracing for packets, set this level\n" +
		"  # to \"low\", \"medium\", or \"maximum\". The higher the level, the less packets\n" +
		"  # that will be seen in monitor output.\n" +
		"  monitor-aggregation: medium\n" +
		"\n" +
		"  # ct-global-max-entries-* specifies the maximum number of connections\n" +
		"  # supported across all endpoints, split by protocol: tcp or other. One pair\n" +
		"  # of maps uses these values for IPv4 connections, and another pair of maps\n" +
		"  # use these values for IPv6 connections.\n" +
		"  #\n" +
		"  # If these values are modified, then during the next Cilium startup the\n" +
		"  # tracking of ongoing connections may be disrupted. This may lead to brief\n" +
		"  # policy drops or a change in loadbalancing decisions for a connection.\n" +
		"  #\n" +
		"  # For users upgrading from Cilium 1.2 or earlier, to minimize disruption\n" +
		"  # during the upgrade process, comment out these options.\n" +
		"  bpf-ct-global-tcp-max: \"524288\"\n" +
		"  bpf-ct-global-any-max: \"262144\"\n" +
		"\n" +
		"  # Pre-allocation of map entries allows per-packet latency to be reduced, at\n" +
		"  # the expense of up-front memory allocation for the entries in the maps. The\n" +
		"  # default value below will minimize memory usage in the default installation;\n" +
		"  # users who are sensitive to latency may consider setting this to \"true\".\n" +
		"  #\n" +
		"  # This option was introduced in Cilium 1.4. Cilium 1.3 and earlier ignore\n" +
		"  # this option and behave as though it is set to \"true\".\n" +
		"  #\n" +
		"  # If this value is modified, then during the next Cilium startup the restore\n" +
		"  # of existing endpoints and tracking of ongoing connections may be disrupted.\n" +
		"  # This may lead to policy drops or a change in loadbalancing decisions for a\n" +
		"  # connection for some time. Endpoints may need to be recreated to restore\n" +
		"  # connectivity.\n" +
		"  #\n" +
		"  # If this option is set to \"false\" during an upgrade from 1.3 or earlier to\n" +
		"  # 1.4 or later, then it may cause one-time disruptions during the upgrade.\n" +
		"  preallocate-bpf-maps: \"false\"\n" +
		"\n" +
		"  # Regular expression matching compatible Istio sidecar istio-proxy\n" +
		"  # container image names\n" +
		"  sidecar-istio-proxy-image: \"cilium/istio_proxy\"\n" +
		"\n" +
		"  # Encapsulation mode for communication between nodes\n" +
		"  # Possible values:\n" +
		"  #   - disabled\n" +
		"  #   - vxlan (default)\n" +
		"  #   - geneve\n" +
		"  tunnel: vxlan\n" +
		"\n" +
		"  # Name of the cluster. Only relevant when building a mesh of clusters.\n" +
		"  cluster-name: default\n" +
		"\n" +
		"  # DNS Polling periodically issues a DNS lookup for each `matchName` from\n" +
		"  # cilium-agent. The result is used to regenerate endpoint policy.\n" +
		"  # DNS lookups are repeated with an interval of 5 seconds, and are made for\n" +
		"  # A(IPv4) and AAAA(IPv6) addresses. Should a lookup fail, the most recent IP\n" +
		"  # data is used instead. An IP change will trigger a regeneration of the Cilium\n" +
		"  # policy for each endpoint and increment the per cilium-agent policy\n" +
		"  # repository revision.\n" +
		"  #\n" +
		"  # This option is disabled by default starting from version 1.4.x in favor\n" +
		"  # of a more powerful DNS proxy-based implementation, see [0] for details.\n" +
		"  # Enable this option if you want to use FQDN policies but do not want to use\n" +
		"  # the DNS proxy.\n" +
		"  #\n" +
		"  # To ease upgrade, users may opt to set this option to \"true\".\n" +
		"  # Otherwise please refer to the Upgrade Guide [1] which explains how to\n" +
		"  # prepare policy rules for upgrade.\n" +
		"  #\n" +
		"  # [0] http://docs.cilium.io/en/stable/policy/language/#dns-based\n" +
		"  # [1] http://docs.cilium.io/en/stable/install/upgrade/#changes-that-may-require-action\n" +
		"  tofqdns-enable-poller: \"false\"\n" +
		"\n" +
		"  # wait-bpf-mount makes init container wait until bpf filesystem is mounted\n" +
		"  wait-bpf-mount: \"false\"\n" +
		"\n" +
		"  # Enable fetching of container-runtime specific metadata\n" +
		"  #\n" +
		"  # By default, the Kubernetes pod and namespace labels are retrieved and\n" +
		"  # associated with endpoints for identification purposes. By integrating\n" +
		"  # with the container runtime, container runtime specific labels can be\n" +
		"  # retrieved, such labels will be prefixed with container:\n" +
		"  #\n" +
		"  # CAUTION: The container runtime labels can include information such as pod\n" +
		"  # annotations which may result in each pod being associated a unique set of\n" +
		"  # labels which can result in excessive security identities being allocated.\n" +
		"  # Please review the labels filter when enabling container runtime labels.\n" +
		"  #\n" +
		"  # Supported values:\n" +
		"  # - containerd\n" +
		"  # - crio\n" +
		"  # - docker\n" +
		"  # - none\n" +
		"  # - auto (automatically detect the container runtime)\n" +
		"  #\n" +
		"  container-runtime: none\n" +
		"\n" +
		"  masquerade: \"true\"\n" +
		"\n" +
		"  install-iptables-rules: \"true\"\n" +
		"  auto-direct-node-routes: \"false\"\n" +
		"  enable-node-port: \"false\"\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/serviceaccount.yaml\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: cilium\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/serviceaccount.yaml\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: cilium-operator\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/clusterrole.yaml\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: cilium\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - networking.k8s.io\n" +
		"  resources:\n" +
		"  - networkpolicies\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
		"  - namespaces\n" +
		"  - services\n" +
		"  - nodes\n" +
		"  - endpoints\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
		"  - pods\n" +
		"  - nodes\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"  - update\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
		"  - nodes\n" +
		"  - nodes/status\n" +
		"  verbs:\n" +
		"  - patch\n" +
		"- apiGroups:\n" +
		"  - extensions\n" +
		"  resources:\n" +
		"  - ingresses\n" +
		"  verbs:\n" +
		"  - create\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - apiextensions.k8s.io\n" +
		"  resources:\n" +
		"  - customresourcedefinitions\n" +
		"  verbs:\n" +
		"  - create\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"  - update\n" +
		"- apiGroups:\n" +
		"  - cilium.io\n" +
		"  resources:\n" +
		"  - ciliumnetworkpolicies\n" +
		"  - ciliumnetworkpolicies/status\n" +
		"  - ciliumendpoints\n" +
		"  - ciliumendpoints/status\n" +
		"  - ciliumnodes\n" +
		"  - ciliumnodes/status\n" +
		"  - ciliumidentities\n" +
		"  - ciliumidentities/status\n" +
		"  verbs:\n" +
		"  - '*'\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/clusterrole.yaml\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: cilium-operator\n" +
		"rules:\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
		"  # to automatically delete [core|kube]dns pods so that are starting to being\n" +
		"  # managed by Cilium\n" +
		"  - pods\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"  - delete\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
		"  # to automatically read from k8s and import the node's pod CIDR to cilium's\n" +
		"  # etcd so all nodes know how to reach another pod running in in a different\n" +
		"  # node.\n" +
		"  - nodes\n" +
		"  # to perform the translation of a CNP that contains `ToGroup` to its endpoints\n" +
		"  - services\n" +
		"  - endpoints\n" +
		"  # to check apiserver connectivity\n" +
		"  - namespaces\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - cilium.io\n" +
		"  resources:\n" +
		"  - ciliumnetworkpolicies\n" +
		"  - ciliumnetworkpolicies/status\n" +
		"  - ciliumendpoints\n" +
		"  - ciliumendpoints/status\n" +
		"  - ciliumnodes\n" +
		"  - ciliumnodes/status\n" +
		"  - ciliumidentities\n" +
		"  - ciliumidentities/status\n" +
		"  verbs:\n" +
		"  - '*'\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/clusterrolebinding.yaml\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: cilium\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: cilium\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: cilium\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/clusterrolebinding.yaml\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: cilium-operator\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: cilium-operator\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: cilium-operator\n" +
		"  namespace: kube-system\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/daemonset.yaml\n" +
		"apiVersion: apps/v1\n" +
		"kind: DaemonSet\n" +
		"metadata:\n" +
		"  labels:\n" +
		"    k8s-app: cilium\n" +
		"    kubernetes.io/cluster-service: \"true\"\n" +
		"  name: cilium\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      k8s-app: cilium\n" +
		"      kubernetes.io/cluster-service: \"true\"\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      annotations:\n" +
		"        # This annotation plus the CriticalAddonsOnly toleration makes\n" +
		"        # cilium to be a critical pod in the cluster, which ensures cilium\n" +
		"        # gets priority scheduling.\n" +
		"        # https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/\n" +
		"        scheduler.alpha.kubernetes.io/critical-pod: \"\"\n" +
		"        scheduler.alpha.kubernetes.io/tolerations: '[{\"key\":\"dedicated\",\"operator\":\"Equal\",\"value\":\"master\",\"effect\":\"NoSchedule\"}]'\n" +
		"      labels:\n" +
		"        k8s-app: cilium\n" +
		"        kubernetes.io/cluster-service: \"true\"\n" +
		"    spec:\n" +
		"      containers:\n" +
		"      - args:\n" +
		"        - --config-dir=/tmp/cilium/config-map\n" +
		"        command:\n" +
		"        - cilium-agent\n" +
		"        env:\n" +
		"        - name: K8S_NODE_NAME\n" +
		"          valueFrom:\n" +
		"            fieldRef:\n" +
		"              apiVersion: v1\n" +
		"              fieldPath: spec.nodeName\n" +
		"        - name: CILIUM_K8S_NAMESPACE\n" +
		"          valueFrom:\n" +
		"            fieldRef:\n" +
		"              apiVersion: v1\n" +
		"              fieldPath: metadata.namespace\n" +
		"        - name: CILIUM_FLANNEL_MASTER_DEVICE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: flannel-master-device\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_FLANNEL_UNINSTALL_ON_EXIT\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: flannel-uninstall-on-exit\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_CLUSTERMESH_CONFIG\n" +
		"          value: /var/lib/cilium/clustermesh/\n" +
		"        - name: CILIUM_CNI_CHAINING_MODE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: cni-chaining-mode\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_CUSTOM_CNI_CONF\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: custom-cni-conf\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        image: \"docker.io/cilium/cilium:v1.6.4\"\n" +
		"        imagePullPolicy: IfNotPresent\n" +
		"        lifecycle:\n" +
		"          postStart:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /cni-install.sh\n" +
		"          preStop:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /cni-uninstall.sh\n" +
		"        livenessProbe:\n" +
		"          exec:\n" +
		"            command:\n" +
		"            - cilium\n" +
		"            - status\n" +
		"            - --brief\n" +
		"          failureThreshold: 10\n" +
		"          # The initial delay for the liveness probe is intentionally large to\n" +
		"          # avoid an endless kill & restart cycle if in the event that the initial\n" +
		"          # bootstrapping takes longer than expected.\n" +
		"          initialDelaySeconds: 120\n" +
		"          periodSeconds: 30\n" +
		"          successThreshold: 1\n" +
		"          timeoutSeconds: 5\n" +
		"        name: cilium-agent\n" +
		"        readinessProbe:\n" +
		"          exec:\n" +
		"            command:\n" +
		"            - cilium\n" +
		"            - status\n" +
		"            - --brief\n" +
		"          failureThreshold: 3\n" +
		"          initialDelaySeconds: 5\n" +
		"          periodSeconds: 30\n" +
		"          successThreshold: 1\n" +
		"          timeoutSeconds: 5\n" +
		"        securityContext:\n" +
		"          capabilities:\n" +
		"            add:\n" +
		"            - NET_ADMIN\n" +
		"            - SYS_MODULE\n" +
		"          privileged: true\n" +
		"        volumeMounts:\n" +
		"        - mountPath: /sys/fs/bpf\n" +
		"          name: bpf-maps\n" +
		"        - mountPath: /var/run/cilium\n" +
		"          name: cilium-run\n" +
		"        - mountPath: /host/opt/cni/bin\n" +
		"          name: cni-path\n" +
		"        - mountPath: /host/etc/cni/net.d\n" +
		"          name: etc-cni-netd\n" +
		"        - mountPath: /var/lib/cilium/clustermesh\n" +
		"          name: clustermesh-secrets\n" +
		"          readOnly: true\n" +
		"        - mountPath: /tmp/cilium/config-map\n" +
		"          name: cilium-config-path\n" +
		"          readOnly: true\n" +
		"          # Needed to be able to load kernel modules\n" +
		"        - mountPath: /lib/modules\n" +
		"          name: lib-modules\n" +
		"          readOnly: true\n" +
		"        - mountPath: /run/xtables.lock\n" +
		"          name: xtables-lock\n" +
		"      hostNetwork: true\n" +
		"      initContainers:\n" +
		"      - command:\n" +
		"        - /init-container.sh\n" +
		"        env:\n" +
		"        - name: CILIUM_ALL_STATE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: clean-cilium-state\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_BPF_STATE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: clean-cilium-bpf-state\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_WAIT_BPF_MOUNT\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: wait-bpf-mount\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        image: \"docker.io/cilium/cilium:v1.6.4\"\n" +
		"        imagePullPolicy: IfNotPresent\n" +
		"        name: clean-cilium-state\n" +
		"        securityContext:\n" +
		"          capabilities:\n" +
		"            add:\n" +
		"            - NET_ADMIN\n" +
		"          privileged: true\n" +
		"        volumeMounts:\n" +
		"        - mountPath: /sys/fs/bpf\n" +
		"          name: bpf-maps\n" +
		"        - mountPath: /var/run/cilium\n" +
		"          name: cilium-run\n" +
		"      restartPolicy: Always\n" +
		"      serviceAccount: cilium\n" +
		"      serviceAccountName: cilium\n" +
		"      terminationGracePeriodSeconds: 1\n" +
		"      tolerations:\n" +
		"      - operator: Exists\n" +
		"      volumes:\n" +
		"        # To keep state between restarts / upgrades\n" +
		"      - hostPath:\n" +
		"          path: /var/run/cilium\n" +
		"          type: DirectoryOrCreate\n" +
		"        name: cilium-run\n" +
		"        # To keep state between restarts / upgrades for bpf maps\n" +
		"      - hostPath:\n" +
		"          path: /sys/fs/bpf\n" +
		"          type: DirectoryOrCreate\n" +
		"        name: bpf-maps\n" +
		"      # To install cilium cni plugin in the host\n" +
		"      - hostPath:\n" +
		"          path:  /opt/cni/bin\n" +
		"          type: DirectoryOrCreate\n" +
		"        name: cni-path\n" +
		"        # To install cilium cni configuration in the host\n" +
		"      - hostPath:\n" +
		"          path: /etc/cni/net.d\n" +
		"          type: DirectoryOrCreate\n" +
		"        name: etc-cni-netd\n" +
		"        # To be able to load kernel modules\n" +
		"      - hostPath:\n" +
		"          path: /lib/modules\n" +
		"        name: lib-modules\n" +
		"        # To access iptables concurrently with other processes (e.g. kube-proxy)\n" +
		"      - hostPath:\n" +
		"          path: /run/xtables.lock\n" +
		"          type: FileOrCreate\n" +
		"        name: xtables-lock\n" +
		"        # To read the clustermesh configuration\n" +
		"      - name: clustermesh-secrets\n" +
		"        secret:\n" +
		"          defaultMode: 420\n" +
		"          optional: true\n" +
		"          secretName: cilium-clustermesh\n" +
		"        # To read the configuration from the config map\n" +
		"      - configMap:\n" +
		"          name: cilium-config\n" +
		"        name: cilium-config-path\n" +
		"  updateStrategy:\n" +
		"    rollingUpdate:\n" +
		"      maxUnavailable: 2\n" +
		"    type: RollingUpdate\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/deployment.yaml\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  labels:\n" +
		"    io.cilium/app: operator\n" +
		"    name: cilium-operator\n" +
		"  name: cilium-operator\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  replicas: 1\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      io.cilium/app: operator\n" +
		"      name: cilium-operator\n" +
		"  strategy:\n" +
		"    rollingUpdate:\n" +
		"      maxSurge: 1\n" +
		"      maxUnavailable: 1\n" +
		"    type: RollingUpdate\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      annotations:\n" +
		"      labels:\n" +
		"        io.cilium/app: operator\n" +
		"        name: cilium-operator\n" +
		"    spec:\n" +
		"      containers:\n" +
		"      - args:\n" +
		"        - --debug=$(CILIUM_DEBUG)\n" +
		"        - --identity-allocation-mode=$(CILIUM_IDENTITY_ALLOCATION_MODE)\n" +
		"        command:\n" +
		"        - cilium-operator\n" +
		"        env:\n" +
		"        - name: CILIUM_K8S_NAMESPACE\n" +
		"          valueFrom:\n" +
		"            fieldRef:\n" +
		"              apiVersion: v1\n" +
		"              fieldPath: metadata.namespace\n" +
		"        - name: K8S_NODE_NAME\n" +
		"          valueFrom:\n" +
		"            fieldRef:\n" +
		"              apiVersion: v1\n" +
		"              fieldPath: spec.nodeName\n" +
		"        - name: CILIUM_DEBUG\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: debug\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_CLUSTER_NAME\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: cluster-name\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_CLUSTER_ID\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: cluster-id\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_IPAM\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: ipam\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_DISABLE_ENDPOINT_CRD\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: disable-endpoint-crd\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_KVSTORE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: kvstore\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: CILIUM_KVSTORE_OPT\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: kvstore-opt\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        - name: AWS_ACCESS_KEY_ID\n" +
		"          valueFrom:\n" +
		"            secretKeyRef:\n" +
		"              key: AWS_ACCESS_KEY_ID\n" +
		"              name: cilium-aws\n" +
		"              optional: true\n" +
		"        - name: AWS_SECRET_ACCESS_KEY\n" +
		"          valueFrom:\n" +
		"            secretKeyRef:\n" +
		"              key: AWS_SECRET_ACCESS_KEY\n" +
		"              name: cilium-aws\n" +
		"              optional: true\n" +
		"        - name: AWS_DEFAULT_REGION\n" +
		"          valueFrom:\n" +
		"            secretKeyRef:\n" +
		"              key: AWS_DEFAULT_REGION\n" +
		"              name: cilium-aws\n" +
		"              optional: true\n" +
		"        - name: CILIUM_IDENTITY_ALLOCATION_MODE\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: identity-allocation-mode\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        image: \"docker.io/cilium/operator:v1.6.4\"\n" +
		"        imagePullPolicy: IfNotPresent\n" +
		"        name: cilium-operator\n" +
		"        livenessProbe:\n" +
		"          httpGet:\n" +
		"            path: /healthz\n" +
		"            port: 9234\n" +
		"            scheme: HTTP\n" +
		"          initialDelaySeconds: 60\n" +
		"          periodSeconds: 10\n" +
		"          timeoutSeconds: 3\n" +
		"\n" +
		"      hostNetwork: true\n" +
		"      restartPolicy: Always\n" +
		"      serviceAccount: cilium-operator\n" +
		"      serviceAccountName: cilium-operator\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/servicemonitor.yaml\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/agent/templates/svc.yaml\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/servicemonitor.yaml\n" +
		"\n" +
		"---\n" +
		"# Source: cilium/charts/operator/templates/svc.yaml\n" +
		""
	return tmpl
}
