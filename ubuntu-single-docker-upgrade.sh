#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${2:-v1.22.6}"
UBUNTU_VERSION=${1:-focal}

vagrant ssh ubuntu-docker-${UBUNTU_VERSION} -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"
