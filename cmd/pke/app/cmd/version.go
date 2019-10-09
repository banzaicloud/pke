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
	"encoding/json"
	"fmt"
	"runtime"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type ClientVersion struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// NewCmdVersion provides the version information of banzai-cloud-pke.
func NewCmdVersion(gitVersion, gitCommit, gitTreeState, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print tool version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(cmd, gitVersion, gitCommit, gitTreeState, buildDate)
		},
	}
	cmd.Flags().StringP(constants.FlagOutput, constants.FlagOutputShort, "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunVersion(cmd *cobra.Command, gitVersion, gitCommit, gitTreeState, buildDate string) error {
	of, err := cmd.Flags().GetString(constants.FlagOutput)
	if err != nil {
		return err
	}

	v := struct {
		ClientVersion `json:"clientVersion"`
	}{ClientVersion{
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		BuildDate:    buildDate,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}}

	switch of {
	case "":
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "kubeadm version: %#v\n", v)
	case "short":
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", v.GitVersion)
	case "yaml":
		y, err := yaml.Marshal(&v)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(y))
	case "json":
		y, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(y))
	default:
		return errors.Errorf("invalid output format: %s", of)
	}

	return nil
}
