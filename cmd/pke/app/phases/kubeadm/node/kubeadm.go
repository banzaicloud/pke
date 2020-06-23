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

package node

import (
	"io"
	"net"
	"strings"
	"text/template"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/cri"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pbnjay/memory"
)

//go:generate templify -t ${GOTMPL} -p node -f kubeadmConfigV1Beta1 kubeadm_v1beta1.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p node -f kubeadmConfigV1Beta2 kubeadm_v1beta2.yaml.tmpl

func (n Node) writeKubeadmConfig(out io.Writer, filename string) error {
	// API server advertisement
	bindPort := "6443"
	if n.advertiseAddress != "" {
		host, port, err := kubeadm.SplitHostPort(n.advertiseAddress, "6443")
		if err != nil {
			return err
		}
		n.advertiseAddress = host
		bindPort = port
	}

	// Control Plane
	if n.apiServerHostPort != "" {
		host, port, err := kubeadm.SplitHostPort(n.apiServerHostPort, "6443")
		if err != nil {
			return err
		}
		n.apiServerHostPort = net.JoinHostPort(host, port)
	}

	ver, err := semver.NewVersion(n.kubernetesVersion)
	if err != nil {
		return errors.Wrapf(err, "unable to parse Kubernetes version %q", n.kubernetesVersion)
	}

	var conf string
	switch ver.Minor() {
	case 15, 16, 17:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		conf = kubeadmConfigV1Beta1Template()
	case 18:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		conf = kubeadmConfigV1Beta2Template()
	default:
		return errors.Errorf("unsupported Kubernetes version %q for kubeadm", n.kubernetesVersion)
	}

	tmpl, err := template.New("kubeadm-config").Parse(conf)
	if err != nil {
		return err
	}

	taints, err := kubernetes.ParseTaints(n.taints)
	if err != nil {
		return err
	}

	// kube reserved resources
	var (
		kubeReservedCPU    = "100m"
		kubeReservedMemory = kubeadm.KubeReservedMemory(memory.TotalMemory())
	)

	// Node labels
	nodeLabels := n.labels
	if n.nodepool != "" {
		nodeLabels = append(nodeLabels, "nodepool.banzaicloud.io/name="+n.nodepool)
	}

	type data struct {
		APIServerAdvertiseAddress string
		APIServerBindPort         string
		CRISocket                 string
		ControlPlaneEndpoint      string
		Token                     string
		CACertHash                string
		CloudProvider             string
		NodeLabels                string
		Taints                    []kubernetes.Taint
		KubeReservedCPU           string
		KubeReservedMemory        string
	}

	d := data{
		APIServerAdvertiseAddress: n.advertiseAddress,
		APIServerBindPort:         bindPort,
		CRISocket:                 cri.GetCRISocket(n.containerRuntime),
		ControlPlaneEndpoint:      n.apiServerHostPort,
		Token:                     n.kubeadmToken,
		CACertHash:                n.caCertHash,
		CloudProvider:             n.cloudProvider,
		NodeLabels:                strings.Join(nodeLabels, ","),
		Taints:                    taints,
		KubeReservedCPU:           kubeReservedCPU,
		KubeReservedMemory:        kubeReservedMemory,
	}

	return file.WriteTemplate(filename, tmpl, d)
}
