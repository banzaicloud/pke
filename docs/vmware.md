## PKE in VMware

This tutorial will walk through the steps of manually installing PKE clusters on a VMware infrastructure.

## Creating the infrastructure

To install PKE, you will need a virtual machine for each of the nodes you plan to install.
At the time, RHEL 7 and CentOS 7 are supported.
Each node needs at least 2 vCPUs, 2GiB RAM and an 8 GB root volume.

The nodes should have unrestricted IP network access to each other, direct Internet access, and allow incoming access on TCP 6443 from the hosts managing it.

Default CentOS and RHEL installations configure firewalld to restrict access to services running on the host.
If you don't want to configure it to meet the above mentioned criteria, turn it off by `systemctl disable firewalld; systemctl stop firewalld` as the root user.

If you want to deploy multi-master clusters, you will also need a round-robin DNS record pointing to the masters.

If your environment does not allow these conditions, please feel free to [contact us](https://banzaicloud.com/contact/) for a solution.

### Requirements for persistent volume integration

If you want to use vSphere datastores for storing Kubernetes persistent volumes, you will also need to supply the username/password pair of a vCenter service user which can list compute and storage resources, deploy (temporary) virtual machines with new disks, and attach volumes to your nodes.

Please ensure that VMware Tools are installed (`yum install open-vm-tools`).

You will need to set `disk.EnableUUID = "TRUE"` in the VMX file describing your VMware virtual machines on each node.
On previous vCenter versions, this could be set on powered off machines' *edit settings* dialog, on the *VM options* tab behind *Edit configuration*.
On recent versions you will have to shut down the machine, browse the VMX file on the data store, download it, modify with a text editor, and replace by an upload, then power on the machine.
If you use VM templates, you can do this once on the template.

For the next steps you will need to define some shell variables for the vSphere persistent volume integration in each shell session where you will install PKE:

```
server=my.vcenter.local          # the hostname or IP of the vCenter instance that hosts the nodes
port=443                         # the TCP port where vCenter listens
fingerprint=AA:...:CC:DD         # the fingerprint of the server certificate of vCenter to use; you can use `openssl s_client -connect $server:$port </dev/null | openssl x509 -fingerprint -noout` to determine it
datacenter=Datacenter            # the name of the vSphere datacenter to use to store persistent volumes (and deploy temporary VMs to create them)
datastore=Datastore              # the name of the vSphere datastore that is in the given datacenter, and is available on all nodes (datastore clusters are not supported)
resourcepool=Cluster/Pool        # the path of the vSphere resource pool to create temporary VMs in during volume creation
folder=folder                    # the name of the vSphere folder (aka blue folder) to create temporary VMs in during volume creation as well as all Kubernetes nodes are there
username=username@vcenter.local  # the name of vCenter SSO user to use for deploying persistent volumes (should be avoided in favor of a K8S secret)
password=Password12              # the password of vCenter SSO user to use for deploying persistent volumes (should be avoided in favor of a K8S secret)
lbrange=192.168.22.200-210       # IPv4 range that will be advertised via ARP and where LoadBalancer Services will be served (optional, example: 192.168.0.100-192.168.0.110)
```

### Single node PKE

Launch a single virtual machine that meets the above mentioned requirements, and log in to it to execute the following commands.

> If you don't need the vSphere persistent volume integration, leave out the `pke install` arguments starting from `--kubernetes-cloud-provider`.

```
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install single \
        --kubernetes-cloud-provider vsphere    \
        --vsphere-server=$server               \
	--vsphere-port=$port                   \
	--vsphere-fingerprint=$fingerprint     \
	--vsphere-datacenter=$datacenter       \
	--vsphere-datastore=$datastore         \
	--vsphere-resourcepool=$resourcepool   \
	--vsphere-folder=$folder               \
	--vsphere-username=$username           \
	--vsphere-password=$password           \
        --lb-range=$lbrange

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

You can use `kubectl` from now on. Test it by executing `kubectl get nodes`

### Multi node PKE

#### Start the master node

> Please take care of setting different hostnames for all the nodes in the cluster. You can do it with `hostnamectl set-hostname my-first-pke-node`.
> You may encounter issues if the hostname is corrected later.

In *addition* to the variables described earlier, you will need the IP address or hostname where the master will be available to the other nodes and the clients (who will access the Kubernetes API server):

```
master=192.0.2.1
```

```
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install master \
        --kubernetes-advertise-address=$master \
        --kubernetes-api-server=$master:6443   \
        --kubernetes-cloud-provider vsphere    \
        --vsphere-server=$server               \
	--vsphere-port=$port                   \
	--vsphere-fingerprint=$fingerprint     \
	--vsphere-datacenter=$datacenter       \
	--vsphere-datastore=$datastore         \
	--vsphere-resourcepool=$resourcepool   \
	--vsphere-folder=$folder               \
	--vsphere-username=$username           \
	--vsphere-password=$password           \
        --lb-range=$lbrange

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```


Please note the token and certhash from the logs printed or issue the following PKE command to print the token and cert hash needed by workers to join the cluster.

```
pke token list # Print token and cert hash needed by workers to join the cluster
```

If the tokens have expired or you'd like to create a new one, issue:

```
pke token create
```

#### Start a worker node

Log in to the worker nodes, and set the TOKEN and CERTHASH variables with the values from the previous step, as well as `$master`.

```
TOKEN=""
CERTHASH=""
```

Install the node:
```
pke install worker \
        --kubernetes-node-token $TOKEN         \
        --kubernetes-api-server-ca-cert-hash $CERTHASH  \
        --kubernetes-api-server $master:6443   \
        --kubernetes-cloud-provider vsphere
```

Note that you can add as many worker nodes as you wish repeating the commands above. You can check the status of the containers by issuing `crictl ps`

### Creating templates

We recommend to create a template in vSphere for nodes. Simply install an operating system to meet your requirements and the ones at the beginning of this guide.
After that, download PKE, and install the components it would normally download and install at for each node separately:

```
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke machine-image
```

### Further options

#### Network providers

Traffic between Pods is by default encapsulated by Weave, which will work in most cases without further configuration, by might not be needed --- and worth the overhead --- in the most popular network scenarios.

If you want to prevent encapsulating packets between your pods, you can use Calico with native IP packets (the default setting, `ipipMode: Never`) when all your nodes are in a single Layer 2 network.
You can achieve this by supplying `--kubernetes-network-provider=calico` to the install commands on each node.

#### Load balancer integration

If you want to use `LoadBalancer` type Kubernetes Service resources to expose your services, you can easily do this by specifying the `--lb-range` option to master installation command.
This will install and configure MetalLB to allocate IP addresses in the specified range, and advertise them via gratuitous ARP on the host's network.

The service will take care of advertising the addresses on exactly one Ready node, and forwarding incoming connections from any node it may hit to one of the specified Pods.
