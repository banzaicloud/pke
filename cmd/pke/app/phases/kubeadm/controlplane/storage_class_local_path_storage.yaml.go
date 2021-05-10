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

// storageClassLocalPathStorageTemplate is a generated function returning the template as a string.
func storageClassLocalPathStorageTemplate() string {
	var tmpl = "apiVersion: v1\n" +
		"kind: ServiceAccount\n" +
		"metadata:\n" +
		"  name: local-path-provisioner-service-account\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1beta1\n" +
		"kind: ClusterRole\n" +
		"metadata:\n" +
		"  name: local-path-provisioner-role\n" +
		"  namespace: kube-system\n" +
		"rules:\n" +
		"- apiGroups: [\"\"]\n" +
		"  resources: [\"nodes\", \"persistentvolumeclaims\"]\n" +
		"  verbs: [\"get\", \"list\", \"watch\"]\n" +
		"- apiGroups: [\"\"]\n" +
		"  resources: [\"endpoints\", \"persistentvolumes\", \"pods\"]\n" +
		"  verbs: [\"*\"]\n" +
		"- apiGroups: [\"\"]\n" +
		"  resources: [\"events\"]\n" +
		"  verbs: [\"create\", \"patch\"]\n" +
		"- apiGroups: [\"storage.k8s.io\"]\n" +
		"  resources: [\"storageclasses\"]\n" +
		"  verbs: [\"get\", \"list\", \"watch\"]\n" +
		"---\n" +
		"apiVersion: rbac.authorization.k8s.io/v1\n" +
		"kind: ClusterRoleBinding\n" +
		"metadata:\n" +
		"  name: local-path-provisioner-bind\n" +
		"  namespace: kube-system\n" +
		"roleRef:\n" +
		"  apiGroup: rbac.authorization.k8s.io\n" +
		"  kind: ClusterRole\n" +
		"  name: local-path-provisioner-role\n" +
		"subjects:\n" +
		"- kind: ServiceAccount\n" +
		"  name: local-path-provisioner-service-account\n" +
		"  namespace: kube-system\n" +
		"---\n" +
		"apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  name: local-path-provisioner\n" +
		"  namespace: kube-system\n" +
		"spec:\n" +
		"  replicas: 1\n" +
		"  selector:\n" +
		"    matchLabels:\n" +
		"      app: local-path-provisioner\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      labels:\n" +
		"        app: local-path-provisioner\n" +
		"    spec:\n" +
		"      serviceAccountName: local-path-provisioner-service-account\n" +
		"      containers:\n" +
		"      - name: local-path-provisioner\n" +
		"        image: {{ .ImageRepository }}/local-path-provisioner:v0.0.9\n" +
		"        imagePullPolicy: Always\n" +
		"        command:\n" +
		"        - local-path-provisioner\n" +
		"        - --debug\n" +
		"        - start\n" +
		"        - --config\n" +
		"        - /etc/config/config.json\n" +
		"        - --provisioner-name=banzaicloud.io/local-path\n" +
		"        volumeMounts:\n" +
		"        - name: config-volume\n" +
		"          mountPath: /etc/config/\n" +
		"        env:\n" +
		"        - name: POD_NAMESPACE\n" +
		"          valueFrom:\n" +
		"            fieldRef:\n" +
		"              fieldPath: metadata.namespace\n" +
		"      volumes:\n" +
		"        - name: config-volume\n" +
		"          configMap:\n" +
		"            name: local-path-config\n" +
		"---\n" +
		"apiVersion: storage.k8s.io/v1\n" +
		"kind: StorageClass\n" +
		"metadata:\n" +
		"  name: local-path\n" +
		"  annotations:\n" +
		"    storageclass.kubernetes.io/is-default-class: \"true\"\n" +
		"provisioner: banzaicloud.io/local-path\n" +
		"volumeBindingMode: WaitForFirstConsumer\n" +
		"reclaimPolicy: Delete\n" +
		"---\n" +
		"kind: ConfigMap\n" +
		"apiVersion: v1\n" +
		"metadata:\n" +
		"  name: local-path-config\n" +
		"  namespace: kube-system\n" +
		"data:\n" +
		"  config.json: |-\n" +
		"        {\n" +
		"                \"nodePathMap\":[\n" +
		"                {\n" +
		"                        \"node\":\"DEFAULT_PATH_FOR_NON_LISTED_NODES\",\n" +
		"                        \"paths\":[\"/opt/local-path-provisioner\"]\n" +
		"                }\n" +
		"                ]\n" +
		"        }\n" +
		""
	return tmpl
}
