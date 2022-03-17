#!/bin/bash -xe

KUBERNETES_VERSION=$1
APISERVER_ADDRESS="${2:-192.168.64.11:6443}"
CONTAINER_RUNTIME="${3:-containerd}"
NETWORK_PROVIDER="${4:-cilium}"

systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://github.com/banzaicloud/pke/releases/download/0.9.0/pke-0.9.0 -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke install single \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-container-runtime="${CONTAINER_RUNTIME}" \
      --kubernetes-advertise-address="${APISERVER_ADDRESS}" \
      --kubernetes-api-server="${APISERVER_ADDRESS}" \
      --kubernetes-network-provider="${NETWORK_PROVIDER}" && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
