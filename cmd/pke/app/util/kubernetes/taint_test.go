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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTaints(t *testing.T) {
	testCases := []struct {
		taint    string
		expected *Taint
		err      bool
	}{
		{"foo=bar:NoSchedule", &Taint{Key: "foo", Value: "bar", Effect: "NoSchedule"}, false},
		{"foo:PreferNoSchedule", &Taint{Key: "foo", Value: "", Effect: "PreferNoSchedule"}, false},
		{"dedicated:NoSchedule-", &Taint{Key: "dedicated", Effect: "NoSchedule-"}, false},
		{"dedicated-", &Taint{Effect: "dedicated-"}, false},
		{"", nil, false},
		{"xxx", nil, true},
	}
	for _, tc := range testCases {
		taint, err := ParseTaints([]string{tc.taint})
		if tc.err {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		if tc.expected == nil {
			require.Empty(t, taint)
		} else {
			require.Equal(t, tc.expected, &taint[0])
		}
	}
}
