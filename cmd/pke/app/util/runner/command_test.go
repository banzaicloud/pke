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

package runner

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	err := c.Run()
	require.NoError(t, err)
}

func TestOutput(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	out, err := c.Output()
	require.NoError(t, err)
	require.Equal(t, []byte("ok\n"), out)
}

func TestPipeOut(t *testing.T) {
	c := Cmd(ioutil.Discard, "echo", "ok")
	o, err := c.StdoutPipe()
	require.NoError(t, err)

	err = c.Start()
	require.NoError(t, err)

	out, err := ioutil.ReadAll(o)
	require.NoError(t, err)
	require.Equal(t, []byte("ok\n"), out)

	err = c.Wait()
	require.NoError(t, err)
}
