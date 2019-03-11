## Requirements

### Operating system

`pke` currently is available for CentOS 7.x and RHEL 7.x. There is work in progress to support Ubuntu 18.04 LTS.

### Network

A flat network between nodes is required. Port `6443` (K8s API server) should be opened if there is need to access K8s API externally.

### The `pke` binary

Check get a particular release please follow the release page on [GitHub](https://github.com/banzaicloud/pke/releases)


```
curl -v https://banzaicloud.com/downloads/pke/pke-0.1.0 -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/
```
