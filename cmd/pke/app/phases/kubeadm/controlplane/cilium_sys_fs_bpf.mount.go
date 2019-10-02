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

// ciliumSysFsBpfTemplate is a generated function returning the template as a string.
func ciliumSysFsBpfTemplate() string {
	var tmpl = "[Unit]\n" +
		"Description=Cilium BPF mounts\n" +
		"Documentation=http://docs.cilium.io/\n" +
		"DefaultDependencies=no\n" +
		"Before=local-fs.target umount.target\n" +
		"After=swap.target\n" +
		"\n" +
		"[Mount]\n" +
		"What=bpffs\n" +
		"Where=/sys/fs/bpf\n" +
		"Type=bpf\n" +
		"\n" +
		"[Install]\n" +
		"WantedBy=multi-user.target\n" +
		""
	return tmpl
}
