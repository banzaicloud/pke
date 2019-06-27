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

// kubeadmVsphereConfigTemplate is a generated function returning the template as a string.
func kubeadmVsphereConfigTemplate() string {
	var tmpl = "[Global]\n" +
		"\n" +
		"[VirtualCenter \"{{ .Server }}\"]\n" +
		"port = \"{{ .Port }}\"\n" +
		"datacenters = \"{{ .Datacenter }}\"\n" +
		"{{ if .Fingerprint }}\n" +
		"    thumbprint = \"{{ .Fingerprint }}\"\n" +
		"{{ end }}\n" +
		"{{ if .Username }}\n" +
		"    user = \"{{ .Username }}\"\n" +
		"{{ end }}\n" +
		"{{ if .Password }}\n" +
		"    password = \"{{ .Password }}\"\n" +
		"{{ end }}\n" +
		"\n" +
		"[Workspace]\n" +
		"server = \"{{ .Server }}\"\n" +
		"datacenter = \"{{ .Datacenter }}\"\n" +
		"default-datastore = \"{{ .Datastore }}\"\n" +
		"resourcepool-path = \"{{ .ResourcePool }}\"\n" +
		"folder = \"{{ .Folder }}\"\n" +
		"\n" +
		"[Disk]\n" +
		"scsicontrollertype = pvscsi\n" +
		""
	return tmpl
}
