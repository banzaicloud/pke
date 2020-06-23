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
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/stretchr/testify/require"
)

func TestWriteKubeadmConfig(t *testing.T) {
	t.SkipNow()
	filename := os.TempDir() + "kubeadm.conf"
	t.Log(filename)

	c := &ControlPlane{
		advertiseAddress:                "192.168.64.11:6443",
		apiServerHostPort:               "192.168.64.11:6443",
		kubeletCertificateAuthority:     "/etc/kubernetes/pki/ca.crt",
		clusterName:                     "my-cluster",
		kubernetesVersion:               "1.18.0",
		serviceCIDR:                     "10.32.0.0/24",
		podNetworkCIDR:                  "10.200.0.0/16",
		cloudProvider:                   constants.CloudProviderAmazon,
		nodepool:                        "pool1",
		controllerManagerSigningCA:      "/etc/kubernetes/pki/cm-signing-ca.crt",
		apiServerCertSANs:               []string{"almafa", "vadkorte"},
		withPluginPSP:                   true,
		withoutPluginDenyEscalatingExec: false,
		taints:                          []string{"node-role.kubernetes.io/master:NoSchedule"},
		withoutAuditLog:                 false,
	}

	err := c.WriteKubeadmConfig(os.Stdout, filename)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	t.Logf("%s\n", b)
}

func TestEnsureAPIServerConnection(t *testing.T) {
	t.SkipNow()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := ensureAPIServerConnection(os.Stdout, ctx, 5, "192.168.64.11")
	require.NoError(t, err)
}
