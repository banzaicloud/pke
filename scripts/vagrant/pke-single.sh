#!/bin/bash -xe

KUBERNETES_VERSION=$1

systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /banzaicloud/pke
    ) && \

    /banzaicloud/pke install single \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-advertise-address=192.168.64.11:6443 \
      --kubernetes-api-server=192.168.64.11:6443 && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
