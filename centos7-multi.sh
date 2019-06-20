#!/bin/bash -e

# prerequisitesSkipping phase due to missing Pipeline API endpoint credentials
jq --version || (echo "Please install jq command line tool. https://stedolan.github.io/jq/" && exit 1)

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.14.3}"

# install first master node
echo ""
echo "= node1 ========================================================================"
vagrant up node1
vagrant ssh node1 -c "sudo /scripts/pke-multi-master1.sh $KUBERNETES_VERSION"
vagrant ssh node1 -c 'sudo cat /etc/kubernetes/admin.conf' > pke-multi-config.yaml
vagrant ssh node1 -c "sudo /banzaicloud/pke token list -o json" > build/token.json

TOKEN=`jq -r '.tokens[] | select(.expired == false) | .token' build/token.json`
CERTHASH=`jq -r '.tokens[] | select(.expired == false) | .hash' build/token.json`

echo ""
echo "Using '$TOKEN' and '$CERTHASH' to join other nodes to the cluster"

# install second master node
echo ""
echo "= node2 ========================================================================"
vagrant up node2
vagrant ssh node2 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '$TOKEN' '$CERTHASH' 192.168.64.12"

# install third master node
echo ""
echo "= node3 ========================================================================"
vagrant up node3
vagrant ssh node3 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '$TOKEN' '$CERTHASH' 192.168.64.13"

# install worker node
echo ""
echo "= node4 ========================================================================"
vagrant up node4
vagrant ssh node4 -c "sudo /scripts/pke-multi-worker.sh '$KUBERNETES_VERSION' '$TOKEN' '$CERTHASH'"

export KUBECONFIG=$PWD/pke-multi-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-multi-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh node1 -c 'sudo -s'"
