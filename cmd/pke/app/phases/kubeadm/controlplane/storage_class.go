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

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

func applyDefaultStorageClass(out io.Writer, disableDefaultStorageClass bool, cloudProvider, azureStorageAccountType, azuerStorageKind string) error {
	if disableDefaultStorageClass {
		return nil
	}

	var err error
	switch cloudProvider {
	case constants.CloudProviderAmazon:
		err = writeStorageClassAmazon(out, storageClassConfig)
	case constants.CloudProviderAzure:
		err = writeStorageClassAzure(out, storageClassConfig, azureStorageAccountType, azuerStorageKind)
	default:
		err = writeStorageClassLocalPathStorage(out, storageClassConfig)
	}
	if err != nil {
		return err
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", storageClassConfig)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	if err != nil {
		return err
	}

	return nil
}

func writeStorageClassAmazon(out io.Writer, filename string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating Amazon default storage class\n", use)
	// https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs
	conf := `kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: gp2
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
`

	tmpl, err := template.New("storage-class-amazon").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	type data struct{}

	d := data{}

	return tmpl.Execute(w, d)
}

func writeStorageClassAzure(out io.Writer, filename string, storageAccountType, kind string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating Azure default storage class\n", use)
	// https://kubernetes.io/docs/concepts/storage/storage-classes/#new-azure-disk-storage-class-starting-from-v1-7-2
	conf := `kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: azure-disk
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  labels:
    kubernetes.io/cluster-service: "true"
provisioner: kubernetes.io/azure-disk
volumeBindingMode: WaitForFirstConsumer
parameters:
  storageaccounttype: {{ .StorageAccountType }}
  kind: {{ .Kind }}
`

	tmpl, err := template.New("storage-class-azure").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	type data struct {
		StorageAccountType string
		Kind               string
	}

	d := data{
		StorageAccountType: storageAccountType,
		Kind:               kind,
	}

	return tmpl.Execute(w, d)
}

func writeStorageClassLocalPathStorage(out io.Writer, filename string) error {
	_, _ = fmt.Fprintf(out, "[%s] creating local default storage class\n", use)
	conf := `apiVersion: v1
kind: ServiceAccount
metadata:
  name: local-path-provisioner-service-account
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: local-path-provisioner-role
  namespace: kube-system
rules:
- apiGroups: [""]
  resources: ["nodes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["endpoints", "persistentvolumes", "pods"]
  verbs: ["*"]
- apiGroups: [""]
  resources: ["events"]
  verbs: ["create", "patch"]
- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: local-path-provisioner-bind
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: local-path-provisioner-role
subjects:
- kind: ServiceAccount
  name: local-path-provisioner-service-account
  namespace: kube-system
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: local-path-provisioner
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: local-path-provisioner
  template:
    metadata:
      labels:
        app: local-path-provisioner
    spec:
      serviceAccountName: local-path-provisioner-service-account
      containers:
      - name: local-path-provisioner
        image: banzaicloud/local-path-provisioner:v0.0.5
        imagePullPolicy: Always
        command:
        - local-path-provisioner
        - --debug
        - start
        - --config
        - /etc/config/config.json
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config/
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      volumes:
        - name: config-volume
          configMap:
            name: local-path-config
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-path
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: banzaicloud.io/local-path
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Delete
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: local-path-config
  namespace: kube-system
data:
  config.json: |-
        {
                "nodePathMap":[
                {
                        "node":"DEFAULT_PATH_FOR_NON_LISTED_NODES",
                        "paths":["/opt/local-path-provisioner"]
                }
                ]
        }
`

	tmpl, err := template.New("storage-class-local-path").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	type data struct{}

	d := data{}

	return tmpl.Execute(w, d)
}
