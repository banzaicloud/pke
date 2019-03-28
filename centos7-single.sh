#!/bin/bash -e

vagrant up node1
vagrant ssh node1 -c "sudo /scripts/pke-single.sh"
vagrant ssh node1 -c 'sudo cat /etc/kubernetes/admin.conf' > pke-single-config.yaml

export KUBECONFIG=$PWD/pke-single-config.yaml

echo ""
echo "You can access your PKE cluster either:"
echo "- from your host machine accessing the cluster via kubectl. Please run:"
echo "export KUBECONFIG=$PWD/pke-single-config.yaml"
echo ""
echo "- or starting a shell in the virtual machine. Please run:"
echo "vagrant ssh node1 -c 'sudo -s'"
