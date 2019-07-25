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
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/kubernetes"
	"github.com/pkg/errors"
)

//go:generate templify -t ${GOTMPL} -p node -f kubeadmConfigV1Alpha3 kubeadm_v1alpha3.yaml.tmpl
//go:generate templify -t ${GOTMPL} -p node -f kubeadmConfigV1Beta1 kubeadm_v1beta1.yaml.tmpl

func (n Node) writeKubeadmConfig(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

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
	case 12, 13:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3
		conf = kubeadmConfigV1Alpha3Template()
	case 14, 15:
		// see https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1
		conf = kubeadmConfigV1Beta1Template()
	default:
		return errors.Errorf("unsupported Kubernetes version %q for kubeadm", n.kubernetesVersion)
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

	taints, err := kubernetes.ParseTaints(n.taints)
	if err != nil {
		return err
	}

	type data struct {
		APIServerAdvertiseAddress string
		APIServerBindPort         string
		ControlPlaneEndpoint      string
		Token                     string
		CACertHash                string
		CloudProvider             string
		Nodepool                  string
		Taints                    []kubernetes.Taint
	}

	d := data{
		APIServerAdvertiseAddress: n.advertiseAddress,
		APIServerBindPort:         bindPort,
		ControlPlaneEndpoint:      n.apiServerHostPort,
		Token:                     n.kubeadmToken,
		CACertHash:                n.caCertHash,
		CloudProvider:             n.cloudProvider,
		Nodepool:                  n.nodepool,
		Taints:                    taints,
	}

	return tmpl.Execute(w, d)
}
