package ready

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/banzaicloud/pipeline/client"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/util/network"
	"github.com/banzaicloud/pke/cmd/pke/app/util/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/pkg/errors"
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
	pipelineAPIEndpoint    string
	pipelineAPIToken       string
	pipelineOrganizationID int32
	pipelineClusterID      int32
	pipelineNodepool       string
}

func NewCommand(out io.Writer, role Role) *cobra.Command {
	return phases.NewCommand(out, &Ready{role: role})
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
	flags.Int32(constants.FlagPipelineOrganizationID, 0, "Organization ID to use with Pipeline API")
	flags.Int32(constants.FlagPipelineClusterID, 0, "Cluster ID to use with Pipeline API")
	flags.String(constants.FlagPipelineNodepool, "", "name of the nodepool the node belongs to")
	flags.String(constants.FlagInfrastructureCIDR, "192.168.64.0/20", "network CIDR for the actual machine")
}

func (r *Ready) Validate(cmd *cobra.Command) error {
	// Run is optional, only validate if Pipeline credentials are present.
	if !pipeline.Enabled(cmd) {
		return nil
	}

	var err error
	if r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineOrganizationID, r.pipelineClusterID, err = pipeline.CommandArgs(cmd); err != nil {
		return err
	}

	if err = pipeline.ValidArgs(r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineOrganizationID, r.pipelineClusterID); err != nil {
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
	_, _ = fmt.Fprintf(out, "[RUNNING] %s\n", r.Use())

	if err := pipeline.ValidArgs(r.pipelineAPIEndpoint, r.pipelineAPIToken, r.pipelineOrganizationID, r.pipelineClusterID); err != nil {
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
	c := pipeline.Client(out, r.pipelineAPIEndpoint, r.pipelineAPIToken)
	req := client.PostReadyPkeNodeRequest{
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

	_, resp, err := c.ClustersApi.PostReadyPKENode(context.Background(), r.pipelineOrganizationID, r.pipelineClusterID, req)
	if err != nil {
		return errors.Wrapf(err, "post node ready call failed. http status code: %d", resp.StatusCode)
	}

	return nil
}
