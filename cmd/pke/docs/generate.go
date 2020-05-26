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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/banzaicloud/pke/cmd/pke/app/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	const fmTemplate = `---
title: %s
generated_file: true
---
`
	const basePath = "/docs/pke/cli/reference/"

	// Customized Hugo output based on https://github.com/spf13/cobra/blob/master/doc/md_docs.md
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return basePath + strings.ToLower(base) + "/"
	}
	c := cmd.NewPKECommand("", "", "", "")
	c.SetOutput(ioutil.Discard)
	err := doc.GenMarkdownTreeCustom(c, ".", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
