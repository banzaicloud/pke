#!/bin/bash -xe

/banzaicloud/pke version -o yaml || ( \
    curl -v https://banzaicloud.com/downloads/pke/pke-0.4.8 -o /banzaicloud/pke && \
    chmod +x /banzaicloud/pke
)

/banzaicloud/pke version -o yaml

systemctl is-active kubelet || ( \
    /banzaicloud/pke install single --kubernetes-advertise-address=192.168.64.11:6443 --kubernetes-api-server=192.168.64.11:6443 && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
