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
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/token/create"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/token/list"
	"github.com/spf13/cobra"
)

// NewCmdToken provides the version information of banzai-cloud-pke.
func NewCmdToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage Kubernetes bootstrap tokens",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
