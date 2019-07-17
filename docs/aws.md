## PKE in AWS (manually)

This tutorial will walk through the steps of manually installing PKE clusters to Amazon's cloud services.

If you would like to supercharge your Kubernetes experience and deploy PKE to AWS clusters automatically, check out the free developer beta of Banzai Cloud Pipeline:
<p align="center">
  <a href="https://beta.banzaicloud.io">
  <img src="https://camo.githubusercontent.com/a487fb3128bcd1ef9fc1bf97ead8d6d6a442049a/68747470733a2f2f62616e7a6169636c6f75642e636f6d2f696d672f7472795f706970656c696e655f627574746f6e2e737667">
  </a>
</p>



## Creating the infrastructure

>While you can create all the resources with the web based console of AWS, we will provide command line examples assuming that you have an AWS CLI set up to your Amazon account.

In most production setups you will need (or already have) a different network layout on your AWS environment, but we will stick to the default networks created by Amazon for simplicity.
Feel free to use any VPC and subnet that TCP/IP connections between the nodes, and allows them to download the needed OS packages and Docker images.

We will use some shell variables to simplify following the guide. Set them to your liking.

```
export CLUSTER_NAME=testcluster
export AWS_SSH_KEY_PAIR_NAME=my-ssh-key
export AWS_DEFAULT_REGION=eu-west-1
```

First, you will need to create IAM Roles and Instance Profiles for the EC2 instances serving your Kubernetes nodes to allow the Amazon integrations work: to create load balancers and persitent volumes, or to retrieve information about their environment.

The easiest way is to use our CloudFormation [template](https://raw.githubusercontent.com/banzaicloud/pipeline/0.14.3/templates/pke/global.cf.yaml) from Pipeline.

>Please note that these resources are region-independent. You can create them in your preferred region, but will see them from all the others as well. They will be suitable for all of your PKE clusters under that specific AWS account.

Submit the CF template with the following command:

```
aws cloudformation create-stack \
--stack-name pke-global \
--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
--template-body https://raw.githubusercontent.com/banzaicloud/pipeline/0.14.0/templates/pke/global.cf.yaml
```

Wait until the template is executed successfully. You can check the status with:

```
aws cloudformation describe-stacks --stack-name pke-global
```

The easiest way to configure the machines is to allow SSH traffic to them. To manage the new Kubernetes cluster remotely, you will also need to access the Kubernetes API server on the master node. In case of multi-node clusters, the nodes of the cluster should access each other as well.

Create a security group for the cluster nodes with SSH access:
```
aws ec2 create-security-group --group-name pke-cluster --description "PKE security group"
aws ec2 authorize-security-group-ingress --group-name pke-cluster --protocol tcp --port 22 --cidr 0.0.0.0/0
PKE_CLUSTER_SEC_GROUP_ID=$(aws ec2 describe-security-groups --group-name pke-cluster --region eu-west-1 --query "SecurityGroups[*].GroupId" --output=text)
aws ec2 authorize-security-group-ingress --group-name pke-cluster --source-group $PKE_CLUSTER_SEC_GROUP_ID --protocol -1
```

Create an additional security group for the master node which allows remote HTTPS access to the API server:
```
aws ec2 create-security-group --group-name pke-master --description "PKE master security group"
aws ec2 authorize-security-group-ingress --group-name pke-master --protocol tcp --port 6443 --cidr 0.0.0.0/0
```

After that, we can create the EC2 instances that will host the nodes of the cluster. You should create one master instance, and any number of worker nodes.
You can use any OS AMI image that meets our requirements. You can check the AMI numbers we use in Pipeline [here](https://github.com/banzaicloud/pipeline/blob/0.14.3/internal/providers/pke/pkeworkflow/create_cluster.go#L29).

To create a master instance, run:
```
aws ec2 run-instances --image-id ami-3548444c \
--count 1 \
--instance-type c3.xlarge \
--key-name $AWS_SSH_KEY_PAIR_NAME \
--tag-specifications "ResourceType=instance,Tags=[{Key=kubernetes.io/cluster/$CLUSTER_NAME,Value=owned},{Key=Name,Value=pke-master}]" \
--security-groups pke-cluster pke-master \
--iam-instance-profile Name=pke-global-master-profile
```

To create a worker node:
```
aws ec2 run-instances --image-id ami-3548444c \
--count 1 \
--instance-type c3.xlarge \
--key-name $AWS_SSH_KEY_PAIR_NAME \
--tag-specifications "ResourceType=instance,Tags=[{Key=kubernetes.io/cluster/$CLUSTER_NAME,Value=owned},{Key=Name,Value=pke-worker}]" \
--security-groups pke-cluster \
--iam-instance-profile Name=pke-global-worker-profile
```

## Setting up PKE
### Single node

Once you single master instance booted up, SSH into it with the key file configured. Run the following commands as root:

```
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install single --kubernetes-cloud-provider=aws
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

You can use `kubectl` from now on. Test it by executing `kubectl get nodes`

### Multi node

In case of multi node clusters, you will have to provide something more information to the PKE tool.

#### Install the master node

Find out the address of the master that is accessible to the other nodes, and the clients you want to use the API server. This can be retrieved with a command like:

```
aws ec2 describe-instances --filters Name=tag:Name,Values=pke-master --query "Reservations[*].Instances[*].PublicIpAddress" --output=text
```

To install the cluster, set the `CLUSTER_NAME` variable, and run:

```
export CLUSTER_NAME=""

export EXTERNAL_IP=$(curl http://169.254.169.254/latest/meta-data/public-ipv4)
export INTERNAL_IP=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)

export MAC=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/)
export VPC_CIDR=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/$MAC/vpc-ipv4-cidr-block/)

curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install master \
--kubernetes-advertise-address=${INTERNAL_IP} \
--kubernetes-api-server=${EXTERNAL_IP}:6443 \
--kubernetes-cloud-provider=aws \
--kubernetes-cluster-name=${CLUSTER_NAME} \
--kubernetes-infrastructure-cidr=${VPC_CIDR}

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

Please get the token and certhash from the logs or issue the following PKE command to print the token and cert hash needed by workers to join the cluster.

```
pke token list
```

If the tokens have expired or you'd like to create a new one, issue:

```
pke token create
```

#### Install and join a worker node

To install a worker node, run the following commands. Take note that you'd need to set the TOKEN and CERTHASH variables from above.

```
curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

# copy values from master node
export TOKEN=""
export CERTHASH=""
# same as $INTERNAL_IP from master node
export API_SERVER_INTERNAL_IP=""

pke install worker \
--kubernetes-node-token $TOKEN \
--kubernetes-api-server-ca-cert-hash $CERTHASH \
--kubernetes-api-server $API_SERVER_INTERNAL_IP:6443
```

You may also need to define storage classes for your cluster:
```
kubectl create -f - <<EOF 
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: gp2
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
EOF
```

For more information, see [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/) in the Kubernetes documentation.

Note that you can add as many worker nodes as you wish repeating the commands above. For basic troubleshooting you can check the status of the containers by issuing `crictl ps`.
