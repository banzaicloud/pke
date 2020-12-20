#!/bin/bash -e

# build latest pke tool
GOOS=linux make pke

KUBERNETES_VERSION="${1:-v1.19.6}"
export VAGRANT_VAGRANTFILE=Vagrantfile-centos8

vagrant ssh centos1 -c "sudo /banzaicloud/pke upgrade master --kubernetes-version='$KUBERNETES_VERSION'"
