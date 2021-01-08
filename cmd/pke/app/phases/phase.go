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

package phases

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Runnable interface for making phased commands.
type Runnable interface {
	Use() string
	Short() string
	RegisterFlags(flags *pflag.FlagSet)
	Validate(cmd *cobra.Command) error
	Run(out io.Writer) error
}

// NewCommand create new command.
func NewCommand(r Runnable) *cobra.Command {
	cmd := &cobra.Command{
		Use:   r.Use(),
		Short: r.Short(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := r.Validate(cmd); err != nil {
				return err
			}
			return r.Run(cmd.OutOrStdout())
		},
	}

	r.RegisterFlags(cmd.Flags())

	return cmd
}

// RunEAllSubcommands runs all sub-commands for a given phase.
func RunEAllSubcommands(cmd *cobra.Command, args []string) error {
	for _, c := range cmd.Commands() {
		if c.HasParent() {
			p := c.Parent()
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				if f := p.Flag(flag.Name); f != nil {
					*flag = *f
				}
			})
		}
		for p := c; p != nil; p = p.Parent() {
			if p.PersistentPreRunE != nil {
				if err := p.PersistentPreRunE(c, args); err != nil {
					return err
				}
				break
			} else if p.PersistentPreRun != nil {
				p.PersistentPreRun(c, args)
				break
			}
		}
		err := c.RunE(c, args)
		if err != nil {
			return err
		}
	}

	return nil
}

// MakeRunnable makes command phase runnable.
func MakeRunnable(cmd *cobra.Command) {
	visitedFlags := make(map[string]bool)
	for _, c := range cmd.Commands() {
		// local flags
		c.Flags().VisitAll(func(flag *pflag.Flag) {
			if visitedFlags[flag.Name] {
				return
			}
			cmd.Flags().AddFlag(flag)
			visitedFlags[flag.Name] = true
		})
		// persistent flags
		c.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if visitedFlags[flag.Name] {
				return
			}
			cmd.PersistentFlags().AddFlag(flag)
			visitedFlags[flag.Name] = true
		})
	}
}
