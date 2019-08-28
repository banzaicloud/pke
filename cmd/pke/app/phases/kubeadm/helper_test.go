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

package kubeadm

import (
	"testing"
)

func TestKubeReservedMemory(t *testing.T) {
	type args struct {
		actualMemory uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"256Gi", args{256 * 1024 << 20}, "15729Mi"},
		{"128Gi", args{128 * 1024 << 20}, "9961Mi"},
		{"64Gi", args{64 * 1024 << 20}, "6029Mi"},
		{"32Gi", args{32 * 1024 << 20}, "4030Mi"},
		{"16Gi", args{16 * 1024 << 20}, "2720Mi"},
		{"8Gi", args{8 * 1024 << 20}, "1843Mi"},
		{"4Gi", args{4 * 1024 << 20}, "995Mi"},
		{"1Gi", args{1 * 1024 << 20}, "307Mi"},
		{"500Mi", args{500 * 1024 << 10}, "200Mi"},
		{"100Mi", args{100 * 1024 << 10}, "0Mi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KubeReservedMemory(tt.args.actualMemory); got != tt.want {
				t.Errorf("kubeReservedMemory(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
