#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.18.4}"

vagrant up ubuntu-docker
vagrant ssh ubuntu-docker -c "sudo /scripts/pke-single.sh '$KUBERNETES_VERSION' '192.168.64.30:6443' docker"
vagrant ssh ubuntu-docker -c 'sudo cat /etc/kubernetes/admin.conf' > pke-single-config.yaml

export KUBECONFIG=$PWD/pke-single-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-single-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh ubuntu-docker -c 'sudo -s'"
