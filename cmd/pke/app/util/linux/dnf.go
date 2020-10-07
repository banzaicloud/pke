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
	"fmt"
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
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg`
)

func DnfInstall(out io.Writer, packages []string) error {
	_, err := runner.Cmd(out, cmdDnf, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}
		if pkg[:1] == "-" {
			continue
		}

		name, ver, rel, arch, err := rpmQuery(out, pkg)
		if err != nil {
			return err
		}
		if name == pkg ||
			name+"-"+ver == pkg ||
			name+"-"+ver+"-"+rel == pkg ||
			name+"-"+ver+"-"+rel+"."+arch == pkg {
			continue
		}
		return errors.New(fmt.Sprintf("expected package version after installation: %q, got: %q", pkg, name+"-"+ver+"-"+rel+"."+arch))
	}

	return nil
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
	p := []string{
		"kubelet-" + kubernetesVersion,
		"kubeadm-" + kubernetesVersion,
		"kubectl-" + kubernetesVersion,
	}

	return DnfInstall(out, p)
}

func (y *DnfInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	// dnf install -y kubeadm --disableexcludes=kubernetes
	pkg := []string{
		"kubelet-" + kubernetesVersion,
		"kubeadm-" + kubernetesVersion,
	}

	return DnfInstall(out, pkg)
}

func (y *DnfInstaller) InstallContainerdPrerequisites(out io.Writer, containerdVersion string) error {
	// dnf install -y libseccomp
	if err := DnfInstall(out, []string{"libseccomp"}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}
