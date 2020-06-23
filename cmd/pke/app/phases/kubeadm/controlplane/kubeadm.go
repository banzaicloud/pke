// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/cri"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pbnjay/memory"
)

//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Beta1 kubeadm_v1beta1.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Beta2 kubeadm_v1beta2.yaml.tmpl

func (c ControlPlane) WriteKubeadmConfig(out io.Writer, filename string) error {
	// API server advertisement
	bindPort := "6443"
	if c.advertiseAddress != "" {
		host, port, err := kubeadm.SplitHostPort(c.advertiseAddress, "6443")
		if err != nil {
			return err
		}
		c.advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if c.apiServerHostPort != "" {
		host, port, err := kubeadm.SplitHostPort(c.apiServerHostPort, "6443")
		if err != nil {
			return err
		}
		c.apiServerHostPort = net.JoinHostPort(host, port)
	}

	ver, err := semver.NewVersion(c.kubernetesVersion)
	if err != nil {
		return errors.Wrapf(err, "unable to parse Kubernetes version %q", c.kubernetesVersion)
	}

	encryptionProviderPrefix := ""

	var conf string
	switch ver.Minor() {
	case 15, 16, 17:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		conf = kubeadmConfigV1Beta1Template()
	case 18:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2
		conf = kubeadmConfigV1Beta2Template()
	default:
		return errors.New(fmt.Sprintf("unsupported Kubernetes version %q for kubeadm", c.kubernetesVersion))
	}

	tmpl, err := template.New("kubeadm-config").Parse(conf)
	if err != nil {
		return err
	}

	taints, err := kubernetes.ParseTaints(c.taints)
	if err != nil {
		return err
	}

	// Cloud provider configuration
	var (
		cloudConfig        bool // Cloud config for: kube-apiserver and kube-controller-manager
		kubeletCloudConfig bool // Cloud config for: kubelet
	)
	// Cloud provider configuration file is provided or not
	cloudConfig = c.cloudProvider != "" && c.cloudProvider != constants.CloudProviderExternal
	// Decide if kubelet needs the cloud-provider config
	switch c.cloudProvider {
	case constants.CloudProviderAzure, constants.CloudProviderVsphere:
		kubeletCloudConfig = true
	}

	// Kube reserved resources
	var (
		kubeReservedCPU    = "100m"
		kubeReservedMemory = kubeadm.KubeReservedMemory(memory.TotalMemory())
	)

	// Node labels
	nodeLabels := c.labels
	if c.nodepool != "" {
		nodeLabels = append(nodeLabels, "nodepool.banzaicloud.io/name="+c.nodepool)
	}

	type data struct {
		APIServerAdvertiseAddress       string
		APIServerBindPort               string
		CRISocket                       string
		ControlPlaneEndpoint            string
		APIServerCertSANs               []string
		KubeletCertificateAuthority     string
		AdmissionConfig                 string
		ClusterName                     string
		KubernetesVersion               string
		UseHyperKubeImage               bool
		ServiceCIDR                     string
		PodCIDR                         string
		CloudProvider                   string
		CloudConfig                     bool
		KubeletCloudConfig              bool
		NodeLabels                      string
		ControllerManagerSigningCA      string
		OIDCIssuerURL                   string
		OIDCClientID                    string
		ImageRepository                 string
		EncryptionProviderPrefix        string
		WithPluginPSP                   bool
		WithoutPluginDenyEscalatingExec bool
		WithAuditLog                    bool
		Taints                          []kubernetes.Taint
		AuditLogDir                     string
		AuditPolicyFile                 string
		EtcdEndpoints                   []string
		EtcdCAFile                      string
		EtcdCertFile                    string
		EtcdKeyFile                     string
		EtcdPrefix                      string
		KubeReservedCPU                 string
		KubeReservedMemory              string
	}

	d := data{
		APIServerAdvertiseAddress:       c.advertiseAddress,
		APIServerBindPort:               bindPort,
		CRISocket:                       cri.GetCRISocket(c.containerRuntime),
		ControlPlaneEndpoint:            c.apiServerHostPort,
		APIServerCertSANs:               c.apiServerCertSANs,
		KubeletCertificateAuthority:     c.kubeletCertificateAuthority,
		AdmissionConfig:                 admissionConfig,
		ClusterName:                     c.clusterName,
		KubernetesVersion:               c.kubernetesVersion,
		UseHyperKubeImage:               c.useHyperKubeImage,
		ServiceCIDR:                     c.serviceCIDR,
		PodCIDR:                         c.podNetworkCIDR,
		CloudProvider:                   c.cloudProvider,
		CloudConfig:                     cloudConfig,
		KubeletCloudConfig:              kubeletCloudConfig,
		NodeLabels:                      strings.Join(nodeLabels, ","),
		ControllerManagerSigningCA:      c.controllerManagerSigningCA,
		OIDCIssuerURL:                   c.oidcIssuerURL,
		OIDCClientID:                    c.oidcClientID,
		ImageRepository:                 c.imageRepository,
		EncryptionProviderPrefix:        encryptionProviderPrefix,
		WithPluginPSP:                   c.withPluginPSP,
		WithoutPluginDenyEscalatingExec: c.withoutPluginDenyEscalatingExec,
		WithAuditLog:                    !c.withoutAuditLog,
		Taints:                          taints,
		AuditLogDir:                     auditLogDir,
		AuditPolicyFile:                 auditPolicyFile,
		EtcdEndpoints:                   c.etcdEndpoints,
		EtcdCAFile:                      c.etcdCAFile,
		EtcdCertFile:                    c.etcdCertFile,
		EtcdKeyFile:                     c.etcdKeyFile,
		EtcdPrefix:                      c.etcdPrefix,
		KubeReservedCPU:                 kubeReservedCPU,
		KubeReservedMemory:              kubeReservedMemory,
	}

	return file.WriteTemplate(filename, tmpl, d)
}

//go:generate templify -t ${GOTMPL} -p controlplane -f auditV1Beta1 audit_v1beta1.yaml.tmpl

func writeAuditPolicyFile(out io.Writer) error {
	filename := auditPolicyFile
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	err = file.Overwrite(filename, auditV1Beta1Template())
	if err != nil {
		return err
	}
	return nil
}
