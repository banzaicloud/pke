#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.14.3}"

# upgrade first master node
echo ""
echo "= node1 ========================================================================"
vagrant ssh node1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"

# upgrade second master node
echo ""
echo "= node2 ========================================================================"
vagrant ssh node2 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade third master node
echo ""
echo "= node3 ========================================================================"
vagrant ssh node3 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade worker node
echo ""
echo "= node4 ========================================================================"
vagrant ssh node4 -c "sudo /banzaicloud/pke upgrade worker --kubernetes-version='$KUBERNETES_VERSION'"
