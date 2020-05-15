#!/bin/bash -xe

KUBERNETES_VERSION=$1
APISERVER_ADDRESS="${2:-192.168.64.11:6443}"
CONTAINER_RUNTIME="${3:-containerd}"

systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke install single \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-container-runtime="${CONTAINER_RUNTIME}" \
      --kubernetes-advertise-address="${APISERVER_ADDRESS}" \
      --kubernetes-api-server="${APISERVER_ADDRESS}" && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
