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
	"io/ioutil"
	"net/url"
	"os"

	"emperror.dev/errors"
	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdApt             = "/usr/bin/apt-get"
	cmdAptKey          = "/usr/bin/apt-key"
	banzaiCloudDEBRepo = "/etc/apt/sources.list.d/banzaicloud.repo"
	k8sDEBRepoFile     = "/etc/apt/sources.list.d/kubernetes.list"
	k8sDEBRepo         = `deb https://apt.kubernetes.io/ kubernetes-xenial main`
	k8sDEBRepoGPG      = "https://packages.cloud.google.com/apt/doc/apt-key.gpg"
)

var _ ContainerdPackages = (*AptInstaller)(nil)
var _ KubernetesPackages = (*AptInstaller)(nil)

type AptInstaller struct{}

func NewAptInstaller() *AptInstaller {
	return &AptInstaller{}
}

func (a *AptInstaller) InstallKubernetesPrerequisites(out io.Writer, kubernetesVersion string) error {
	if err := SwapOff(out); err != nil {
		return err
	}

	if err := ModprobeKubeProxyIPVSModules(out); err != nil {
		return err
	}

	if err := SysctlLoadAllFiles(out); err != nil {
		return errors.Wrapf(err, "unable to load all sysctl rules from files")
	}

	if _, err := os.Stat(banzaiCloudDEBRepo); err != nil {
		// Add kubernetes repo
		err = file.Overwrite(k8sDEBRepoFile, k8sDEBRepo)
		if err != nil {
			return err
		}
	}

	// curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
	// Download Kubernetes repo apt key.
	f, err := ioutil.TempFile("", "kubernetes-apt-key")
	if err != nil {
		return errors.Wrapf(err, "unable to create temporary file: %q", f.Name())
	}
	defer func() { _ = f.Close() }()
	u, err := url.Parse(k8sDEBRepoGPG)
	if err != nil {
		return errors.Wrapf(err, "unable to parse Kubernetes repo apt key. url: %q", k8sDEBRepoGPG)
	}
	if err = file.Download(u, f.Name()); err != nil {
		return errors.Wrapf(err, "unable to download Kubernetes repo apt key. url: %q", u.String())
	}
	if _, err := runner.Cmd(out, cmdAptKey, "add", f.Name()).CombinedOutputAsync(); err != nil {
		return errors.Wrap(err, "unable to add Kubernetes repo apt key")
	}

	_, err = runner.Cmd(out, cmdApt, "update").CombinedOutputAsync()
	return err
}

func (a *AptInstaller) InstallKubernetesPackages(out io.Writer, kubernetesVersion string) error {
	p := []string{
		mapAptPackageVersion(kubelet, kubernetesVersion),
		mapAptPackageVersion(kubeadm, kubernetesVersion),
		mapAptPackageVersion(kubectl, kubernetesVersion),
		mapAptPackageVersion(kubernetescni, kubernetesVersion),
	}

	return AptInstall(out, p)
}

func (a *AptInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	p := []string{
		mapAptPackageVersion(kubeadm, kubernetesVersion),
		mapAptPackageVersion(kubelet, kubernetesVersion),       // kubeadm dependency
		mapAptPackageVersion(kubernetescni, kubernetesVersion), // kubeadm dependency
	}

	return AptInstall(out, p)
}

func (a *AptInstaller) InstallContainerdPrerequisites(out io.Writer, containerdVersion string) error {
	// apt-get install -y libseccomp
	if err := AptInstall(out, []string{"libseccomp2"}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}

func AptInstall(out io.Writer, packages []string) error {
	_, err := runner.Cmd(out, cmdApt, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
	return err
}

func mapAptPackageVersion(pkg, kubernetesVersion string) string {
	switch pkg {
	case kubeadm:
		return "kubeadm=" + getAptPackageVersion(kubernetesVersion)

	case kubectl:
		return "kubectl=" + getAptPackageVersion(kubernetesVersion)

	case kubelet:
		return "kubelet=" + getAptPackageVersion(kubernetesVersion)

	case kubernetescni:
		return "kubernetes-cni=0.8.6-00"

	default:
		return ""
	}
}

func getAptPackageVersion(kubernetesVersion string) string {
	ver, _ := semver.NewVersion(kubernetesVersion)
	// There was an issue with bundled CNI plugin so new package was released in case of versions below. (https://github.com/kubernetes/kubernetes/issues/92242)
	c, _ := semver.NewConstraint("=1.16.11 || =1.17.7 || =1.18.4")
	if c.Check(ver) {
		return kubernetesVersion + "-01"
	}
	return kubernetesVersion + "-00"
}
