## pke install worker kubernetes-node

Kubernetes worker node installation

### Synopsis

Kubernetes worker node installation

```
pke install worker kubernetes-node [flags]
```

### Options

```
  -h, --help                                        help for kubernetes-node
      --kubernetes-api-server string                Kubernetes API Server host port
      --kubernetes-api-server-ca-cert-hash string   CA cert hash
      --kubernetes-cloud-provider string            cloud provider. example: aws
      --kubernetes-node-token string                PKE join token
      --kubernetes-version string                   Kubernetes version (default "1.14.0")
      --pipeline-cluster-id int32                   Cluster ID to use with Pipeline API
      --pipeline-nodepool string                    name of the nodepool the node belongs to
      --pipeline-org-id int32                       Organization ID to use with Pipeline API
  -t, --pipeline-token string                       Token for accessing Pipeline API
  -u, --pipeline-url string                         Pipeline API server url
```

### SEE ALSO

* [pke install worker](pke_install_worker.md)	 - Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Worker node

