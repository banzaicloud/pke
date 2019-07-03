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

// cniBridgeTemplate is a generated function returning the template as a string.
func cniBridgeTemplate() string {
	var tmpl = "{\n" +
		"    \"cniVersion\": \"0.3.1\",\n" +
		"    \"name\": \"bridge\",\n" +
		"    \"type\": \"bridge\",\n" +
		"    \"bridge\": \"cnio0\",\n" +
		"    \"isGateway\": true,\n" +
		"    \"ipMasq\": true,\n" +
		"    \"ipam\": {\n" +
		"        \"type\": \"host-local\",\n" +
		"        \"ranges\": [\n" +
		"          [{\"subnet\": \"{{ .PodNetworkCIDR }}\"}]\n" +
		"        ],\n" +
		"        \"routes\": [{\"dst\": \"0.0.0.0/0\"}]\n" +
		"    }\n" +
		"}\n" +
		""
	return tmpl
}
