---
title: pke machine-image
generated_file: true
---
## pke machine-image

Machine image build helper for Banzai Cloud Pipeline Kubernetes Engine (PKE)

### Synopsis

Machine image build helper for Banzai Cloud Pipeline Kubernetes Engine (PKE)

```
pke machine-image [flags]
```

### Options

```
  -h, --help                                  help for machine-image
      --image-repository string               Prefix for image repository (default "banzaicloud")
      --kubernetes-container-runtime string   Kubernetes container runtime (default "containerd")
      --kubernetes-version string             Kubernetes version (default "1.19.10")
      --use-image-repo-for-k8s                Use defined image repository for K8s Images as well
```

### SEE ALSO

* [pke](/docs/pke/cli/reference/pke/)	 - Bootstrap a secure Kubernetes cluster with Banzai Cloud Pipeline Kubernetes Engine (PKE)
* [pke machine-image container-runtime](/docs/pke/cli/reference/pke_machine-image_container-runtime/)	 - Container runtime installation
* [pke machine-image image-pull](/docs/pke/cli/reference/pke_machine-image_image-pull/)	 - Pull images used by PKE tool
* [pke machine-image kubernetes-runtime](/docs/pke/cli/reference/pke_machine-image_kubernetes-runtime/)	 - Kubernetes runtime installation
* [pke machine-image kubernetes-version](/docs/pke/cli/reference/pke_machine-image_kubernetes-version/)	 - Check Kubernetes version is supported or not
* [pke machine-image write-config](/docs/pke/cli/reference/pke_machine-image_write-config/)	 - Write configuration file

