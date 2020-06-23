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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidVersion(t *testing.T) {
	testCases := []struct {
		version string
		valid   bool
	}{
		{"0.0.1", false},
		{"1.12.0", false},
		{"1.12.6", false},
		{"1.13.0", false},
		{"1.13.1", false},
		{"1.14.0", false},
		{"1.14.1", false},
		{"v1.14.1-beta.0", false},
		{"v1.15.0", true},
		{"v1.16.0", true},
		{"v1.17.0", true},
		{"v1.18.0", true},
		{"v1.19.0", false},
	}

	for _, tc := range testCases {
		err := validVersion(tc.version, constraint)
		if !tc.valid {
			assert.Error(t, err, tc.version)
		} else {
			assert.NoError(t, err, tc.version)
		}
	}
}
