## Requirements

### Operating system

`pke` currently is available for CentOS 7.x and RHEL 7.x. There is work in progress to support Ubuntu 18.04 LTS.

### Network

A flat network between nodes is required. Port `6443` (K8s API server) should be opened if there is need to access K8s API externally.

### The `pke` binary

You can download a particular binary release from the project's release page on [GitHub](https://github.com/banzaicloud/pke/releases). Our guides assume that the executable is available as `pke` in the system PATH.

You can also use the following commands as root to achieve this:


```
curl -v https://banzaicloud.com/downloads/pke/pke-0.4.5 -o /usr/local/bin/pke
chmod +x /usr/local/bin/pke
export PATH=$PATH:/usr/local/bin/
```
