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
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

const (
	cmdYum                    = "/bin/yum"
	cmdRpm                    = "/bin/rpm"
	kubeadm                   = "kubeadm"
	kubectl                   = "kubectl"
	kubelet                   = "kubelet"
	kubernetescni             = "kubernetes-cni"
	disableExcludesKubernetes = "--disableexcludes=kubernetes"

	dotS  = "."
	dashS = "-"
)

var (
	errorUnableToParseRPMOutput = errors.New("Unable to parse rpm output")
	errorNotInstalled           = errors.New("package is not installed")
)

func YumInstall(out io.Writer, packages []string) error {
	// check packages are already installed or not
	n := 0
	p := make([]string, 0)
	for _, pkg := range packages {
		if !validPacakageName(pkg) {
			p = append(p, pkg)
			continue
		}
		n++

		name, ver, rel, arch, err := rpmQuery(out, pkg)
		if err != nil {
			if err == errorNotInstalled {
				// package not installed, add it to the list to be installed
				p = append(p, pkg)
				continue
			}
			return err
		}

		if matchPackage(pkg, name, ver, rel, arch) {
			n--
			continue
		}
		p = append(p, pkg)
	}

	if n == 0 {
		return nil
	}
	packages = p

	// install packages
	// TODO: handle downgrade
	_, err := runner.Cmd(out, cmdYum, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	// validate package version after installation
	for _, pkg := range packages {
		if !validPacakageName(pkg) {
			continue
		}

		if err = validateInstalledPackage(out, pkg); err != nil {
			return err
		}
	}

	// everything is installed
	return nil
}

func validPacakageName(pkg string) bool {
	// starts with '-', not a valid package name
	return pkg != "" && strings.TrimLeft(pkg, " ")[:1] != "-"
}

func matchPackage(pkg, name, ver, rel, arch string) bool {
	return name == pkg ||
		name+"-"+ver == pkg ||
		name+"-"+ver+"-"+rel == pkg ||
		name+"-"+ver+"-"+rel+"."+arch == pkg
}

func validateInstalledPackage(out io.Writer, pkg string) error {
	name, ver, rel, arch, err := rpmQuery(out, pkg)
	if err != nil {
		return err
	}
	if matchPackage(pkg, name, ver, rel, arch) {
		return nil
	}
	return errors.New(fmt.Sprintf("expected packgae version after installation: %q, got: %q", pkg, name+"-"+ver+"-"+rel+"."+arch))
}

func rpmQuery(out io.Writer, pkg string) (name, version, release, arch string, err error) {
	b, err := runner.Cmd(out, cmdRpm, []string{"-q", pkg}...).Output()
	if err != nil {
		if bytes.Contains(b, []byte("is not installed")) {
			err = errorNotInstalled
			return
		}
		return
	}

	return parseRpmPackageOutput(string(b))
}

func parseRpmPackageOutput(pkg string) (name, version, release, arch string, err error) {
	idx := strings.LastIndex(pkg, dotS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	arch = pkg[idx+1:]

	pkg = pkg[:idx]
	idx = strings.LastIndex(pkg, dashS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	release = pkg[idx+1:]

	pkg = pkg[:idx]
	idx = strings.LastIndex(pkg, dashS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	version = pkg[idx+1:]
	name = pkg[:idx]

	return
}

var _ KubernetesPackages = (*YumInstaller)(nil)
var _ ContainerDPackages = (*YumInstaller)(nil)

type YumInstaller struct{}

func NewYumInstaller() *YumInstaller {
	return &YumInstaller{}
}

func (y *YumInstaller) InstallKubernetesPackages(out io.Writer, kubernetesVersion string) error {
	p := []string{
		mapYumPackageVersion(kubelet, kubernetesVersion),
		mapYumPackageVersion(kubeadm, kubernetesVersion),
		mapYumPackageVersion(kubectl, kubernetesVersion),
		mapYumPackageVersion(kubernetescni, kubernetesVersion),
		disableExcludesKubernetes,
	}

	return YumInstall(out, p)
}

func (y *YumInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	pkg := []string{
		mapYumPackageVersion(kubeadm, kubernetesVersion),
		mapYumPackageVersion(kubelet, kubernetesVersion),       // kubeadm dependency
		mapYumPackageVersion(kubernetescni, kubernetesVersion), // kubeadm dependency
		disableExcludesKubernetes,
	}
	return YumInstall(out, pkg)
}

func (y *YumInstaller) InstallPrerequisites(out io.Writer, containerDVersion string) error {
	// yum install -y libseccomp
	if err := YumInstall(out, []string{"libseccomp"}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	return nil
}

func mapYumPackageVersion(pkg, kubernetesVersion string) string {
	switch pkg {
	case kubeadm:
		return "kubeadm-" + kubernetesVersion + "-0"

	case kubectl:
		return "kubectl-" + kubernetesVersion + "-0"

	case kubelet:
		return "kubelet-" + kubernetesVersion + "-0"

	case kubernetescni:
		ver, _ := semver.NewVersion(kubernetesVersion)
		c, _ := semver.NewConstraint(">=1.12.7,<1.13.x || >=1.13.5")
		if c.Check(ver) {
			return "kubernetes-cni-0.7.5-0"
		}
		return "kubernetes-cni-0.6.0-0"

	default:
		return ""
	}
}
