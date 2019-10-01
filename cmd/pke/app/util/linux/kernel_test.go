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

package linux

import (
	"bytes"
	"testing"
)

func TestKernelVersionConstraint(t *testing.T) {
	type args struct {
		constraint string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{">=0.0.1", args{">=0.0.1-0"}, false},
		{"<0.0.1", args{"<0.0.1-0"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := KernelVersionConstraint(out, tt.args.constraint)
			if (err != nil) != tt.wantErr {
				t.Errorf("KernelVersionConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
