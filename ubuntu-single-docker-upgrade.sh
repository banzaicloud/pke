#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.18.2}"

vagrant ssh ubuntu-docker -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION' --container-runtime=docker"
