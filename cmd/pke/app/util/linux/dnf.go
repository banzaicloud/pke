// Copyright © 2020 Banzai Cloud
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

package linux

import (
	"io"
	"os"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdDnf        = "/bin/dnf"
	k8sRPMRepoDnf = `[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
excludepkgs=[kube*]`
)

func DnfInstall(out io.Writer, packages packages) error {
	_, err := runner.Cmd(out, cmdDnf, append([]string{"install", "-y", disableExcludesKubernetes}, packages.strings()...)...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	return checkRPMPackages(out, packages)
}

var _ ContainerdPackages = (*DnfInstaller)(nil)
var _ KubernetesPackages = (*DnfInstaller)(nil)

type DnfInstaller struct{}

func NewDnfInstaller() *DnfInstaller {
	return &DnfInstaller{}
}

func (y *DnfInstaller) InstallKubernetesPrerequisites(out io.Writer, kubernetesVersion string) error {
	if err := SwapOff(out); err != nil {
		return err
	}

	if err := ModprobeKubeProxyIPVSModules(out); err != nil {
		return err
	}

	if err := SysctlLoadAllFiles(out); err != nil {
		return errors.Wrapf(err, "unable to load all sysctl rules from files")
	}

	if _, err := os.Stat(banzaiCloudRPMRepo); err != nil {
		err = file.Overwrite(k8sRPMRepoFile, k8sRPMRepoDnf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (y *DnfInstaller) InstallKubernetesPackages(out io.Writer, kubernetesVersion string) error {
	// dnf install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
	pkg := packages{
		{"kubelet", kubernetesVersion},
		{"kubeadm", kubernetesVersion},
		{"kubectl", kubernetesVersion},
	}

	return DnfInstall(out, pkg)
}

func (y *DnfInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	// dnf install -y kubeadm --disableexcludes=kubernetes
	pkg := packages{
		{"kubelet", kubernetesVersion},
		{"kubeadm", kubernetesVersion},
	}

	return DnfInstall(out, pkg)
}

func (y *DnfInstaller) InstallContainerdPrerequisites(out io.Writer, containerdVersion string) error {
	// dnf install -y libseccomp
	if err := DnfInstall(out, packages{{"libseccomp", ""}}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}
