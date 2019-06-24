#!/bin/bash -e

KUBERNETES_VERSION=$1
TOKEN="$2"
CERTHASH="$3"


systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -v https://banzaicloud.com/downloads/pke/pke-0.4.9 -o /banzaicloud/pke && \
        chmod +x /usr/local/bin/pke
    ) && \

    /banzaicloud/pke machine-image --kubernetes-version="$KUBERNETES_VERSION" && \

    /banzaicloud/pke install worker \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-api-server=192.168.64.11:6443 \
      --kubernetes-node-token ${TOKEN} \
      --kubernetes-api-server-ca-cert-hash ${CERTHASH}
)
