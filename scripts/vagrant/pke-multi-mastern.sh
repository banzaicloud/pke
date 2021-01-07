#!/bin/bash -xe

KUBERNETES_VERSION=$1
APISERVER_ADDRESS="${2:-192.168.64.11:6443}"
TOKEN="$3"
CERTHASH="$4"
ADVERTISE_ADDRESS="$5"
NETWORK_PROVIDER="${6:-cilium}"

systemctl is-active kubelet || ( \
    tar -xvzf /banzaicloud/certs.tgz -C / && \

    mkdir -p /etc/kubernetes/admission-control/ && \
    cp /banzaicloud/encryption-provider-config.yaml /etc/kubernetes/admission-control/encryption-provider-config.yaml && \

    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke machine-image --kubernetes-version="$KUBERNETES_VERSION" && \

    /banzaicloud/pke install master \
      --kubernetes-master-mode=ha \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-advertise-address="${ADVERTISE_ADDRESS}:6443" \
      --kubernetes-api-server="${APISERVER_ADDRESS}" \
      --kubernetes-node-token "${TOKEN}" \
      --kubernetes-api-server-ca-cert-hash "${CERTHASH}" \
      --kubernetes-join-control-plane \
      --kubernetes-network-provider="${NETWORK_PROVIDER}" && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
