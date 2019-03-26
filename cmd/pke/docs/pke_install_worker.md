## pke install worker

Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Worker node

### Synopsis

Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Worker node

```
pke install worker [flags]
```

### Options

```
  -h, --help                                        help for worker
      --image-repository string                     Prefix for image repository (default "banzaicloud")
      --kubernetes-api-server string                Kubernetes API Server host port
      --kubernetes-api-server-ca-cert-hash string   CA cert hash
      --kubernetes-cloud-provider string            cloud provider. example: aws
      --kubernetes-infrastructure-cidr string       network CIDR for the actual machine (default "192.168.64.0/20")
      --kubernetes-node-token string                PKE join token
      --kubernetes-version string                   Kubernetes version (default "1.14.0")
      --pipeline-cluster-id int32                   Cluster ID to use with Pipeline API
      --pipeline-nodepool string                    name of the nodepool the node belongs to
      --pipeline-org-id int32                       Organization ID to use with Pipeline API
  -t, --pipeline-token string                       Token for accessing Pipeline API
  -u, --pipeline-url string                         Pipeline API server url
```

### SEE ALSO

* [pke install](pke_install.md)	 - Install a single Banzai Cloud Pipeline Kubernetes Engine (PKE) machine
* [pke install worker container-runtime](pke_install_worker_container-runtime.md)	 - Container runtime installation
* [pke install worker kubernetes-node](pke_install_worker_kubernetes-node.md)	 - Kubernetes worker node installation
* [pke install worker kubernetes-runtime](pke_install_worker_kubernetes-runtime.md)	 - Kubernetes runtime installation
* [pke install worker kubernetes-version](pke_install_worker_kubernetes-version.md)	 - Check Kubernetes version is supported or not
* [pke install worker pipeline-ready](pke_install_worker_pipeline-ready.md)	 - Register node as ready at Pipeline

