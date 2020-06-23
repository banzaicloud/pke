#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.18.4}"

# upgrade first master node
echo ""
echo "= centos1 ========================================================================"
vagrant ssh centos1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"

# upgrade second master node
echo ""
echo "= centos2 ========================================================================"
vagrant ssh centos2 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade third master node
echo ""
echo "= centos3 ========================================================================"
vagrant ssh centos3 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --kubernetes-additional-control-plane"

# upgrade worker node
echo ""
echo "= centos4 ========================================================================"
vagrant ssh centos4 -c "sudo /banzaicloud/pke upgrade worker --kubernetes-version='$KUBERNETES_VERSION'"
