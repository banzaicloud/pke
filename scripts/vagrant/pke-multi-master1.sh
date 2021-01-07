#!/bin/bash -xe

KUBERNETES_VERSION=$1
APISERVER_ADDRESS="${2:-192.168.64.11:6443}"
NETWORK_PROVIDER="${3:-cilium}"

systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke machine-image --kubernetes-version="$KUBERNETES_VERSION" && \

    /banzaicloud/pke install master \
      --kubernetes-master-mode=ha \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-advertise-address="${APISERVER_ADDRESS}" \
      --kubernetes-api-server="${APISERVER_ADDRESS}" \
      --kubernetes-network-provider="${NETWORK_PROVIDER}" && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)

tar -cvzf /banzaicloud/certs.tgz \
  /etc/kubernetes/pki/*ca.* \
  /etc/kubernetes/pki/ca.* \
  /etc/kubernetes/pki/etcd/*ca.* \
  /etc/kubernetes/pki/*sa*

cp /etc/kubernetes/admission-control/encryption-provider-config.yaml /banzaicloud/encryption-provider-config.yaml
