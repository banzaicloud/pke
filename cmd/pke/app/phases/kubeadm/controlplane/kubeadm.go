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
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pkg/errors"
)

//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Alpha3 kubeadm_v1alpha3.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p controlplane -f kubeadmConfigV1Beta1 kubeadm_v1beta1.yaml.tmpl

func (c ControlPlane) WriteKubeadmConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

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
	if ver.LessThan(semver.MustParse("1.13.0")) {
		encryptionProviderPrefix = "experimental-"
	}

	var conf string
	switch ver.Minor() {
	case 12, 13:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
		conf = kubeadmConfigV1Alpha3Template()
	case 14, 15:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		conf = kubeadmConfigV1Beta1Template()
	default:
		return errors.New(fmt.Sprintf("unsupported Kubernetes version %q for kubeadm", c.kubernetesVersion))
	}

	tmpl, err := template.New("kubeadm-config").Parse(conf)
	if err != nil {
		return err
	}

	// create and truncate write only file
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	taints, err := kubernetes.ParseTaints(c.taints)
	if err != nil {
		return err
	}

	type data struct {
		APIServerAdvertiseAddress   string
		APIServerBindPort           string
		ControlPlaneEndpoint        string
		APIServerCertSANs           []string
		KubeletCertificateAuthority string
		AdmissionConfig             string
		ClusterName                 string
		KubernetesVersion           string
		ServiceCIDR                 string
		PodCIDR                     string
		CloudProvider               string
		Nodepool                    string
		ControllerManagerSigningCA  string
		OIDCIssuerURL               string
		OIDCClientID                string
		ImageRepository             string
		EncryptionProviderPrefix    string
		WithPluginPSP               bool
		WithAuditLog                bool
		Taints                      []kubernetes.Taint
		AuditLogDir                 string
		AuditPolicyFile             string
		EtcdEndpoints               []string
		EtcdCAFile                  string
		EtcdCertFile                string
		EtcdKeyFile                 string
		EtcdPrefix                  string
	}

	d := data{
		APIServerAdvertiseAddress:   c.advertiseAddress,
		APIServerBindPort:           bindPort,
		ControlPlaneEndpoint:        c.apiServerHostPort,
		APIServerCertSANs:           c.apiServerCertSANs,
		KubeletCertificateAuthority: c.kubeletCertificateAuthority,
		AdmissionConfig:             admissionConfig,
		ClusterName:                 c.clusterName,
		KubernetesVersion:           c.kubernetesVersion,
		ServiceCIDR:                 c.serviceCIDR,
		PodCIDR:                     c.podNetworkCIDR,
		CloudProvider:               c.cloudProvider,
		Nodepool:                    c.nodepool,
		ControllerManagerSigningCA:  c.controllerManagerSigningCA,
		OIDCIssuerURL:               c.oidcIssuerURL,
		OIDCClientID:                c.oidcClientID,
		ImageRepository:             c.imageRepository,
		EncryptionProviderPrefix:    encryptionProviderPrefix,
		WithPluginPSP:               c.withPluginPSP,
		WithAuditLog:                c.withAuditLog,
		Taints:                      taints,
		AuditLogDir:                 auditLogDir,
		AuditPolicyFile:             auditPolicyFile,
		EtcdEndpoints:               c.etcdEndpoints,
		EtcdCAFile:                  c.etcdCAFile,
		EtcdCertFile:                c.etcdCertFile,
		EtcdKeyFile:                 c.etcdKeyFile,
		EtcdPrefix:                  c.etcdPrefix,
	}

	return tmpl.Execute(w, d)
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
