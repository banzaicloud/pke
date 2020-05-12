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

// checkApiserverTemplate is a generated function returning the template as a string.
func checkApiserverTemplate() string {
	var tmpl = "#!/bin/sh\n" +
		"errorExit() {\n" +
		"    echo \"*** $*\" 1>&2\n" +
		"    exit 1\n" +
		"}\n" +
		"\n" +
		"export no_proxy=127.0.0.1 NO_PROXY=127.0.0.1\n" +
		"curl --silent --max-time 2 --insecure https://127.0.0.1:6443/healthz -o /dev/null || errorExit \"Error GET https://127.0.0.1:6443/healthz\"\n" +
		""
	return tmpl
}
