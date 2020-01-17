## PKE in Vagrant

To *try out* a single- or multi-node PKE Kubernetes cluster you can also use Vagrant.

### Install Vagrant

To install Vagrant on a macOS, complete the following steps:

1. Install VirtualBox: `brew cask install virtualbox`
    > You may need to download VirtualBox 6.0 manually, because [VirtualBox 6.1 support](https://github.com/hashicorp/vagrant/pull/11250) of Vagrant is not yet released (as of early 2020).
1. Install Vagrant: `brew cask install vagrant`
1. Install the vagrant-vbguest plugin `vagrant plugin install vagrant-vbguest`
1. Clone the PKE repository. If you do not have git installed, you can also [download and unzip it into the pke directory](https://github.com/banzaicloud/pke/releases/latest):

    ```bash
    git clone git@github.com:banzaicloud/pke.git
    cd pke
    ```

1. Decide whether you want to install a single-node or a multi-node cluster. 
    * To start a single-node cluster, proceed to (Single node PKE)[#single-node-pke].
    * To start a multi-node cluster, proceed to (Multi node PKE)[#single-node-pke].

### Single node PKE

1. Start a machine with the following command: `vagrant up centos1`
1. Wait until the node starts.
1. Issue the following commands:

    ```bash
    vagrant ssh centos1 -c 'sudo -s'

    curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
    chmod +x /usr/local/bin/pke
    export PATH=$PATH:/usr/local/bin/

    pke install single
    mkdir -p $HOME/.kube
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    chown $(id -u):$(id -g) $HOME/.kube/config
    ```
    
1. You can use `kubectl` from now on. Execute the following command to test it: `kubectl get nodes`

### Multi node PKE

1. Start a machine with the following command: `vagrant up centos1 centos2`
1. Wait until the node starts.
1. Issue the following commands to start the master node:

    ```bash
    vagrant ssh centos1 -c 'sudo -s'

    curl -vL https://banzaicloud.com/downloads/pke/latest -o /usr/local/bin/pke
    chmod +x /usr/local/bin/pke
    export PATH=$PATH:/usr/local/bin/

    pke install master --kubernetes-advertise-address=192.168.64.11 --kubernetes-api-server=192.168.64.11:6443 
    mkdir -p $HOME/.kube
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    chown $(id -u):$(id -g) $HOME/.kube/config
    ```
    
1. Get the token and certhash from the logs, or issue the following PKE command to print the token and cert hash needed by workers to join the cluster.

    ```bash
    pke token list # Print token and cert hash needed by workers to join the cluster
    ```

    > Note: If the tokens have expired or you'd like to create a new one, issue:
    > 
    >    ```bash 
    >    pke token create
    >    ```

1. Export the TOKEN and CERTHASH environment variables with the values you have retrieved in the previous step.
1. Start a worker node using the following commands:

    ```bash
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

1. Repeat the previous step to add additional worker nodes. You can check the status of the containers by issuing the `crictl ps` command.
