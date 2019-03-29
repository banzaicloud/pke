#!/bin/bash -xe

/usr/local/bin/pke version -o yaml || ( \
    curl -v https://banzaicloud.com/downloads/pke/pke-0.2.3 -o /usr/local/bin/pke && \
    chmod +x /usr/local/bin/pke
)

systemctl is-active kubelet || ( \
    /usr/local/bin/pke install single --kubernetes-advertise-address=192.168.64.11:6443 --kubernetes-api-server=192.168.64.11:6443 && \
    mkdir -p $HOME/.kube && \
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
    chown $(id -u):$(id -g) $HOME/.kube/config
)
