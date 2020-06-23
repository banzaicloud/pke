#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.18.4}"

# upgrade first master node
echo ""
echo "= ubuntu1 ========================================================================"
vagrant ssh ubuntu1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"

# upgrade second master node
echo ""
echo "= ubuntu2 ========================================================================"
vagrant ssh ubuntu2 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade third master node
echo ""
echo "= ubuntu3 ========================================================================"
vagrant ssh ubuntu3 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade worker node
echo ""
echo "= ubuntu4 ========================================================================"
vagrant ssh ubuntu4 -c "sudo /banzaicloud/pke upgrade worker --kubernetes-version='$KUBERNETES_VERSION'"
