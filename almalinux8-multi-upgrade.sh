#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.23.3}"

# upgrade first master node
echo ""
echo "= almalinux1 ========================================================================"
vagrant ssh almalinux1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"

# waiting 10 seconds because of apiserver
sleep 10

# upgrade second master node
echo ""
echo "= almalinux2 ========================================================================"
vagrant ssh almalinux2 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade third master node
echo ""
echo "= almalinux3 ========================================================================"
vagrant ssh almalinux3 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade worker node
echo ""
echo "= almalinux4 ========================================================================"
vagrant ssh almalinux4 -c "sudo /banzaicloud/pke upgrade worker --kubernetes-version='$KUBERNETES_VERSION'"
