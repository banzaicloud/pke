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

package certificates

import (
	"context"
	"fmt"
	"io"
	"os"

	"emperror.dev/errors"
	"github.com/antihax/optional"
	"github.com/banzaicloud/pke/.gen/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/config"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm"
	"github.com/banzaicloud/pke/cmd/pke/app/util/file"
	pipelineutil "github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "pipeline-certificates"
	short = "Pipeline pre-generated certificate download"

	etcdDir                 = "/etc/kubernetes/pki/etcd"
	etcdCACert              = "/etc/kubernetes/pki/etcd/ca.crt"
	etcdCAKey               = "/etc/kubernetes/pki/etcd/ca.key"
	kubernetesCASigningCert = "/etc/kubernetes/pki/cm-signing-ca.crt"
	kubernetesCACert        = "/etc/kubernetes/pki/ca.crt"
	kubernetesCAKey         = "/etc/kubernetes/pki/ca.key"
	frontProxyCACert        = "/etc/kubernetes/pki/front-proxy-ca.crt"
	frontProxyCAKey         = "/etc/kubernetes/pki/front-proxy-ca.key"
	saPub                   = "/etc/kubernetes/pki/sa.pub"
	saKey                   = "/etc/kubernetes/pki/sa.key"
)

var _ phases.Runnable = (*Certificates)(nil)

type Certificates struct {
	config config.Config

	pipelineEnabled        bool
	pipelineAPIEndpoint    string
	pipelineAPIToken       string
	pipelineAPIInsecure    bool
	pipelineOrganizationID int32
	pipelineClusterID      int32
	kubernetesVersion      string
}

func NewCommand(config config.Config) *cobra.Command {
	return phases.NewCommand(&Certificates{config: config})
}

func (c *Certificates) Use() string {
	return use
}

func (c *Certificates) Short() string {
	return short
}

func (c *Certificates) RegisterFlags(flags *pflag.FlagSet) {
	// Pipeline
	flags.StringP(constants.FlagPipelineAPIEndpoint, constants.FlagPipelineAPIEndpointShort, "", "Pipeline API server url")
	flags.StringP(constants.FlagPipelineAPIToken, constants.FlagPipelineAPITokenShort, "", "Token for accessing Pipeline API")
	flags.Bool(constants.FlagPipelineAPIInsecure, false, "If the Pipeline API should not verify the API's certificate")
	flags.Int32(constants.FlagPipelineOrganizationID, 0, "Organization ID to use with Pipeline API")
	flags.Int32(constants.FlagPipelineClusterID, 0, "Cluster ID to use with Pipeline API")
	// Kubernetes version
	flags.String(constants.FlagKubernetesVersion, c.config.Kubernetes.Version, "Kubernetes version")
}

func (c *Certificates) Validate(cmd *cobra.Command) error {
	if !pipelineutil.Enabled(cmd) {
		// TODO: Warning
		return nil
	}
	c.pipelineEnabled = true

	var err error
	c.pipelineAPIEndpoint, c.pipelineAPIToken, c.pipelineAPIInsecure, c.pipelineOrganizationID, c.pipelineClusterID, err = pipelineutil.CommandArgs(cmd)
	if err != nil {
		return err
	}
	err = pipelineutil.ValidArgs(c.pipelineAPIEndpoint, c.pipelineAPIToken, c.pipelineAPIInsecure, c.pipelineOrganizationID, c.pipelineClusterID)
	if err != nil {
		return err
	}

	c.kubernetesVersion, err = cmd.Flags().GetString(constants.FlagKubernetesVersion)
	return err
}

func (c *Certificates) Run(out io.Writer) error {
	if !c.pipelineEnabled {
		return nil
	}

	_, _ = fmt.Fprintf(out, "[%s] running\n", c.Use())

	if err := pipelineutil.ValidArgs(c.pipelineAPIEndpoint, c.pipelineAPIToken, c.pipelineAPIInsecure, c.pipelineOrganizationID, c.pipelineClusterID); err != nil {
		_, _ = fmt.Fprintf(out, "[WARNING] Skipping %s phase due to missing Pipeline API endpoint. err: %v\n", use, err)
		return nil
	}

	var err error
	req := &pipeline.GetSecretsOpts{
		Type_:  optional.NewString("pkecert"),
		Values: optional.NewBool(true),
		Tags:   optional.NewInterface([]string{fmt.Sprintf("clusterID:%d", c.pipelineClusterID)}),
	}
	pc := pipelineutil.Client(out, c.pipelineAPIEndpoint, c.pipelineAPIToken, c.pipelineAPIInsecure)
	secrets, _, err := pc.SecretsApi.GetSecrets(context.Background(), c.pipelineOrganizationID, req)
	if err != nil {
		return err
	}
	if n := len(secrets); n <= 0 || n > 1 {
		ids := func(secrets []pipeline.SecretItem) []string {
			ids := make([]string, len(secrets))
			for k, s := range secrets {
				ids[k] = s.Id
			}
			return ids
		}(secrets)

		return errors.Errorf("multiple or none PKE certificates are returned for cluster: %q", ids)
	}

	secret := secrets[0]

	_, _ = fmt.Fprintf(out, "[%s] creating directory: %q\n", use, etcdDir)
	err = os.MkdirAll(etcdDir, 0750)
	if err != nil {
		return err
	}
	// /etc/kubernetes/pki/etcd/ca.crt
	if err = write(out, etcdCACert, secret.Values["etcdCaCert"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/etcd/ca.key
	if err = write(out, etcdCAKey, secret.Values["etcdCaKey"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/cm-signing-ca.crt
	if err = write(out, kubernetesCASigningCert, secret.Values["kubernetesCaSigningCert"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/ca.crt
	if err = write(out, kubernetesCACert, secret.Values["kubernetesCaCert"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/ca.key
	if err = write(out, kubernetesCAKey, secret.Values["kubernetesCaKey"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/front-proxy-ca.crt
	if err = write(out, frontProxyCACert, secret.Values["frontProxyCaCert"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/front-proxy-ca.key
	if err = write(out, frontProxyCAKey, secret.Values["frontProxyCaKey"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/sa.pub
	if err = write(out, saPub, secret.Values["saPub"]); err != nil {
		return err
	}

	// /etc/kubernetes/pki/sa.key
	if err = write(out, saKey, secret.Values["saKey"]); err != nil {
		return err
	}

	if key, ok := secret.Values["enc"].(string); ok {
		return kubeadm.WriteEncryptionProviderConfig(out, kubeadm.EncryptionProviderConfig, c.kubernetesVersion, key)
	}

	return nil
}

func write(out io.Writer, filename string, value interface{}) error {
	_, _ = fmt.Fprintf(out, "[%s] writing file: %s\n", use, filename)
	if v, ok := value.(string); ok {
		return file.Overwrite(filename, v)
	}

	return errors.Errorf("unexpected interface type. expected string, got: %T", value)
}
