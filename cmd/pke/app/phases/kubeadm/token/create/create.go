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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/token"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "create"
	short = "Create Kubernetes bootstrap token"

	cmdKubeadm = "kubeadm"
	kubeConfig = "/etc/kubernetes/admin.conf"
	caCertFile = "/etc/kubernetes/pki/ca.crt"
)

var _ phases.Runnable = (*Create)(nil)

type Create struct {
	o string
}

func NewCommand() *cobra.Command {
	return phases.NewCommand(&Create{})
}

func (*Create) Use() string {
	return use
}

func (*Create) Short() string {
	return short
}

func (*Create) RegisterFlags(flags *pflag.FlagSet) {
	flags.StringP(constants.FlagOutput, constants.FlagOutputShort, "", "Output format; available options are 'yaml', 'json' and 'short'")
}

func (c *Create) Validate(cmd *cobra.Command) error {
	var err error
	c.o, err = cmd.Flags().GetString(constants.FlagOutput)

	return err
}

func (c *Create) Run(out io.Writer) error {
	hash, err := token.CertHash(ioutil.Discard, caCertFile)
	if err != nil {
		return errors.Wrap(err, "failed to generate certificate hash")
	}

	cmd := runner.Cmd(ioutil.Discard, cmdKubeadm, "token", "create")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	o, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	var t *token.Token

	scn := bufio.NewScanner(bytes.NewReader(o))
	for scn.Scan() {
		line := scn.Text()
		if line == "" {
			continue
		}
		if err := scn.Err(); err != nil {
			return errors.Wrapf(err, "failed to scan output: %s", o)
		}

		idx := strings.IndexRune(line, '.')
		if idx < 0 {
			return errors.New("creation error: invalid token format")
		}

		t, err = token.Get(ioutil.Discard, "bootstrap-token-"+line[:idx], hash)
		if err != nil {
			return errors.Wrapf(err, "failed to get token for %q", line)
		}
	}

	switch c.o {
	default:
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintf(tw, "Token\tTTL\tExpires\tExpired\tCert Hash\n")
		_, _ = fmt.Fprintf(tw, "%s\t%dh\t%s\t%t\t%s\n", t.Token, t.TTL, t.Expires, t.Expired, t.CertHash)
		_ = tw.Flush()
		_ = tw.Flush()

	case "yaml":
		y, err := yaml.Marshal(&t)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	case "json":
		y, err := json.MarshalIndent(&t, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	}

	return nil
}
