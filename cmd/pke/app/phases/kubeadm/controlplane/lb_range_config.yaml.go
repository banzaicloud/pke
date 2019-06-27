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

// lbRangeConfigTemplate is a generated function returning the template as a string.
func lbRangeConfigTemplate() string {
	var tmpl = "apiVersion: v1\n" +
		"kind: ConfigMap\n" +
		"metadata:\n" +
		"  name: config\n" +
		"  namespace: metallb-system\n" +
		"data:\n" +
		"  config: |\n" +
		"    address-pools:\n" +
		"    - name: default\n" +
		"      protocol: layer2\n" +
		"      addresses:\n" +
		"      - {{ .Range }}\n" +
		""
	return tmpl
}
