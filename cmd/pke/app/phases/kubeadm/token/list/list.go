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

package list

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
)

const (
	use   = "list"
	short = "List Kubernetes bootstrap token(s)"
)

var _ phases.Runnable = (*List)(nil)

type List struct {
	kubernetesVersion string
	imageRepository   string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &List{})
}

func (*List) Use() string {
	return use
}

func (*List) Short() string {
	return short
}

func (*List) RegisterFlags(flags *pflag.FlagSet) {}

func (*List) Validate(cmd *cobra.Command) error {
	return nil
}

func (*List) Run(out io.Writer) error {
	return nil
}
