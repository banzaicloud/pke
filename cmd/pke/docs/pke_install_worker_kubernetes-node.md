---
title: pke install worker kubernetes-node
generated_file: true
---
## pke install worker kubernetes-node

Kubernetes worker node installation

### Synopsis

Kubernetes worker node installation

```
pke install worker kubernetes-node [flags]
```

### Options

```
      --azure-loadbalancer-sku string               Sku of Load Balancer and Public IP. Candidate values are: basic and standard (default "basic")
      --azure-route-table-name string               The name of the route table attached to the subnet that the cluster is deployed in (default "kubernetes-routes")
      --azure-security-group-name string            The name of the security group attached to the cluster's subnet
      --azure-subnet-name string                    The name of the subnet that the cluster is deployed in
      --azure-tenant-id string                      The AAD Tenant ID for the Subscription that the cluster is deployed in
      --azure-vm-type string                        The type of azure nodes. Candidate values are: vmss and standard (default "standard")
      --azure-vnet-name string                      The name of the VNet that the cluster is deployed in
      --azure-vnet-resource-group string            The name of the resource group that the Vnet is deployed in
  -h, --help                                        help for kubernetes-node
      --kubernetes-api-server string                Kubernetes API Server host port
      --kubernetes-api-server-ca-cert-hash string   CA cert hash
      --kubernetes-cloud-provider string            cloud provider. example: aws
      --kubernetes-container-runtime string         Kubernetes container runtime (default "containerd")
      --kubernetes-node-labels strings              Specifies the labels the Node should be registered with
      --kubernetes-node-token string                PKE join token
      --kubernetes-pod-network-cidr string          range of IP addresses for the pod network on the current node
      --kubernetes-version string                   Kubernetes version (default "1.19.10")
      --pipeline-cluster-id int32                   Cluster ID to use with Pipeline API
      --pipeline-insecure                           If the Pipeline API should not verify the API's certificate
      --pipeline-nodepool string                    name of the nodepool the node belongs to
      --pipeline-org-id int32                       Organization ID to use with Pipeline API
  -t, --pipeline-token string                       Token for accessing Pipeline API
  -u, --pipeline-url string                         Pipeline API server url
      --reset-on-failure                            Roll back changes after failures
      --taints strings                              Specifies the taints the Node should be registered with
```

### SEE ALSO

* [pke install worker](/docs/pke/cli/reference/pke_install_worker/)	 - Installs Banzai Cloud Pipeline Kubernetes Engine (PKE) Worker node

