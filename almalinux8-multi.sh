#!/bin/bash -e

# prerequisitesSkipping phase due to missing Pipeline API endpoint credentials
jq --version || (echo "Please install jq command line tool. https://stedolan.github.io/jq/" && exit 1)

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.22.6}"

# install first master node
echo ""
echo "= almalinux1 ========================================================================"
vagrant up almalinux1
vagrant ssh almalinux1 -c "sudo sh -c 'echo -n "LANG=en_US.utf-8\nLC_ALL=en_US.utf-8" > /etc/environment'"
vagrant ssh almalinux1 -c "sudo /scripts/pke-multi-master1.sh '$KUBERNETES_VERSION' '192.168.64.11:6443'"
vagrant ssh almalinux1 -c 'sudo cat /etc/kubernetes/admin.conf' > pke-multi-config.yaml
vagrant ssh almalinux1 -c "sudo /banzaicloud/pke token list -o json" > build/token.json


TOKEN=`jq -r '.tokens[] | select(.expired == false) | .token' build/token.json`
CERTHASH=`jq -r '.tokens[] | select(.expired == false) | .hash' build/token.json`

echo ""
echo "Using '$TOKEN' and '$CERTHASH' to join other nodes to the cluster"

# install second master node
echo ""
echo "= almalinux2 ========================================================================"
vagrant up almalinux2
vagrant ssh almalinux2 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '192.168.64.11:6443' '$TOKEN' '$CERTHASH' 192.168.64.12"

# install third master node
echo ""
echo "= almalinux3 ========================================================================"
vagrant up almalinux3
vagrant ssh almalinux3 -c "sudo /scripts/pke-multi-mastern.sh '$KUBERNETES_VERSION' '192.168.64.11:6443' '$TOKEN' '$CERTHASH' 192.168.64.13"

# install worker node
echo ""
echo "= almalinux4 ========================================================================"
vagrant up almalinux4
vagrant ssh almalinux4 -c "sudo /scripts/pke-multi-worker.sh '$KUBERNETES_VERSION' '192.168.64.11:6443' '$TOKEN' '$CERTHASH'"

export KUBECONFIG=$PWD/pke-multi-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-multi-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh almalinux1 -c 'sudo -s'"
