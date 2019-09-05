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

package file

import (
	"os"
	"path/filepath"
	"text/template"
)

// WriteTemplate write template output to file
func WriteTemplate(filename string, tmpl *template.Template, data interface{}) error {
	return WriteTemplateFlagPerm(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640, tmpl, data)
}

// WriteTemplateFlagPerm write template output to file with given flag and permission
func WriteTemplateFlagPerm(filename string, flag int, perm os.FileMode, tmpl *template.Template, data interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), perm|0110)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	return tmpl.Execute(w, data)
}
