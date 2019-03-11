package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
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
func NewCmdVersion(out io.Writer, gitVersion, gitCommit, gitTreeState, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print tool version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(out, cmd, gitVersion, gitCommit, gitTreeState, buildDate)
		},
	}
	cmd.Flags().StringP(constants.FlagOutput, constants.FlagOutputShort, "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunVersion(out io.Writer, cmd *cobra.Command, gitVersion, gitCommit, gitTreeState, buildDate string) error {
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
		_, _ = fmt.Fprintf(out, "kubeadm version: %#v\n", v)
	case "short":
		_, _ = fmt.Fprintf(out, "%s\n", v.GitVersion)
	case "yaml":
		y, err := yaml.Marshal(&v)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	case "json":
		y, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	default:
		return errors.Errorf("invalid output format: %s", of)
	}

	return nil
}
