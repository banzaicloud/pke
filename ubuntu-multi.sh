#!/bin/bash -e

# prerequisitesSkipping phase due to missing Pipeline API endpoint credentials
jq --version || (echo "Please install jq command line tool. https://stedolan.github.io/jq/" && exit 1)

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.18.4}"

# install first master node
echo ""
echo "= ubuntu1 ========================================================================"
vagrant up ubuntu1
vagrant ssh ubuntu1 -c "sudo /scripts/pke-multi-master1.sh '$KUBERNETES_VERSION' '192.168.64.21:6443'"
vagrant ssh ubuntu1 -c 'sudo cat /etc/kubernetes/admin.conf' > pke-multi-config.yaml
vagrant ssh ubuntu1 -c "sudo /banzaicloud/pke token list -o json" > build/token.json

TOKEN=`jq -r '.tokens[] | select(.expired == false) | .token' build/token.json`
CERTHASH=`jq -r '.tokens[] | select(.expired == false) | .hash' build/token.json`

echo ""
echo "Using '$TOKEN' and '$CERTHASH' to join other nodes to the cluster"

# install second master node
echo ""
echo "= ubuntu2 ========================================================================"
vagrant up ubuntu2
vagrant ssh ubuntu2 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '192.168.64.21:6443' '$TOKEN' '$CERTHASH' 192.168.64.22"

# install third master node
echo ""
echo "= ubuntu3 ========================================================================"
vagrant up ubuntu3
vagrant ssh ubuntu3 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '192.168.64.21:6443' '$TOKEN' '$CERTHASH' 192.168.64.23"

# install worker node
echo ""
echo "= ubuntu4 ========================================================================"
vagrant up ubuntu4
vagrant ssh ubuntu4 -c "sudo /scripts/pke-multi-worker.sh '$KUBERNETES_VERSION' '192.168.64.21:6443' '$TOKEN' '$CERTHASH'"

export KUBECONFIG=$PWD/pke-multi-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-multi-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh ubuntu1 -c 'sudo -s'"
