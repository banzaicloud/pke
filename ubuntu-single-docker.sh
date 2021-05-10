#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${2:-v1.19.10}"
UBUNTU_VERSION=${1:-focal}

vagrant up ubuntu-docker-${UBUNTU_VERSION}
vagrant ssh ubuntu-docker-${UBUNTU_VERSION} -c "sudo /scripts/pke-single.sh '$KUBERNETES_VERSION' '192.168.64.30:6443' docker"
vagrant ssh ubuntu-docker-${UBUNTU_VERSION} -c 'sudo cat /etc/kubernetes/admin.conf' > pke-single-config.yaml

export KUBECONFIG=$PWD/pke-single-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-single-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh ubuntu-docker-${UBUNTU_VERSION} -c 'sudo -s'"
