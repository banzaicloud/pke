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

package ready

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"emperror.dev/errors"
	"github.com/banzaicloud/pke/.gen/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/network"
	pipelineutil "github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	use   = "pipeline-ready"
	short = "Register node as ready at Pipeline"

	kubeConfig = "/etc/kubernetes/admin.conf"

	RoleMaster Role = "master"
	RoleWorker Role = "worker"
)

type Role string

var _ phases.Runnable = (*Ready)(nil)

type Ready struct {
	role                   Role // accepted values: master, worker
	cidr                   string
	pipelineEnabled        bool
	pipelineAPIEndpoint    string
	pipelineAPIToken       string
	pipelineAPIInsecure    bool
	pipelineOrganizationID int32
	pipelineClusterID      int32
	pipelineNodepool       string
}

func NewCommand(role Role) *cobra.Command {
	return phases.NewCommand(&Ready{role: role})
}

func (r *Ready) Use() string {
	return use
}

func (r *Ready) Short() string {
	return short
}

func (r *Ready) RegisterFlags(flags *pflag.FlagSet) {
	flags.StringP(constants.FlagPipelineAPIEndpoint, constants.FlagPipelineAPIEndpointShort, "", "Pipeline API server url")
	flags.StringP(constants.FlagPipelineAPIToken, constants.FlagPipelineAPITokenShort, "", "Token for accessing Pipeline API")
	flags.Bool(constants.FlagPipelineAPIInsecure, false, "If the Pipeline API should not verify the API's certificate")
	flags.Int32(constants.FlagPipelineOrganizationID, 0, "Organization ID to use with Pipeline API")
	flags.Int32(constants.FlagPipelineClusterID, 0, "Cluster ID to use with Pipeline API")
	flags.String(constants.FlagPipelineNodepool, "", "name of the nodepool the node belongs to")
	flags.String(constants.FlagInfrastructureCIDR, "192.168.64.0/20", "network CIDR for the actual machine")
}

func (r *Ready) Validate(cmd *cobra.Command) error {
	// Run is optional, only validate if Pipeline credentials are present.
	if !pipelineutil.Enabled(cmd) {
		return nil
	}
	r.pipelineEnabled = true

	var err error
	if r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineAPIInsecure, r.pipelineOrganizationID, r.pipelineClusterID, err = pipelineutil.CommandArgs(cmd); err != nil {
		return err
	}

	if err = pipelineutil.ValidArgs(r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineAPIInsecure, r.pipelineOrganizationID, r.pipelineClusterID); err != nil {
		return err
	}

	if r.pipelineNodepool, err = cmd.Flags().GetString(constants.FlagPipelineNodepool); err != nil {
		return err
	}

	if err := validator.NotEmpty(map[string]interface{}{
		"Pipeline nodepool": r.pipelineNodepool,
	}); err != nil {
		return err
	}

	r.cidr, err = cmd.Flags().GetString(constants.FlagInfrastructureCIDR)

	return err
}

// Run Register installed machine at Pipeline.
// Optional step.
// Skipped if no Pipeline credentials are provided.
func (r *Ready) Run(out io.Writer) error {
	if !r.pipelineEnabled {
		return nil
	}

	_, _ = fmt.Fprintf(out, "[%s] running\n", r.Use())

	if err := pipelineutil.ValidArgs(r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineAPIInsecure, r.pipelineOrganizationID, r.pipelineClusterID); err != nil {
		_, _ = fmt.Fprintf(out, "[WARNING][%s] Skipping phase due to missing Pipeline API endpoint credentials. %s\n", use, err)
		return nil
	}

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// ip
	ips, err := network.IPv4Addresses()
	if err != nil {
		return err
	}
	ip, err := network.ContainsFirst(r.cidr, ips)
	if err != nil {
		return err
	}

	// post node ready
	c := pipelineutil.Client(out, r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineAPIInsecure)
	req := pipeline.PostReadyPkeNodeRequest{
		Name:     hostname,
		NodePool: r.pipelineNodepool,
		Ip:       ip.String(),
	}
	_, _ = fmt.Fprintf(out, "[%s] %v\n", use, req)

	if r.role == RoleMaster {
		b, err := ioutil.ReadFile(kubeConfig)
		if err != nil {
			return errors.Wrapf(err, "unable to read file: %s", kubeConfig)
		}
		req.Config = base64.StdEncoding.EncodeToString(b)
	}

	_, _, err = c.ClustersApi.PostReadyPKENode(context.Background(), r.pipelineOrganizationID, r.pipelineClusterID, req)
	if err != nil {
		return errors.Wrapf(err, "post node ready call failed.")
	}

	return nil
}
