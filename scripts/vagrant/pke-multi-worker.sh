#!/bin/bash -e

KUBERNETES_VERSION=$1
APISERVER_ADDRESS="${2:-192.168.64.11:6443}"
TOKEN="$3"
CERTHASH="$4"


systemctl is-active kubelet || ( \
    /banzaicloud/pke version -o yaml || ( \
        curl -vL https://banzaicloud.com/downloads/pke/latest -o /banzaicloud/pke && \
        chmod +x /usr/local/bin/pke
    ) && \

    /banzaicloud/pke machine-image --kubernetes-version="$KUBERNETES_VERSION" && \

    /banzaicloud/pke install worker \
      --kubernetes-version="${KUBERNETES_VERSION}" \
      --kubernetes-api-server="${APISERVER_ADDRESS}" \
      --kubernetes-node-token ${TOKEN} \
      --kubernetes-api-server-ca-cert-hash ${CERTHASH}
)
