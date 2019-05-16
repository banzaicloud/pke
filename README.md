## Pipeline Kubernetes Engine - PKE

PKE is an extremely simple [CNCF certifed](https://www.cncf.io/certification/software-conformance/#logos) Kubernetes installer and distribution, designed to work on any cloud, VM or bare metal. The `pke` tool supports cloud provider integrations, multi-phased installs (requires only an OS), pre-warmed machine image builds, and more. PKE is the preferred Kubernetes run-time of the [Banzai Cloud Pipeline platform](https://beta.banzaicloud.io/), which supercharges the development, deployment and scaling of container-based applications with native support for multi-, hybrid-, and edge-cloud environments.

If you would like to supercharge your Kubernetes experience using Banzai Cloud Pipeline, check out the free developer beta:
<p align="center">
  <a href="https://beta.banzaicloud.io">
  <img src="https://camo.githubusercontent.com/a487fb3128bcd1ef9fc1bf97ead8d6d6a442049a/68747470733a2f2f62616e7a6169636c6f75642e636f6d2f696d672f7472795f706970656c696e655f627574746f6e2e737667">
  </a>
</p>

## Create clusters

Please review the [requirements](/docs/requirements.md) before creating Kubernetes clusters. Note that the `pke` tool will install all required dependencies (like CRI, CNI, etc).

### Single-node PKE

Creating a single-node Kubernetes cluster is as simple as running the following command as root:

`pke install single`

### Multi-node PKE

To create the Kubernetes API server:

```
export MASTER_IP_ADDRESS=""
pke install master --kubernetes-api-server=$MASTER_IP_ADDRESS:6443
```

>Please get the token and certhash from the logs or issue the `pke token list` command to print the token and cert hash needed by workers to join the cluster.
>

Once the API server is up and running you can add as many nodes as needed:

```
export TOKEN=""
export CERTHASH=""
export MASTER_IP_ADDRESS=""
pke install worker --kubernetes-node-token $TOKEN --kubernetes-api-server-ca-cert-hash $CERTHASH --kubernetes-api-server $MASTER_IP_ADDRESS:6443
```

### Using `kubectl`

To use `kubectl` and other command line tools on the master node, set up its config:

```
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes
```

### Other options

You can create PKE clusters on any of the cloud providers, Vagrant, virtual machines, etc using the `pke` tool or let the [Pipeline platform](https://beta.banzaicloud.io/) do it for you using all the complementary features, such as: centralized log collection, federated monitoring, autoscaling, Vault based secret management, disaster recovery, security scans and a lot more.

- [PKE in Vagrant](/docs/vagrant.md)
- [PKE on AWS](/docs/aws.md)
- [PKE on Azure](/docs/azure.md)
- [PKE with Pipeline](https://beta.banzaicloud.io/)

### Contributing

Thank you for your contribution and being part of our community. Please read [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details on the code of conduct, and the process for submitting pull requests. When you are opening a PR to PKE for the first time we will require you to sign a standard CLA.

## Community

If you have any questions about PKE, and would like to talk to us and the other members of the Banzai Cloud community, please join our [Slack](https://slack.banzaicloud.io/).

### License

Copyright (c) 2017-2019 [Banzai Cloud, Inc.](https://banzaicloud.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
