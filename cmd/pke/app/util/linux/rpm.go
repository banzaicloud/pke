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
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdRpm                    = "/bin/rpm"
	kubeadm                   = "kubeadm"
	kubectl                   = "kubectl"
	kubelet                   = "kubelet"
	kubernetescni             = "kubernetes-cni"
	kubernetesCNIVersion      = "0.8.7"
	disableExcludesKubernetes = "--disableexcludes=kubernetes"
	selinuxConfig             = "/etc/selinux/config"
)

type pkg struct {
	name    string
	version string
}

func (p pkg) String() string {
	if p.version == "" {
		return p.name
	}
	return p.name + "-" + p.version
}

type packages []pkg

func (p packages) strings() (out []string) {
	for _, pkg := range p {
		out = append(out, pkg.String())
	}
	return out
}

func checkRPMPackages(out io.Writer, packages packages) error {
	for _, pkg := range packages {
		output, err := rpmQuery(out, pkg.name)
		if err != nil {
			return errors.WrapIff(err, "failed to query installed package %q by name", pkg.name)
		}

		_, err = rpmQuery(out, pkg.String())
		if err != nil {
			return errors.WrapIff(err, "failed to query installed package %q by name and version (have %q instead)", pkg.String(), output)
		}
	}
	return nil
}

func rpmQuery(out io.Writer, pkg string) (string, error) {
	b, err := runner.Cmd(out, cmdRpm, []string{"-q", pkg}...).Output()
	return strings.TrimSpace(string(b)), err
}
