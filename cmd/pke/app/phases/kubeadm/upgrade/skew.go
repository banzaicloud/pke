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

package upgrade

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"gopkg.in/yaml.v2"
)

const (
	cmdKubectl = "kubectl"

	MaximumAllowedMinorVersionUpgradeSkew = 1
)

func RunWithSkewCheck(out io.Writer, use, kubernetesVersion, kubeConfig string, minor, patch func(out io.Writer, from, to *semver.Version) error) error {
	// query current kubernetes version
	clientVersion, serverVersion, err := kubectlVersion(out, kubeConfig)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "[%s] client version: %s\n", use, clientVersion)
	_, _ = fmt.Fprintf(out, "[%s] server version: %s\n", use, serverVersion)

	cmVer, err := configmapVersion(out, kubeConfig)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(out, "[%s] configmap version: %s\n", use, cmVer)

	// minor or patch update
	srvVer, err := semver.NewVersion(serverVersion)
	if err != nil {
		return err
	}
	ver, err := semver.NewVersion(kubernetesVersion)
	if err != nil {
		return err
	}
	if srvVer.Major() != ver.Major() {
		return errors.New(fmt.Sprintf(
			"major version upgrade not supported. trying to upgrade from %d.x to %d.x",
			srvVer.Major(),
			ver.Major(),
		))
	}

	if srvVer.Minor() != ver.Minor() {
		if srvVer.Minor() > ver.Minor() {
			return errors.New(fmt.Sprintf(
				"downgrade not supported. trying to upgrade from %s to %s",
				srvVer,
				ver,
			))
		}
		if srvVer.Minor()+MaximumAllowedMinorVersionUpgradeSkew < ver.Minor() {
			return errors.New(fmt.Sprintf(
				"only %d minor version can be updated at a time. trying to upgrade from %d.%d to %d.%d",
				MaximumAllowedMinorVersionUpgradeSkew,
				srvVer.Major(),
				srvVer.Minor(),
				ver.Major(),
				ver.Minor(),
			))
		}
		// Minor version bump
		return minor(out, srvVer, ver)
	}

	if srvVer.Patch() > ver.Patch() {
		return errors.New(fmt.Sprintf(
			"downgrade not supported. trying to upgrade from %s to %s",
			srvVer,
			ver,
		))
	}

	return patch(out, srvVer, ver)
}

func kubectlVersion(out io.Writer, kubeConfig string) (clientVersion, serverVersion string, err error) {
	cmd := runner.Cmd(ioutil.Discard, cmdKubectl, "version", "-o", "json")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	ver, err := cmd.Output()
	if err != nil {
		return
	}
	serverVer := struct {
		ClientVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"clientVersion"`
		ServerVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"serverVersion"`
	}{}
	err = json.Unmarshal(ver, &serverVer)
	if err != nil {
		return
	}

	clientVersion = serverVer.ClientVersion.GitVersion
	serverVersion = serverVer.ServerVersion.GitVersion
	return
}

func configmapVersion(out io.Writer, kubeConfig string) (version string, err error) {
	cmd := runner.Cmd(ioutil.Discard, cmdKubectl, "-n", "kube-system", "get", "cm", "kubeadm-config", "-ojsonpath={.data.ClusterConfiguration}")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	o, err := cmd.Output()
	if err != nil {
		return
	}

	cmVer := struct {
		KubernetesVersion string `yaml:"kubernetesVersion"`
	}{}
	err = yaml.Unmarshal(o, &cmVer)
	if err != nil {
		return
	}

	return cmVer.KubernetesVersion, nil
}
