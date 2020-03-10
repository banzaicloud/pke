## pke install single kubernetes-controlplane

Kubernetes Control Plane installation

### Synopsis

Kubernetes Control Plane installation

```
pke install single kubernetes-controlplane [flags]
```

### Options

```
      --azure-loadbalancer-sku string                     Sku of Load Balancer and Public IP. Candidate values are: basic and standard (default "basic")
      --azure-route-table-name string                     The name of the route table attached to the subnet that the cluster is deployed in (default "kubernetes-routes")
      --azure-security-group-name string                  The name of the security group attached to the cluster's subnet
      --azure-storage-account-type string                 Azure storage account Sku tier (default "Standard_LRS")
      --azure-storage-kind string                         Possible values are shared, dedicated, and managed (default "dedicated")
      --azure-subnet-name string                          The name of the subnet that the cluster is deployed in
      --azure-tenant-id string                            The AAD Tenant ID for the Subscription that the cluster is deployed in
      --azure-vm-type string                              The type of azure nodes. Candidate values are: vmss and standard (default "standard")
      --azure-vnet-name string                            The name of the VNet that the cluster is deployed in
      --azure-vnet-resource-group string                  The name of the resource group that the Vnet is deployed in
      --disable-default-storage-class                     Do not deploy a default storage class
      --encryption-secret string                          Use this key to encrypt secrets (32 byte base64 encoded)
      --etcd-ca-file string                               An SSL Certificate Authority file used to secure etcd communication
      --etcd-cert-file string                             An SSL certification file used to secure etcd communication
      --etcd-endpoints strings                            Endpoints of etcd members
      --etcd-key-file string                              An SSL key file used to secure etcd communication
      --etcd-prefix string                                The prefix to prepend to all resource paths in etcd
  -h, --help                                              help for kubernetes-controlplane
      --image-repository string                           Prefix for image repository (default "banzaicloud")
      --kubelet-certificate-authority string              Path to a cert file for the certificate authority. Used for kubelet server certificate verify. (default "/etc/kubernetes/pki/ca.crt")
      --kubernetes-advertise-address string               Kubernetes API Server advertise address
      --kubernetes-api-server string                      Kubernetes API Server host port
      --kubernetes-api-server-ca-cert-hash string         CA cert hash
      --kubernetes-api-server-cert-sans strings           sets extra Subject Alternative Names for the API Server signing cert
      --kubernetes-cloud-provider string                  cloud provider. example: aws
      --kubernetes-cluster-name string                    Kubernetes cluster name (default "pke")
      --kubernetes-controller-manager-signing-ca string   Kubernetes Controller Manager signing cert
      --kubernetes-infrastructure-cidr string             network CIDR for the actual machine (default "192.168.64.0/20")
      --kubernetes-join-control-plane                     Join an another control plane node
      --kubernetes-master-mode string                     Kubernetes cluster mode (default "default")
      --kubernetes-mtu uint                               maximum transmission unit. 0 means default value of the Kubernetes network provider is used
      --kubernetes-network-provider string                Kubernetes network provider (default "calico")
      --kubernetes-node-labels strings                    Specifies the labels the Node should be registered with
      --kubernetes-node-token string                      PKE join token
      --kubernetes-oidc-client-id string                  A client ID that all OIDC tokens must be issued for
      --kubernetes-oidc-issuer-url string                 URL of the OIDC provider which allows the API server to discover public signing keys
      --kubernetes-pod-network-cidr string                range of IP addresses for the pod network (default "10.20.0.0/16")
      --kubernetes-service-cidr string                    range of IP address for service VIPs (default "10.10.0.0/16")
      --kubernetes-version string                         Kubernetes version (default "1.16.0")
      --lb-range string                                   Advertise the specified IPv4 range via ARP and allocate addresses for LoadBalancer Services (non-cloud only, example: 192.168.0.100-192.168.0.110)
      --pipeline-cluster-id int32                         Cluster ID to use with Pipeline API
      --pipeline-insecure                                 If the Pipeline API should not verify the API's certificate
      --pipeline-nodepool string                          name of the nodepool the node belongs to
      --pipeline-org-id int32                             Organization ID to use with Pipeline API
  -t, --pipeline-token string                             Token for accessing Pipeline API
  -u, --pipeline-url string                               Pipeline API server url
      --taints strings                                    Specifies the taints the Node should be registered with (default [node-role.kubernetes.io/master:NoSchedule])
      --vsphere-datacenter string                         The name of the datacenter to use to store persistent volumes (and deploy temporary VMs to create them)
      --vsphere-datastore string                          The name of the datastore that is in the given datacenter, and is available on all nodes
      --vsphere-fingerprint string                        The fingerprint of the server certificate of vCenter to use
      --vsphere-folder string                             The name of the folder (aka blue folder) to create temporary VMs in during volume creation, as well as all Kubernetes nodes are in
      --vsphere-password string                           The password of vCenter SSO user to use for deploying persistent volumes (should be avoided in favor of a K8S secret)
      --vsphere-port int                                  The TCP port where vCenter listens (default 443)
      --vsphere-resourcepool string                       The path of the resource pool to create temporary VMs in during volume creation (for example "Cluster/Pool")
      --vsphere-server string                             The hostname or IP of vCenter to use
      --vsphere-username string                           The name of vCenter SSO user to use for deploying persistent volumes (Should be avoided in favor of a K8S secret)
      --with-plugin-psp                                   Enable PodSecurityPolicy admission plugin
      --without-plugin-deny-escalating-exec               Disable DenyEscalatingExec admission plugin
      --without-audit-log                                 Disable apiserver audit log
```

### SEE ALSO

* [pke install single](pke_install_single.md)	 - Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) on a single machine

