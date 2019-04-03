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
      --kubelet-certificate-authority string              Path to a cert file for the certificate authority. Used for kubelet server certificate verify. (default "/etc/kubernetes/pki/ca.crt")
      --kubernetes-advertise-address string               Kubernetes API Server advertise address
      --kubernetes-api-server string                      Kubernetes API Server host port
      --kubernetes-api-server-ca-cert-hash string         CA cert hash
      --kubernetes-api-server-cert-sans stringArray       sets extra Subject Alternative Names for the API Server signing cert
      --kubernetes-cloud-provider string                  cloud provider. example: aws
      --kubernetes-cluster-name string                    Kubernetes cluster name (default "pke")
      --kubernetes-controller-manager-signing-ca string   Kubernetes Controller Manager signing cert
      --kubernetes-infrastructure-cidr string             network CIDR for the actual machine (default "192.168.64.0/20")
      --kubernetes-join-control-plane                     Join an another control plane node
      --kubernetes-master-mode string                     Kubernetes cluster mode (default "default")
      --kubernetes-network-provider string                Kubernetes network provider (default "weave")
      --kubernetes-node-token string                      PKE join token
      --kubernetes-oidc-client-id string                  A client ID that all OIDC tokens must be issued for
      --kubernetes-oidc-issuer-url string                 URL of the OIDC provider which allows the API server to discover public signing keys
      --kubernetes-pod-network-cidr string                range of IP addresses for the pod network (default "10.20.0.0/16")
      --kubernetes-service-cidr string                    range of IP address for service VIPs (default "10.10.0.0/16")
      --kubernetes-version string                         Kubernetes version (default "1.14.0")
      --pipeline-cluster-id int32                         Cluster ID to use with Pipeline API
      --pipeline-nodepool string                          name of the nodepool the node belongs to
      --pipeline-org-id int32                             Organization ID to use with Pipeline API
  -t, --pipeline-token string                             Token for accessing Pipeline API
  -u, --pipeline-url string                               Pipeline API server url
      --with-plugin-psp                                   Enable PodSecurityPolicy admission plugin
```

### SEE ALSO

* [pke install master](pke_install_master.md)	 - Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Master node

