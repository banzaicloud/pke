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
	"strings"

	"emperror.dev/errors"
)

type Taint struct {
	Key    string
	Value  string
	Effect string
}

func isTaint(s string) bool {
	return strings.Contains(s, "=") || strings.Contains(s, ":") || strings.HasSuffix(s, "-")
}

func ParseTaints(taints []string) ([]Taint, error) {
	t := make([]Taint, 0)
	for _, taint := range taints {
		// Remove leading and trailing spaces
		taint = strings.Trim(taint, " ")
		// Skip empty
		if taint == "" {
			continue
		}
		// Validate taint
		if !isTaint(taint) {
			return nil, errors.New("invalid taint spec: " + taint)
		}

		var (
			key    string
			value  string
			effect string
		)

		// key=value
		eq := strings.Index(taint, "=")
		if eq > -1 {
			key = taint[:eq]
		}

		// effect
		colon := strings.Index(taint, ":")
		if colon > -1 {
			effect = taint[colon+1:]
			if eq > -1 {
				value = taint[eq+1 : colon]
			} else {
				key = taint[eq+1 : colon]
			}
		} else {
			effect = taint
		}

		t = append(t, Taint{
			Key:    key,
			Value:  value,
			Effect: effect,
		})
	}

	return t, nil
}
