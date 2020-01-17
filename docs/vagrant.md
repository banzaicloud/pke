## PKE in Vagrant

In order to *try out* a single/multi node PKE Kubernetes cluster you can also use Vagrant.

To install Vagrant on a Mac OS X follow these steps below:

- Install VirtualBox: `brew cask install virtualbox`
- Install Vagrant: `brew cask install vagrant`
- Install the vagrant-vbguest plugin `vagrant plugin install vagrant-vbguest`

> You may need to download VirtualBox 6.0 manually, because [VirtualBox 6.1 support](https://github.com/hashicorp/vagrant/pull/11250) is not yet released (as of early 2020).

You are set, and ready to install PKE on your machine. Now you should clone the PKE repo and follow these steps:

```
git clone git@github.com:banzaicloud/pke.git
cd pke
```

### Single node PKE

Start a machine with the following command: `vagrant up centos1`

Once the node is up follow these instructions:
```
vagrant ssh centos1 -c 'sudo -s'

curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install single
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

You can use `kubectl` from now on. Test it by executing `kubectl get nodes`

### Multi node PKE

Start a machine with the following command: `vagrant up centos1 centos2`

Once the node is up follow these instructions:

#### Start the master node

```
vagrant ssh centos1 -c 'sudo -s'

curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

pke install master --kubernetes-advertise-address=192.168.64.11 --kubernetes-api-server=192.168.64.11:6443 
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

Please get the token and certhash from the logs or issue the following PKE command to print the token and cert hash needed by workers to join the cluster.

```
pke token list # Print token and cert hash needed by workers to join the cluster
```

If the tokens have expired or you'd like to create a new one, issue:

```
pke token create
```

#### Start a worker node

Take note that you'd need to export the TOKEN and CERTHASH environment variables from above.

```
vagrant ssh centos2
sudo -s

curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/

# copy values from centos1
export TOKEN=""
export CERTHASH=""
pke install worker --kubernetes-node-token $TOKEN --kubernetes-api-server-ca-cert-hash $CERTHASH --kubernetes-api-server 192.168.64.11:6443
```

Note that you can add as many worker nodes as you wish repeating the commands above. You can check the status of the containers by issueing `crictl ps`
