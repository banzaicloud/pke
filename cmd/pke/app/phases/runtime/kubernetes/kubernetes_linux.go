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

package kubernetes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
)

const (
	kubeletKernelparams = "/etc/sysctl.d/90-kubelet.conf"
)

func (r *Runtime) installRuntime(out io.Writer) error {
	pm, err := linux.KubernetesPackagesImpl(out)
	if err != nil {
		return err
	}
	return install(out, r.kubernetesVersion, pm)
}

func install(out io.Writer, kubernetesVersion string, pm linux.KubernetesPackages) error {

	if err := writeKubeletKernelParams(out, kubeletKernelparams); err != nil {
		return errors.Wrapf(err, "unable to write kubernetes kernel params to %s", kubeletKernelparams)
	}

	if err := pm.InstallKubernetesPrerequisites(out, kubernetesVersion); err != nil {
		return err
	}

	// Install kubelet kubeadm and kubectl.
	if err := pm.InstallKubernetesPackages(out, kubernetesVersion); err != nil {
		return errors.Wrap(err, "unable to install packages")
	}

	_ = linux.SystemctlDisableAndStop(out, "kubelet")

	return nil
}

//go:generate templify -t ${GOTMPL} -p kubernetes -f kubeletKernelParams kubelet_kernel_params.yaml.tmpl

func writeKubeletKernelParams(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}

	return file.Overwrite(filename, kubeletKernelParamsTemplate())
}
