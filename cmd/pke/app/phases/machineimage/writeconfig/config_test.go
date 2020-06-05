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

package writeconfig

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteConfig(t *testing.T) {
	fileName := os.TempDir() + "/pke.yaml"
	t.Cleanup(func() {
		_ = os.Remove(fileName)
	})

	t.Log(fileName)

	c := &WriteConfig{
		kubernetesVersion: "1.18.0",
		containerRuntime:  "containerd",
	}

	err := c.WriteConfig(os.Stdout, fileName)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(fileName)
	require.NoError(t, err)

	const expected = `kubernetes:
  version: 1.18.0
  installed: true
containerRuntime:
  type: containerd
  installed: true
`

	assert.Equal(t, expected, string(b))
}
