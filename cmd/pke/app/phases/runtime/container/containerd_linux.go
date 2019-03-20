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

package container

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
	"github.com/pkg/errors"
)

const (
	containerDVersion     = "1.2.4"
	containerDSHA256      = "3391758c62d17a56807ddac98b05487d9e78e5beb614a0602caab747b0eda9e0"
	containerDURL         = "https://storage.googleapis.com/cri-containerd-release/cri-containerd-%s.linux-amd64.tar.gz"
	containerDVersionPath = "/opt/containerd/cluster/version"

	criConfFile = "/etc/sysctl.d/99-kubernetes-cri.conf"
	criConf     = `net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
`

	containerDConfFile = "/etc/systemd/system/kubelet.service.d/0-containerd.conf"
	containerDConf     = `[Service]
Environment="KUBELET_EXTRA_ARGS=--container-runtime=remote --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
`
)

func installRuntime(out io.Writer) error {
	if ver, err := linux.CentOSVersion(out); err == nil {
		if ver == "7" {
			return installCentOS7(out)
		}
		return constants.ErrUnsupportedOS
	}

	return constants.ErrUnsupportedOS
}

func installCentOS7(out io.Writer) error {
	// modprobe overlay
	if err := linux.ModprobeOverlay(out); err != nil {
		return errors.Wrap(err, "missing overlay Linux Kernel module")
	}

	// modprobe br_netfilter
	if err := linux.ModprobeBFNetFilter(out); err != nil {
		return errors.Wrap(err, "missing br_netfilter Linux Kernel module")
	}

	// Ensure network settings
	// cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
	// net.bridge.bridge-nf-call-iptables  = 1
	// net.bridge.bridge-nf-call-ip6tables = 1
	// net.ipv4.ip_forward                 = 1
	// EOF
	if err := file.Overwrite(criConfFile, criConf); err != nil {
		return err
	}

	if err := linux.SysctlLoadAllFiles(out); err != nil {
		return errors.Wrapf(err, "unable to load all sysctl rules from files")
	}

	// yum install -y libseccomp
	if err := linux.YumInstall(out, []string{"libseccomp"}); err != nil {
		return errors.Wrap(err, "unable to install libseccomp package")
	}

	_ = linux.SystemctlDisableAndStop(out, "containerd")

	// Check ContainerD installed or not
	if err := installContainerD(out); err != nil {
		return err
	}

	// # Start containerd.
	if err := linux.SystemctlEnableAndStart(out, "containerd"); err != nil {
		return err
	}

	_ = linux.SystemctlDisableAndStop(out, "kubelet")

	// cat > /etc/systemd/system/kubelet.service.d/0-containerd.conf <<EOF
	// [Service]
	// Environment="KUBELET_EXTRA_ARGS=--container-runtime=remote --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
	// EOF
	if err := os.MkdirAll(filepath.Dir(containerDConfFile), 0750); err != nil {
		return err
	}
	if err := file.Overwrite(containerDConfFile, containerDConf); err != nil {
		return err
	}

	// systemctl daemon-reload
	if err := linux.SystemctlReload(out); err != nil {
		return err
	}

	return nil
}

func installContainerD(out io.Writer) error {
	// Check ContainerD installed or not
	if _, err := os.Stat(containerDVersionPath); !os.IsNotExist(err) {
		// TODO: check ContainerD version
		_, _ = fmt.Fprintln(out, "ContainerD already installed, skipping download")
		return nil
	}
	// Download ContainerD tar.
	f, err := ioutil.TempFile("", "download_test")
	if err != nil {
		return errors.Wrapf(err, "unable to create temporary file: %q", f.Name())
	}
	// export CONTAINERD_VERSION="1.2.0"
	// export CONTAINERD_SHA256="ee076c6260de140f9aa6dee30b0e360abfb80af252d271e697982d1209ca5dee"
	// wget https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz
	dl := fmt.Sprintf(containerDURL, containerDVersion)
	u, err := url.Parse(dl)
	if err != nil {
		return errors.Wrapf(err, "failed to parse url: %q", dl)
	}
	_, _ = fmt.Fprintf(out, "wget %q -O %s\n", u.String(), f.Name())
	if err = file.Download(u, f.Name()); err != nil {
		return errors.Wrapf(err, "unable to download containerd. url: %q", u.String())
	}
	// echo "${CONTAINERD_SHA256} cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz" | sha256sum --check -
	_, _ = fmt.Fprintf(out, "echo \"%s %s\" | sha256sum --check -\n", containerDSHA256, f.Name())
	err = file.SHA256File(f.Name(), containerDSHA256)
	if err != nil {
		return errors.Wrapf(err, "hash mismatch. hash: %q", containerDSHA256)
	}

	// # Unpack.
	// tar --no-overwrite-dir -C / -xzf cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz
	fh, err := os.Open(f.Name())
	if err != nil {
		return err
	}
	err = file.Untar(out, fh)
	if err != nil {
		return err
	}

	return nil
}
