#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.21.0}"

vagrant ssh almalinux1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"
