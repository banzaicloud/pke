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

package file

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	contents       = "xxxx"
	contentsLonger = "yyyyyy"
)

func TestOverwrite(t *testing.T) {
	f, err := ioutil.TempFile("", "write_test")
	require.NoError(t, err)
	_ = f.Close()

	err = Overwrite(f.Name(), contents)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(f.Name())
	require.Equal(t, contents, string(b))

	err = Overwrite(f.Name(), contentsLonger)
	require.NoError(t, err)

	b, err = ioutil.ReadFile(f.Name())
	require.Equal(t, contentsLonger, string(b))

	err = Overwrite(f.Name(), contents)
	require.NoError(t, err)

	b, err = ioutil.ReadFile(f.Name())
	require.Equal(t, contents, string(b))
}
