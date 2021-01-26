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
		"# Source: cilium/templates/cilium-agent-serviceaccount.yaml\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: cilium\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"# Source: cilium/templates/cilium-operator-serviceaccount.yaml\n" +
		"apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: cilium-operator\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"# Source: cilium/templates/cilium-configmap.yaml\n" +
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
		"  cilium-endpoint-gc-interval: \"5m0s\"\n" +
		"\n" +
		"  # If you want to run cilium in debug mode change this value to true\n" +
		"  debug: \"false\"\n" +
		"  # The agent can be put into the following three policy enforcement modes\n" +
		"  # default, always and never.\n" +
		"  # https://docs.cilium.io/en/latest/policy/intro/#policy-enforcement-modes\n" +
		"  enable-policy: \"default\"\n" +
		"\n" +
		"  # Enable IPv4 addressing. If enabled, all endpoints are allocated an IPv4\n" +
		"  # address.\n" +
		"  enable-ipv4: \"true\"\n" +
		"\n" +
		"  # Enable IPv6 addressing. If enabled, all endpoints are allocated an IPv6\n" +
		"  # address.\n" +
		"  enable-ipv6: \"false\"\n" +
		"  # Users who wish to specify their own custom CNI configuration file must set\n" +
		"  # custom-cni-conf to \"true\", otherwise Cilium may overwrite the configuration.\n" +
		"  custom-cni-conf: \"false\"\n" +
		"  enable-bpf-clock-probe: \"true\"\n" +
		"  # If you want cilium monitor to aggregate tracing for packets, set this level\n" +
		"  # to \"low\", \"medium\", or \"maximum\". The higher the level, the less packets\n" +
		"  # that will be seen in monitor output.\n" +
		"  monitor-aggregation: medium\n" +
		"\n" +
		"  # The monitor aggregation interval governs the typical time between monitor\n" +
		"  # notification events for each allowed connection.\n" +
		"  #\n" +
		"  # Only effective when monitor aggregation is set to \"medium\" or higher.\n" +
		"  monitor-aggregation-interval: 5s\n" +
		"\n" +
		"  # The monitor aggregation flags determine which TCP flags which, upon the\n" +
		"  # first observation, cause monitor notifications to be generated.\n" +
		"  #\n" +
		"  # Only effective when monitor aggregation is set to \"medium\" or higher.\n" +
		"  monitor-aggregation-flags: all\n" +
		"  # Specifies the ratio (0.0-1.0) of total system memory to use for dynamic\n" +
		"  # sizing of the TCP CT, non-TCP CT, NAT and policy BPF maps.\n" +
		"  bpf-map-dynamic-size-ratio: \"0.0025\"\n" +
		"  # bpf-policy-map-max specifies the maximum number of entries in endpoint\n" +
		"  # policy map (per endpoint)\n" +
		"  bpf-policy-map-max: \"16384\"\n" +
		"  # bpf-lb-map-max specifies the maximum number of entries in bpf lb service,\n" +
		"  # backend and affinity maps.\n" +
		"  bpf-lb-map-max: \"65536\"\n" +
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
		"  # As a result, reply packets may be dropped and the load-balancing decisions\n" +
		"  # for established connections may change.\n" +
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
		"  # Enables L7 proxy for L7 policy enforcement and visibility\n" +
		"  enable-l7-proxy: \"true\"\n" +
		"\n" +
		"  # wait-bpf-mount makes init container wait until bpf filesystem is mounted\n" +
		"  wait-bpf-mount: \"false\"\n" +
		"\n" +
		"  masquerade: \"true\"\n" +
		"  enable-bpf-masquerade: \"true\"\n" +
		"\n" +
		"  enable-xt-socket-fallback: \"true\"\n" +
		"  install-iptables-rules: \"true\"\n" +
		"\n" +
		"  auto-direct-node-routes: \"false\"\n" +
		"  enable-bandwidth-manager: \"false\"\n" +
		"  enable-local-redirect-policy: \"false\"\n" +
		"  kube-proxy-replacement:  \"disabled\"\n" +
		"  kube-proxy-replacement-healthz-bind-address: \"\"\n" +
		"  enable-health-check-nodeport: \"true\"\n" +
		"  node-port-bind-protection: \"true\"\n" +
		"  enable-auto-protect-node-port-range: \"true\"\n" +
		"  enable-session-affinity: \"true\"\n" +
		"  enable-endpoint-health-checking: \"true\"\n" +
		"  enable-health-checking: \"true\"\n" +
		"  enable-well-known-identities: \"false\"\n" +
		"  enable-remote-node-identity: \"true\"\n" +
		"  operator-api-serve-addr: \"127.0.0.1:9234\"\n" +
		"  # Enable Hubble gRPC service.\n" +
		"  enable-hubble: \"true\"\n" +
		"  # UNIX domain socket for Hubble server to listen to.\n" +
		"  hubble-socket-path:  \"/var/run/cilium/hubble.sock\"\n" +
		"  ipam: \"cluster-pool\"\n" +
		"  cluster-pool-ipv4-cidr: \"{{ .PodCIDR }}\"\n" +
		"  cluster-pool-ipv4-mask-size: \"24\"\n" +
		"  disable-cnp-status-updates: \"true\"\n" +
		"---\n" +
		"# Source: cilium/templates/cilium-agent-clusterrole.yaml\n" +
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
		"  - discovery.k8s.io\n" +
		"  resources:\n" +
		"  - endpointslices\n" +
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
		"  - pods/finalizers\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"  - update\n" +
		"  - delete\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
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
		"  - apiextensions.k8s.io\n" +
		"  resources:\n" +
		"  - customresourcedefinitions\n" +
		"  verbs:\n" +
		"  # Deprecated for removal in v1.10\n" +
		"  - create\n" +
		"  - list\n" +
		"  - watch\n" +
		"  - update\n" +
		"\n" +
		"  # This is used when validating policies in preflight. This will need to stay\n" +
		"  # until we figure out how to avoid \"get\" inside the preflight, and then\n" +
		"  # should be removed ideally.\n" +
		"  - get\n" +
		"- apiGroups:\n" +
		"  - cilium.io\n" +
		"  resources:\n" +
		"  - ciliumnetworkpolicies\n" +
		"  - ciliumnetworkpolicies/status\n" +
		"  - ciliumnetworkpolicies/finalizers\n" +
		"  - ciliumclusterwidenetworkpolicies\n" +
		"  - ciliumclusterwidenetworkpolicies/status\n" +
		"  - ciliumclusterwidenetworkpolicies/finalizers\n" +
		"  - ciliumendpoints\n" +
		"  - ciliumendpoints/status\n" +
		"  - ciliumendpoints/finalizers\n" +
		"  - ciliumnodes\n" +
		"  - ciliumnodes/status\n" +
		"  - ciliumnodes/finalizers\n" +
		"  - ciliumidentities\n" +
		"  - ciliumidentities/finalizers\n" +
		"  - ciliumlocalredirectpolicies\n" +
		"  - ciliumlocalredirectpolicies/status\n" +
		"  - ciliumlocalredirectpolicies/finalizers\n" +
		"  verbs:\n" +
		"  - '*'\n" +
		"---\n" +
		"# Source: cilium/templates/cilium-operator-clusterrole.yaml\n" +
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
		"  - discovery.k8s.io\n" +
		"  resources:\n" +
		"  - endpointslices\n" +
		"  verbs:\n" +
		"  - get\n" +
		"  - list\n" +
		"  - watch\n" +
		"- apiGroups:\n" +
		"  - \"\"\n" +
		"  resources:\n" +
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
		"  - ciliumnetworkpolicies/finalizers\n" +
		"  - ciliumclusterwidenetworkpolicies\n" +
		"  - ciliumclusterwidenetworkpolicies/status\n" +
		"  - ciliumclusterwidenetworkpolicies/finalizers\n" +
		"  - ciliumendpoints\n" +
		"  - ciliumendpoints/status\n" +
		"  - ciliumendpoints/finalizers\n" +
		"  - ciliumnodes\n" +
		"  - ciliumnodes/status\n" +
		"  - ciliumnodes/finalizers\n" +
		"  - ciliumidentities\n" +
		"  - ciliumidentities/status\n" +
		"  - ciliumidentities/finalizers\n" +
		"  - ciliumlocalredirectpolicies\n" +
		"  - ciliumlocalredirectpolicies/status\n" +
		"  - ciliumlocalredirectpolicies/finalizers\n" +
		"  verbs:\n" +
		"  - '*'\n" +
		"- apiGroups:\n" +
		"  - apiextensions.k8s.io\n" +
		"  resources:\n" +
		"  - customresourcedefinitions\n" +
		"  verbs:\n" +
		"  - create\n" +
		"  - get\n" +
		"  - list\n" +
		"  - update\n" +
		"  - watch\n" +
		"# For cilium-operator running in HA mode.\n" +
		"#\n" +
		"# Cilium operator running in HA mode requires the use of ResourceLock for Leader Election\n" +
		"# between mulitple running instances.\n" +
		"# The preferred way of doing this is to use LeasesResourceLock as edits to Leases are less\n" +
		"# common and fewer objects in the cluster watch \"all Leases\".\n" +
		"# The support for leases was introduced in coordination.k8s.io/v1 during Kubernetes 1.14 release.\n" +
		"# In Cilium we currently don't support HA mode for K8s version < 1.14. This condition make sure\n" +
		"# that we only authorize access to leases resources in supported K8s versions.\n" +
		"- apiGroups:\n" +
		"  - coordination.k8s.io\n" +
		"  resources:\n" +
		"  - leases\n" +
		"  verbs:\n" +
		"  - create\n" +
		"  - get\n" +
		"  - update\n" +
		"---\n" +
		"# Source: cilium/templates/cilium-agent-clusterrolebinding.yaml\n" +
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
		"---\n" +
		"# Source: cilium/templates/cilium-operator-clusterrolebinding.yaml\n" +
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
		"---\n" +
		"# Source: cilium/templates/cilium-agent-daemonset.yaml\n" +
		"apiVersion: apps/v1\n" +
		"kind: DaemonSet\n" +
		"metadata:\n" +
		"  labels:\n" +
		"    k8s-app: cilium\n" +
		"  name: cilium\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      k8s-app: cilium\n" +
		"  updateStrategy:\n" +
		"    rollingUpdate:\n" +
		"      maxUnavailable: 2\n" +
		"    type: RollingUpdate\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      annotations:\n" +
		"        # This annotation plus the CriticalAddonsOnly toleration makes\n" +
		"        # cilium to be a critical pod in the cluster, which ensures cilium\n" +
		"        # gets priority scheduling.\n" +
		"        # https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/\n" +
		"        scheduler.alpha.kubernetes.io/critical-pod: \"\"\n" +
		"      labels:\n" +
		"        k8s-app: cilium\n" +
		"    spec:\n" +
		"      affinity:\n" +
		"        podAntiAffinity:\n" +
		"          requiredDuringSchedulingIgnoredDuringExecution:\n" +
		"          - labelSelector:\n" +
		"              matchExpressions:\n" +
		"              - key: k8s-app\n" +
		"                operator: In\n" +
		"                values:\n" +
		"                - cilium\n" +
		"            topologyKey: kubernetes.io/hostname\n" +
		"      containers:\n" +
		"      - args:\n" +
		"        - --config-dir=/tmp/cilium/config-map\n" +
		"        command:\n" +
		"        - cilium-agent\n" +
		"        livenessProbe:\n" +
		"          httpGet:\n" +
		"            host: '127.0.0.1'\n" +
		"            path: /healthz\n" +
		"            port: 9876\n" +
		"            scheme: HTTP\n" +
		"            httpHeaders:\n" +
		"            - name: \"brief\"\n" +
		"              value: \"true\"\n" +
		"          failureThreshold: 10\n" +
		"          # The initial delay for the liveness probe is intentionally large to\n" +
		"          # avoid an endless kill & restart cycle if in the event that the initial\n" +
		"          # bootstrapping takes longer than expected.\n" +
		"          initialDelaySeconds: 120\n" +
		"          periodSeconds: 30\n" +
		"          successThreshold: 1\n" +
		"          timeoutSeconds: 5\n" +
		"        readinessProbe:\n" +
		"          httpGet:\n" +
		"            host: '127.0.0.1'\n" +
		"            path: /healthz\n" +
		"            port: 9876\n" +
		"            scheme: HTTP\n" +
		"            httpHeaders:\n" +
		"            - name: \"brief\"\n" +
		"              value: \"true\"\n" +
		"          failureThreshold: 3\n" +
		"          initialDelaySeconds: 5\n" +
		"          periodSeconds: 30\n" +
		"          successThreshold: 1\n" +
		"          timeoutSeconds: 5\n" +
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
		"        image: \"{{ .ImageRepository }}/cilium:v1.9.1\"\n" +
		"        imagePullPolicy: IfNotPresent\n" +
		"        lifecycle:\n" +
		"          postStart:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - \"/cni-install.sh\"\n" +
		"              - \"--enable-debug=false\"\n" +
		"          preStop:\n" +
		"            exec:\n" +
		"              command:\n" +
		"              - /cni-uninstall.sh\n" +
		"        name: cilium-agent\n" +
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
		"        image: \"{{ .ImageRepository }}/cilium:v1.9.1\"\n" +
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
		"          mountPropagation: HostToContainer\n" +
		"        - mountPath: /var/run/cilium\n" +
		"          name: cilium-run\n" +
		"        resources:\n" +
		"          requests:\n" +
		"            cpu: 100m\n" +
		"            memory: 100Mi\n" +
		"      restartPolicy: Always\n" +
		"      priorityClassName: system-node-critical\n" +
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
		"---\n" +
		"# Source: cilium/templates/cilium-operator-deployment.yaml\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  labels:\n" +
		"    io.cilium/app: operator\n" +
		"    name: cilium-operator\n" +
		"  name: cilium-operator\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  # We support HA mode only for Kubernetes version > 1.14\n" +
		"  # See docs on ServerCapabilities.LeasesResourceLock in file pkg/k8s/version/version.go\n" +
		"  # for more details.\n" +
		"  {{ if .Single }}\n" +
		"  replicas: 1\n" +
		"  {{else}}\n" +
		"  replicas: 2\n" +
		"  {{end}}\n" +
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
		"      # In HA mode, cilium-operator pods must not be scheduled on the same\n" +
		"      # node as they will clash with each other.\n" +
		"      affinity:\n" +
		"        podAntiAffinity:\n" +
		"          requiredDuringSchedulingIgnoredDuringExecution:\n" +
		"          - labelSelector:\n" +
		"              matchExpressions:\n" +
		"              - key: io.cilium/app\n" +
		"                operator: In\n" +
		"                values:\n" +
		"                - operator\n" +
		"            topologyKey: kubernetes.io/hostname\n" +
		"      containers:\n" +
		"      - args:\n" +
		"        - --config-dir=/tmp/cilium/config-map\n" +
		"        - --debug=$(CILIUM_DEBUG)\n" +
		"        command:\n" +
		"        - cilium-operator-generic\n" +
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
		"        - name: CILIUM_DEBUG\n" +
		"          valueFrom:\n" +
		"            configMapKeyRef:\n" +
		"              key: debug\n" +
		"              name: cilium-config\n" +
		"              optional: true\n" +
		"        image: \"{{ .ImageRepository }}/cilium-operator:v1.9.1\"\n" +
		"        imagePullPolicy: IfNotPresent\n" +
		"        name: cilium-operator\n" +
		"        livenessProbe:\n" +
		"          httpGet:\n" +
		"            host: '127.0.0.1'\n" +
		"            path: /healthz\n" +
		"            port: 9234\n" +
		"            scheme: HTTP\n" +
		"          initialDelaySeconds: 60\n" +
		"          periodSeconds: 10\n" +
		"          timeoutSeconds: 3\n" +
		"        volumeMounts:\n" +
		"        - mountPath: /tmp/cilium/config-map\n" +
		"          name: cilium-config-path\n" +
		"          readOnly: true\n" +
		"      hostNetwork: true\n" +
		"      restartPolicy: Always\n" +
		"      priorityClassName: system-cluster-critical\n" +
		"      serviceAccount: cilium-operator\n" +
		"      serviceAccountName: cilium-operator\n" +
		"      tolerations:\n" +
		"      - operator: Exists\n" +
		"      volumes:\n" +
		"        # To read the configuration from the config map\n" +
		"      - configMap:\n" +
		"          name: cilium-config\n" +
		"        name: cilium-config-path\n" +
		""
	return tmpl
}
