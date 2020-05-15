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
	"text/template"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	"github.com/banzaicloud/pke/cmd/pke/app/util/linux"
)

const (
	containerdVersion     = "1.3.3"
	containerdSHA256      = "24ce7ad6b489fb25d07d2a3bb50e443fcce1ac3318f8cc0831e00668c2c9fd86"
	containerdURL         = "https://storage.googleapis.com/cri-containerd-release/cri-containerd-%s.linux-amd64.tar.gz"
	containerdVersionPath = "/opt/containerd/cluster/version"
	containerdConf        = "/etc/containerd/config.toml"

	criConfFile = "/etc/sysctl.d/99-kubernetes-cri.conf"
	criConf     = `net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
`
)

func (r *Runtime) installContainerd(out io.Writer) error {
	pm, err := linux.ContainerdPackagesImpl(out)
	if err != nil {
		return err
	}

	// modprobe overlay
	if err := linux.ModprobeOverlay(out); err != nil {
		return errors.Wrap(err, "missing overlay Linux Kernel module")
	}

	// modprobe br_netfilter
	if err := linux.ModprobeBRNetFilter(out); err != nil {
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

	if err := pm.InstallContainerdPrerequisites(out, containerdVersion); err != nil {
		return errors.Wrap(err, "unable to install containerd prerequisites")
	}

	_ = linux.SystemctlDisableAndStop(out, "containerd")

	// Check containerd installed or not
	if err := installContainerd(out, r.imageRepository); err != nil {
		return err
	}

	// # Start containerd.
	if err := linux.SystemctlEnableAndStart(out, "containerd"); err != nil {
		return err
	}

	_ = linux.SystemctlDisableAndStop(out, "kubelet")

	// systemctl daemon-reload
	return linux.SystemctlReload(out)
}

func installContainerd(out io.Writer, imageRepository string) error {
	// Check containerd installed or not
	if _, err := os.Stat(containerdVersionPath); !os.IsNotExist(err) {
		// TODO: check containerd version
		_, _ = fmt.Fprintln(out, "containerd already installed, skipping download")
		return nil
	}
	// Download containerd tar.
	f, err := ioutil.TempFile("", "containerd")
	if err != nil {
		return errors.Wrapf(err, "unable to create temporary file: %q", f.Name())
	}
	defer func() { _ = f.Close() }()
	// export CONTAINERD_VERSION="1.3.3"
	// export CONTAINERD_SHA256="24ce7ad6b489fb25d07d2a3bb50e443fcce1ac3318f8cc0831e00668c2c9fd86"
	// wget https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz
	dl := fmt.Sprintf(containerdURL, containerdVersion)
	u, err := url.Parse(dl)
	if err != nil {
		return errors.Wrapf(err, "failed to parse url: %q", dl)
	}
	_, _ = fmt.Fprintf(out, "wget %q -O %s\n", u.String(), f.Name())
	if err = file.Download(u, f.Name()); err != nil {
		return errors.Wrapf(err, "unable to download containerd. url: %q", u.String())
	}
	// echo "${CONTAINERD_SHA256} cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz" | sha256sum --check -
	_, _ = fmt.Fprintf(out, "echo \"%s %s\" | sha256sum --check -\n", containerdSHA256, f.Name())
	err = file.SHA256File(f.Name(), containerdSHA256)
	if err != nil {
		return errors.Wrapf(err, "hash mismatch. hash: %q", containerdSHA256)
	}

	// Unpack.
	// tar --no-overwrite-dir -C / -xzf cri-containerd-${CONTAINERD_VERSION}.linux-amd64.tar.gz
	fh, err := os.Open(f.Name())
	if err != nil {
		return err
	}
	defer func() { _ = fh.Close() }()

	err = file.Untar(out, fh)
	if err != nil {
		return err
	}

	return writeContainerdConfig(out, containerdConf, imageRepository)
}

//go:generate templify -t ${GOTMPL} -p container -f containerdConfig containerd_config.toml.tmpl

func writeContainerdConfig(out io.Writer, filename, imageRepository string) error {
	tmpl, err := template.New("containerd-config").Parse(containerdConfigTemplate())
	if err != nil {
		return err
	}

	type data struct {
		ImageRepository string
	}

	d := data{
		ImageRepository: imageRepository,
	}

	return file.WriteTemplate(filename, tmpl, d)
}
