#!/bin/bash -xe

KUBERNETES_VERSION=$1

systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke machine-image --kubernetes-version="$KUBERNETES_VERSION" && \

    /banzaicloud/pke install master \
      --kubernetes-master-mode=ha \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-advertise-address=192.168.64.11:6443 \
      --kubernetes-api-server=192.168.64.11:6443 && \
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
