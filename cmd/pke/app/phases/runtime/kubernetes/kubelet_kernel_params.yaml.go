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

// kubeletKernelParamsTemplate is a generated function returning the template as a string.
func kubeletKernelParamsTemplate() string {
	var tmpl = "vm.overcommit_memory=1\n" +
		"\n" +
		"# vm.oom_kill_allocating_task\n" +
		"# If this is set to zero, the OOM killer will scan through the entire\n" +
		"# tasklist and select a task based on heuristics to kill.\n" +
		"# If this is set to non-zero, the OOM killer simply kills the task that\n" +
		"# triggered the out-of-memory condition.\n" +
		"# The default value is 0.\n" +
		"vm.oom_kill_allocating_task=1\n" +
		"\n" +
		"kernel.panic=10\n" +
		"kernel.panic_on_oops=1\n" +
		""
	return tmpl
}
