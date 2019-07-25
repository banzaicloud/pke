# PKE in Azure (manually)

This tutorial will walk through the steps of manually installing PKE clusters to Azure Cloud Services.

If you would like to supercharge your Kubernetes experience and deploy PKE automatically, check out the free developer beta of Banzai Cloud Pipeline:
<p align="center">
  <a href="https://beta.banzaicloud.io">
  <img src="https://camo.githubusercontent.com/a487fb3128bcd1ef9fc1bf97ead8d6d6a442049a/68747470733a2f2f62616e7a6169636c6f75642e636f6d2f696d672f7472795f706970656c696e655f627574746f6e2e737667">
  </a>
</p>


## Creating the infrastructure

> While you can create all the resources with the web based console of Azure, we will provide command line examples assuming that you have an Azure CLI set up to your Azure account.

We will use some shell variables to simplify following the guide. Set them to your liking.

```bash
export LOCATION=westeurope
export RG=pke-azure
export CLUSTER_NAME=$RG
export VNET=$RG-vnet
export SUBNET=$RG-subnet
export NSG=$RG-nsg
export BE_POOL=$CLUSTER_NAME-be-pool
export FE_POOL=$CLUSTER_NAME-fe-pool
export INFRA_CIDR=10.240.0.0/24
export PUBLIC_IP=$RG-pip
export APISERVER_PROBE=$RG-apiserver-probe
export APISERVER_RULE=$RG-apiserver-rule
export ROUTES=$RG-routes
export IMAGE=OpenLogic:CentOS-CI:7-CI:7.6.20190306
```

First, you will need to create a Resource Group.

```bash
az group create --name $RG --location $LOCATION
```

Create a Virtual Network with a subnet large enough to assign a private IP address to each node in the PKE Kubernetes cluster.

```bash
az network vnet create -g $RG \
  -n $VNET \
  --address-prefix $INFRA_CIDR \
  --subnet-name $SUBNET \
  --location $LOCATION
```

Create a Firewall (Network Security Group) and assign it to the previously created subnet.

```bash
az network nsg create -g $RG -n $NSG \
  --location $LOCATION

az network vnet subnet update -g $RG \
  -n $SUBNET \
  --vnet-name $VNET \
  --network-security-group $NSG
```

Add firewall rules to allow SSH and HTTPS accees to Kubernetes API Server.

```bash
az network nsg rule create -g $RG \
  -n kubernetes-allow-ssh \
  --access allow \
  --destination-address-prefix '*' \
  --destination-port-range 22 \
  --direction inbound \
  --nsg-name $NSG \
  --protocol tcp \
  --source-address-prefix '*' \
  --source-port-range '*' \
  --priority 1000

az network nsg rule create -g $RG \
  -n kubernetes-allow-api-server \
  --access allow \
  --destination-address-prefix '*' \
  --destination-port-range 6443 \
  --direction inbound \
  --nsg-name $NSG \
  --protocol tcp \
  --source-address-prefix '*' \
  --source-port-range '*' \
  --priority 1001
```

You can verify the created rules.

```bash
# List the firewall rules
az network nsg rule list -g $RG --nsg-name $NSG --query "[].{Name:name, \
  Direction:direction, Priority:priority, Port:destinationPortRange}" -o table
```

Allocate a static IP address that will be attached to the external load balancer fronting the Kubernetes API Server.

```bash
az network public-ip create -n $PUBLIC_IP -g $RG \
    --sku Standard \
    --location $LOCATION
az network lb create -g $RG \
  -n $CLUSTER_NAME \
  --sku Standard \
  --location $LOCATION \
  --backend-pool-name $BE_POOL \
  --frontend-ip-name $FE_POOL \
  --public-ip-address $PUBLIC_IP
```

Check IP address was created correctly in the Resource Group and chosen region.

```bash
az network public-ip list --query="[?name=='$PUBLIC_IP'].{ResourceGroup:resourceGroup, \
  Region:location,Allocation:publicIpAllocationMethod,IP:ipAddress}" -o table
```

Create the load balancer health probe as a pre-requesite for the lb rule that follows.

```bash
az network lb probe create -g $RG \
  --lb-name $CLUSTER_NAME \
  --name $APISERVER_PROBE \
  --port 6443 \
  --protocol tcp
```

Create the external load balancer network resource.

```bash
az network lb rule create -g $RG \
  -n $APISERVER_RULE \
  --protocol tcp \
  --lb-name $CLUSTER_NAME \
  --frontend-ip-name $FE_POOL \
  --frontend-port 6443 \
  --backend-pool-name $BE_POOL \
  --backend-port 6443 \
  --probe-name $APISERVER_PROBE
```

Kubernetes internal network cummication needs to manage a route table, let us create one.

```bash
az network route-table create -g $RG -n $ROUTES
az network vnet subnet update -g $RG \
  -n $SUBNET \
  --vnet-name $VNET \
  --route-table $ROUTES
```



### Master node

Create a Virtual Machine for the Kubernetes Controller.

> From now on we will use a private static IP address for every machine for the ease of this tutorial.

```bash
az vm availability-set create -g $RG \
        -n $CLUSTER_NAME \
        --location $LOCATION \
        --validate

az network nic create -g $RG \
        -n controller-0-nic \
        --private-ip-address 10.240.0.10 \
        --vnet $VNET \
        --subnet $SUBNET \
        --ip-forwarding \
        --lb-name $CLUSTER_NAME \
        --lb-address-pools $BE_POOL

az vm create -g $RG \
        -n controller-0 \
        --location $LOCATION \
        --availability-set $CLUSTER_NAME \
        --image $IMAGE \
        --size Standard_B2s \
        --storage-sku Standard_LRS \
        --assign-identity \
        --generate-ssh-keys \
        --admin-username azureuser \
        --nics controller-0-nic
```

PKE will use the integrated Kubernetes Cloud Controller. Controller's Identity rights needs to be elevated for automated resource management to work.

```bash
export CONTROLLER_IDENTITY=$(az vm list -d -g $RG --query '[?contains(name, `controller-0`)].identity.principalId' -o tsv)
az role assignment create --role Owner --assignee-object-id $CONTROLLER_IDENTITY --resource-group $RG
```

Since this installation is manual, we need to access the machines through SSH. No public IP address is associated with the controller, so we create a rule for the load-balancer to enable SSH.

```bash
az network lb inbound-nat-rule create -g $RG \
    -n ssh \
    --lb-name $CLUSTER_NAME \
    --protocol Tcp \
    --frontend-port 50000 \
    --backend-port 22 \
    --frontend-ip-name $FE_POOL \
    --enable-tcp-reset

export CONTROLLER_IP_CONFIG=$(az network nic ip-config list \
    -g $RG \
    --nic-name controller-0-nic \
    --query '[0].name' \
    -o json | tr -d '"')

az network nic ip-config inbound-nat-rule add \
    --inbound-nat-rule ssh \
    -g $RG \
    --lb-name $CLUSTER_NAME \
    --nic-name controller-0-nic \
    --ip-config-name $CONTROLLER_IP_CONFIG
```



### Worker node (optional, needed by the multi-node installation)

This exmaple creates two worker node vms, but you can change it to your liking.

> Edit the third number after seq at the line starting with for. The maximal vaule is zero indexed, use desired amount of workers minus one.
> For a single worker node use: `seq 0 1 0`
> For three worker nodes use: `seq 0 1 2`

```bash
for i in $(seq 0 1 1); do
az network nic create -g $RG \
        -n worker-${i}-nic \
        --private-ip-address 10.240.0.2${i} \
        --vnet $VNET \
        --subnet $SUBNET \
        --ip-forwarding

az vm create -g $RG \
        -n worker-${i} \
        --location $LOCATION \
        --availability-set $CLUSTER_NAME \
        --image $IMAGE \
        --size Standard_B2s \
        --storage-sku Standard_LRS \
        --assign-identity \
        --generate-ssh-keys \
        --admin-username azureuser \
        --nics worker-${i}-nic
done
```



## Setting up PKE

### Single node

Connect using SSH.

```bash
export CONTROLLER_IP=$(az network public-ip show -g $RG -n $PUBLIC_IP --query ipAddress -otsv)
echo $CONTROLLER_IP
ssh azureuser@$CONTROLLER_IP -p 50000
[azureuser@controller-0 ~]$ sudo -s
```

We are going to run every command from now on as root.

> Hint: tenant id is printed during az login.

```bash
export TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 
export RG=pke-azure
export VNET=$RG-vnet
export SUBNET=$RG-subnet
export NSG=$RG-nsg
export LB_IP=$RG-lb
export CONTROLLER_IP=<copy the printed value from above>
export INFRA_CIDR=10.240.0.0/24
export PRIVATEIP=10.240.0.10
export ROUTES=$RG-routes
```

Install PKE.

> `--kubernetes-cluster-name` is used for load balancer naming.

```bash
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install master --kubernetes-cloud-provider=azure \
--azure-tenant-id=$TENANT_ID \
--azure-subnet-name=$SUBNET \
--azure-security-group-name=$NSG \
--azure-vnet-name=$VNET \
--azure-vnet-resource-group=$RG \
--azure-vm-type=standard \
--azure-loadbalancer-sku=standard \
--azure-route-table-name=$ROUTES \
--kubernetes-master-mode=single \
--kubernetes-cluster-name=$CLUSTER_NAME \
--kubernetes-advertise-address=$PRIVATEIP:6443 \
--kubernetes-api-server=$PRIVATEIP:6443 \
--kubernetes-infrastructure-cidr=$INFRA_CIDR \
--kubernetes-api-server-cert-sans="$CONTROLLER_IP"
```



### Multi node

#### Master node

Meanwhile the controller-0 machine is up and running, we can connect to it using SSH.

```bash
export CONTROLLER_IP=$(az network public-ip show -g $RG -n $PUBLIC_IP --query ipAddress -otsv)
echo $CONTROLLER_IP
ssh azureuser@$CONTROLLER_IP -p 50000
[azureuser@controller-0 ~]$ sudo -s
```

We are going to run every command from now on as root.

> Hint: tenant id is printed during az login.

```bash
export TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 
export RG=pke-azure
export VNET=$RG-vnet
export SUBNET=$RG-subnet
export NSG=$RG-nsg
export LB_IP=$RG-lb
export CONTROLLER_IP=<copy the printed value from above>
export INFRA_CIDR=10.240.0.0/24
export PRIVATEIP=10.240.0.10
export ROUTES=$RG-routes
```

Install PKE.

> `--kubernetes-cluster-name` is used for load balancer naming.

```bash
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install master --kubernetes-cloud-provider=azure \
--azure-tenant-id=$TENANT_ID \
--azure-subnet-name=$SUBNET \
--azure-security-group-name=$NSG \
--azure-vnet-name=$VNET \
--azure-vnet-resource-group=$RG \
--azure-vm-type=standard \
--azure-loadbalancer-sku=standard \
--azure-route-table-name=$ROUTES \
--kubernetes-cluster-name=$CLUSTER_NAME \
--kubernetes-advertise-address=$PRIVATEIP:6443 \
--kubernetes-api-server=$PRIVATEIP:6443 \
--kubernetes-infrastructure-cidr=$INFRA_CIDR \
--kubernetes-api-server-cert-sans="$CONTROLLER_IP"

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

Getting credentials for additional nodes

```bash
pke token list
```

Above command will print Token and Cert Hash, please remember it, it will be required in the following step.

#### Worker node(s)

From a new terminal window SSH to worker-0 using controller-0 as jump host (remember no public IP is associated with the machines).

> These commands must be repeated on every worker node by incrementing the last value of the 10.240.0.20 IP address (10.240.0.21, 10.240.0.22...) based on the number of worker nodes.

```bash
export CONTROLLER_IP=$(az network public-ip show -g $RG -n $PUBLIC_IP --query ipAddress -otsv)
ssh -J azureuser@$CONTROLLER_IP:50000 azureuser@10.240.0.20
[azureuser@worker-0 ~]$ sudo -s
```

We are going to run every command from now on as root.

> Hint: tenant id is printed during az login.

```bash
export TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 
export RG=pke-azure
export VNET=$RG-vnet
export SUBNET=$RG-subnet
export NSG=$RG-nsg
export LB_IP=$RG-lb
export CONTROLLER_IP=<copy the printed value from above>
export INFRA_CIDR=10.240.0.0/24
export PRIVATEIP=10.240.0.10
export ROUTES=$RG-routes
export TOKEN=<copy here value from previous step>
export CERTHASH=<copy here value from previous step>
```

Install PKE.

```bash
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install worker --kubernetes-cloud-provider=azure \
--azure-tenant-id=$TENANT_ID \
--azure-subnet-name=$SUBNET \
--azure-security-group-name=$NSG \
--azure-vnet-name=$VNET \
--azure-vnet-resource-group=$RG \
--azure-vm-type=standard \
--azure-loadbalancer-sku=standard \
--azure-route-table-name=$ROUTES \
--kubernetes-api-server=$PRIVATEIP:6443 \
--kubernetes-infrastructure-cidr=$INFRA_CIDR \
--kubernetes-node-token=$TOKEN \
--kubernetes-api-server-ca-cert-hash=$CERTHASH \
--kubernetes-pod-network-cidr=""
```
