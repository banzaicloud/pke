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

package controlplane

import (
	"io"
	"os"
	"text/template"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
)

const (
	metalLbConfig   = "/etc/kubernetes/metallb-config.yaml"
	metalLbManifest = "https://raw.githubusercontent.com/google/metallb/v0.8.1/manifests/metallb.yaml"
)

func applyLbRange(out io.Writer, lbRange, cloudProvider string) error {
	if (cloudProvider != "" && cloudProvider != constants.CloudProviderVsphere) || lbRange == "" {
		return nil
	}

	cmd := runner.Cmd(out, cmdKubectl, "apply", "-f", metalLbManifest)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err := cmd.CombinedOutputAsync()
	if err != nil {
		return err
	}

	err = writeLbRangeConfig(out, metalLbConfig, lbRange)
	if err != nil {
		return err
	}

	cmd = runner.Cmd(out, cmdKubectl, "apply", "-f", metalLbConfig)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	_, err = cmd.CombinedOutputAsync()
	return err
}

//go:generate templify -t ${GOTMPL} -p controlplane -f lbRangeConfig lb_range_config.yaml.tmpl

func writeLbRangeConfig(out io.Writer, filename, lbRange string) error {
	tmpl, err := template.New("metallb-config").Parse(lbRangeConfigTemplate())
	if err != nil {
		return err
	}

	type data struct {
		Range string
	}

	d := data{
		Range: lbRange,
	}

	return file.WriteTemplate(filename, tmpl, d)
}
