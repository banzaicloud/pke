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

package node

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/stretchr/testify/require"
)

func TestWriteKubeadmConfig(t *testing.T) {
	t.SkipNow()
	filename := os.TempDir() + "kubeadm.conf"
	t.Log(filename)

	n := &Node{
		kubernetesVersion: "1.14.0",
		advertiseAddress:  "1.2.3.1:1234",
		apiServerHostPort: "1.2.3.4:1234",
		kubeadmToken:      "my.token",
		caCertHash:        "sha256:xxx",
		cloudProvider:     constants.CloudProviderAmazon,
		nodepool:          "pool2",
		taints:            []string{"node-role.kubernetes.io/master:NoSchedule-"},
	}

	err := n.writeKubeadmConfig(os.Stdout, filename)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	t.Logf("%s\n", b)
}
