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

// keepalivedConfTemplate is a generated function returning the template as a string.
func keepalivedConfTemplate() string {
	var tmpl = "! Copyright (C) 2020 Banzai Cloud\n" +
		"! Configuration File for keepalived\n" +
		"global_defs {\n" +
		"  router_id LVS_DEVEL\n" +
		"}\n" +
		"vrrp_script check_apiserver {\n" +
		"  script \"/etc/keepalived/check_apiserver.sh\"\n" +
		"  interval 3\n" +
		"  weight -2\n" +
		"  fall 10\n" +
		"  rise 2\n" +
		"}\n" +
		"vrrp_instance VI_1 {\n" +
		"    state {{ .state }}\n" +
		"    interface {{ .iface }}\n" +
		"    virtual_router_id 51\n" +
		"    priority {{ .priority }}\n" +
		"    authentication {\n" +
		"        auth_type PASS\n" +
		"        auth_pass {{ .pass }}\n" +
		"    }\n" +
		"    virtual_ipaddress {\n" +
		"          {{ .vip }}\n" +
		"    }\n" +
		"    track_script {\n" +
		"        check_apiserver\n" +
		"    }\n" +
		"}\n" +
		""
	return tmpl
}
