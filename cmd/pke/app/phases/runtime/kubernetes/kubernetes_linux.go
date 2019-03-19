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

	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

const (
	k8sRepoFile = "/etc/yum.repos.d/kubernetes.repo"
	k8sRepo     = `[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kube*`

	selinuxConfig       = "/etc/selinux/config"
	kubeletKernelparams = "/etc/sysctl.d/90-kubelet.conf"
)

func (r *Runtime) installRuntime(out io.Writer, kubernetesVersion string) error {
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

	if err = linux.SwapOff(out); err != nil {
		return err
	}

	// modprobe nf_conntrack_ipv4
	if err := linux.Modprobe(out, "nf_conntrack_ipv4"); err != nil {
		return errors.Wrap(err, "missing nf_conntrack_ipv4 Linux Kernel module")
	}

	// modprobe ip_vs
	if err := linux.Modprobe(out, "ip_vs"); err != nil {
		return errors.Wrap(err, "missing ip_vs Linux Kernel module")
	}

	// modprobe ip_vs_rr
	if err := linux.Modprobe(out, "ip_vs_rr"); err != nil {
		return errors.Wrap(err, "missing ip_vs_rr Linux Kernel module")
	}

	// modprobe ip_vs_wrr
	if err := linux.Modprobe(out, "ip_vs_wrr"); err != nil {
		return errors.Wrap(err, "missing ip_vs_wrr Linux Kernel module")
	}

	// modprobe ip_vs_sh
	if err := linux.Modprobe(out, "ip_vs_sh"); err != nil {
		return errors.Wrap(err, "missing ip_vs_sh Linux Kernel module")
	}

	if err := writeKubeletKernelParams(out, kubeletKernelparams); err != nil {
		return errors.Wrapf(err, "unable to write kubernetes kernel params to %s", kubeletKernelparams)
	}

	if err := linux.SysctlLoadAllFiles(out); err != nil {
		return errors.Wrapf(err, "unable to load all sysctl rules from files")
	}

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
	err = file.Overwrite(k8sRepoFile, k8sRepo)
	if err != nil {
		return err
	}

	// Install kubelet kubeadm and kubectl.
	// yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
	if err := linux.YumInstall(out, yumPackages(kubernetesVersion)); err != nil {
		return errors.Wrap(err, "unable to install packages")
	}

	_ = linux.SystemctlDisableAndStop(out, "kubelet")

	return nil
}

func writeKubeletKernelParams(out io.Writer, filename string) error {
	dir := filepath.Dir(filename)

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, dir)
	err := os.MkdirAll(dir, 0640)
	if err != nil {
		return err
	}

	conf := `vm.overcommit_memory=1
kernel.panic=10
kernel.panic_on_oops=1
`

	return file.Overwrite(filename, conf)
}

func yumPackages(kubernetesVersion string) []string {
	return []string{
		"kubelet-" + kubernetesVersion + "-0",
		"kubeadm-" + kubernetesVersion + "-0",
		"kubectl-" + kubernetesVersion + "-0",
		"--disableexcludes=kubernetes",
	}
}
