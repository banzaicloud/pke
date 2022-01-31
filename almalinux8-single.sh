#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.20.6}"

vagrant up almalinux1
vagrant ssh almalinux1 -c "sudo /scripts/pke-single.sh '$KUBERNETES_VERSION' '192.168.64.11:6443' containerd cilium"
vagrant ssh almalinux1 -c 'sudo cat /etc/kubernetes/admin.conf' > pke-single-config.yaml

export KUBECONFIG=$PWD/pke-single-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-single-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh almalinux1 -c 'sudo -s'"
