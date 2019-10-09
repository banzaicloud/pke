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

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	cmdModprobe = "/sbin/modprobe"
)

func Modprobe(out io.Writer, module string) error {
	return runner.Cmd(out, cmdModprobe, module).Run()
}

func ModprobeOverlay(out io.Writer) error {
	return Modprobe(out, "overlay")
}

func ModprobeBRNetFilter(out io.Writer) error {
	return Modprobe(out, "br_netfilter")
}

func kubeProxyIPVSModules(out io.Writer) ([]string, error) {
	modules := []string{
		"ip_vs",
		"ip_vs_rr",
		"ip_vs_wrr",
		"ip_vs_sh",
	}

	// https://github.com/kubernetes/kubernetes/blob/d5d7db476d044c7489c120e20f07b1283a81310d/pkg/proxy/ipvs/README.md#prerequisite
	conntrack := "nf_conntrack_ipv4"
	err := KernelVersionConstraint(out, "<4.19")
	if errors.Is(err, constants.ErrUnsupportedKernelVersion) {
		err = nil
		conntrack = "nf_conntrack"
	}
	if err != nil {
		return nil, err
	}
	modules = append(modules, conntrack)

	return modules, nil
}

func ModprobeKubeProxyIPVSModules(out io.Writer) error {
	modules, err := kubeProxyIPVSModules(out)
	if err != nil {
		return err
	}
	for _, module := range modules {
		err = Modprobe(out, module)
		if err != nil {
			return errors.Wrapf(err, "missing %s Linux Kernel module", module)
		}
	}

	return nil
}
