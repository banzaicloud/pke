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

package linux

import (
	"io"
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdYum             = "/bin/yum"
	banzaiCloudRPMRepo = "/etc/yum.repos.d/banzaicloud.repo"
	k8sRPMRepoFile     = "/etc/yum.repos.d/kubernetes.repo"
	k8sRPMRepo         = `[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*`
)

func yumErrorMatcher(text string) bool {
	return strings.Contains(strings.ToLower(text), "error") || strings.Contains(text, "No package ")
}

func YumInstall(out io.Writer, packages packages) error {
	cmd := runner.Cmd(out, cmdYum, append([]string{"install", "-y", disableExcludesKubernetes}, packages.strings()...)...)
	cmd.ErrorMatcher(yumErrorMatcher)
	text, err := cmd.CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = checkRPMPackages(out, packages)
	return errors.WrapIff(err, "yum installation failed [%s]", text)
}

var _ ContainerdPackages = (*YumInstaller)(nil)
var _ KubernetesPackages = (*YumInstaller)(nil)

type YumInstaller struct{}

func (y *YumInstaller) InstallKubernetesPrerequisites(out io.Writer, kubernetesVersion string) error {
	// Set SELinux in permissive mode (effectively disabling it)
	// setenforce 0
	err := runner.Cmd(out, "setenforce", "0").Run()
	if err != nil {
		return err
	}
	// sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config
	err = runner.Cmd(out, "sed", "-i", "s/^SELINUX=enforcing$/SELINUX=permissive/", selinuxConfig).Run()
	if err != nil {
		return err
	}

	if err = SwapOff(out); err != nil {
		return err
	}

	if err := ModprobeKubeProxyIPVSModules(out); err != nil {
		return err
	}

	if err := SysctlLoadAllFiles(out); err != nil {
		return errors.Wrapf(err, "unable to load all sysctl rules from files")
	}

	if _, err := os.Stat(banzaiCloudRPMRepo); err != nil {
		// Add kubernetes repo
		// cat <<EOF > /etc/yum.repos.d/kubernetes.repo
		// [kubernetes]
		// name=Kubernetes
		// baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
		// enabled=1
		// gpgcheck=1
		// repo_gpgcheck=1
		// gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
		// exclude=kube*
		// EOF
		err = file.Overwrite(k8sRPMRepoFile, k8sRPMRepo)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewYumInstaller() *YumInstaller {
	return &YumInstaller{}
}

func (y *YumInstaller) InstallKubernetesPackages(out io.Writer, kubernetesVersion string) error {
	// yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
	pkg := packages{{kubeadm, kubernetesVersion},
		{kubelet, kubernetesVersion},
		{kubectl, kubernetesVersion},
		{kubernetescni, kubernetesCNIVersion}}

	return YumInstall(out, pkg)
}

func (y *YumInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	// yum install -y kubeadm --disableexcludes=kubernetes
	pkg := []pkg{{kubeadm, kubernetesVersion},
		{kubelet, kubernetesVersion},
		{kubernetescni, kubernetesCNIVersion}}

	return YumInstall(out, pkg)
}

func (y *YumInstaller) InstallContainerdPrerequisites(out io.Writer, containerdVersion string) error {
	// yum install -y libseccomp
	if err := YumInstall(out, packages{{"libseccomp", ""}}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}
