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

package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewPKECommand(in io.Reader, out io.Writer, gitVersion, gitCommit, gitTreeState, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "pke",
		Short:             "Bootstrap a secure Kubernetes cluster with Banzai Cloud Pipeline Kubernetes Engine (PKE)",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.ResetFlags()

	cmd.AddCommand(NewCmdInstall(out))
	cmd.AddCommand(NewCmdImage(out))
	cmd.AddCommand(NewCmdToken(out))
	cmd.AddCommand(NewCmdVersion(out, gitVersion, gitCommit, gitTreeState, buildDate))

	return cmd
}
