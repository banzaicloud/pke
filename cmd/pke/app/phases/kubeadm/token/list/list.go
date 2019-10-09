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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	use   = "list"
	short = "List Kubernetes bootstrap token(s)"

	cmdKubectl = "kubectl"
	kubeConfig = "/etc/kubernetes/admin.conf"
	caCertFile = "/etc/kubernetes/pki/ca.crt"
)

var _ phases.Runnable = (*List)(nil)

type List struct {
	o string
}

func NewCommand() *cobra.Command {
	return phases.NewCommand(&List{})
}

func (*List) Use() string {
	return use
}

func (*List) Short() string {
	return short
}

func (*List) RegisterFlags(flags *pflag.FlagSet) {
	flags.StringP(constants.FlagOutput, constants.FlagOutputShort, "", "Output format; available options are 'yaml', 'json' and 'short'")
}

func (l *List) Validate(cmd *cobra.Command) error {
	var err error
	l.o, err = cmd.Flags().GetString(constants.FlagOutput)

	return err
}

func (l *List) Run(out io.Writer) error {
	hash, err := token.CertHash(ioutil.Discard, caCertFile)
	if err != nil {
		return errors.Wrap(err, "failed to generate certificate hash")
	}

	// kubectl get secret -n kube-system -o jsonpath='{range .items[?(.type=="bootstrap.kubernetes.io/token")]}{.metadata.name}{"\n"}{end}'
	args := []string{"get", "secret", "-n", "kube-system", "-o", `jsonpath={range .items[?(.type=="bootstrap.kubernetes.io/token")]}{.metadata.name}{"\n"}{end}`}
	cmd := runner.Cmd(ioutil.Discard, cmdKubectl, args...)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeConfig)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "failed to read secret")
	}

	var list = token.Output{}

	scn := bufio.NewScanner(bytes.NewReader(o))
	for scn.Scan() {
		line := scn.Text()
		if line == "" {
			continue
		}
		if err := scn.Err(); err != nil {
			return errors.Wrapf(err, "failed to scan output: %s", o)
		}

		t, err := token.Get(ioutil.Discard, line, hash)
		if err != nil {
			return errors.Wrapf(err, "failed to get token for %q", line)
		}

		list.Tokens = append(list.Tokens, t)
	}

	switch l.o {
	default:
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintf(tw, "Token\tTTL\tExpires\tExpired\tCert Hash\n")
		for _, row := range list.Tokens {
			_, _ = fmt.Fprintf(tw, "%s\t%dh\t%s\t%t\t%s\n", row.Token, row.TTL, row.Expires, row.Expired, row.CertHash)
		}
		_ = tw.Flush()

	case "yaml":
		y, err := yaml.Marshal(&list)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	case "json":
		y, err := json.MarshalIndent(&list, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(out, string(y))
	}

	return nil
}
