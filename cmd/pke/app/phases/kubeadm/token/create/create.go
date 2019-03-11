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

package create

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "create"
	short = "Create Kubernetes bootstrap token"
)

var _ phases.Runnable = (*Create)(nil)

type Create struct {
	kubernetesVersion string
	imageRepository   string
}

func NewCommand(out io.Writer) *cobra.Command {
	return phases.NewCommand(out, &Create{})
}

func (*Create) Use() string {
	return use
}

func (*Create) Short() string {
	return short
}

func (*Create) RegisterFlags(flags *pflag.FlagSet) {}

func (*Create) Validate(cmd *cobra.Command) error {
	return nil
}

func (*Create) Run(out io.Writer) error {
	return nil
}
