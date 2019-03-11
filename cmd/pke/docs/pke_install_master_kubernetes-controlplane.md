## pke install master kubernetes-controlplane

Kubernetes Control Plane installation

### Synopsis

Kubernetes Control Plane installation

```
pke install master kubernetes-controlplane [flags]
```

### Options

```
  -h, --help                                              help for kubernetes-controlplane
      --image-repository string                           Prefix for image repository (default "banzaicloud")
      --kubernetes-advertise-address string               Kubernetes advertise address
      --kubernetes-api-server string                      Kubernetes API Server host port
      --kubernetes-api-server-cert-sans stringArray       sets extra Subject Alternative Names for the API Server signing cert
      --kubernetes-cloud-provider string                  cloud provider. example: aws
      --kubernetes-cluster-name string                    Kubernetes cluster name (default "pke")
      --kubernetes-controller-manager-signing-ca string   Kubernetes Controller Manager signing cert
      --kubernetes-master-mode string                     Kubernetes cluster mode (default "default")
      --kubernetes-network-provider string                Kubernetes network provider (default "weave")
      --kubernetes-oidc-client-id string                  A client ID that all OIDC tokens must be issued for
      --kubernetes-oidc-issuer-url string                 URL of the OIDC provider which allows the API server to discover public signing keys
      --kubernetes-pod-network-cidr string                range of IP addresses for the pod network (default "10.20.0.0/16")
      --kubernetes-service-cidr string                    range of IP address for service VIPs (default "10.10.0.0/16")
      --kubernetes-version string                         Kubernetes version (default "1.13.3")
      --pipeline-nodepool string                          name of the nodepool the node belongs to
```

### SEE ALSO

* [pke install master](pke_install_master.md)	 - Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Master node

