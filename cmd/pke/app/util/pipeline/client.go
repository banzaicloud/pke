package pipeline

import (
	"io"
	"time"

	"github.com/banzaicloud/pipeline/client"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// Client initializes Pipeline API client
func Client(out io.Writer, endpoint, token string) *client.APIClient {
	config := client.NewConfiguration()
	config.BasePath = endpoint
	config.UserAgent = "banzai-cli/1.0.0/go"
	config.HTTPClient = oauth2.NewClient(nil, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	config.HTTPClient.Timeout = 30 * time.Second
	config.HTTPClient.Transport = &transportLogger{
		roundTripper: config.HTTPClient.Transport,
		output:       out,
	}

	return client.NewAPIClient(config)
}

// CommandArgs extracts args needed for Pipeline API client.
func CommandArgs(cmd *cobra.Command) (endpoint, token string, orgID, clusterID int32, err error) {
	endpoint, err = cmd.Flags().GetString(constants.FlagPipelineAPIEndpoint)
	if err != nil {
		return
	}

	token, err = cmd.Flags().GetString(constants.FlagPipelineAPIToken)
	if err != nil {
		return
	}

	orgID, err = cmd.Flags().GetInt32(constants.FlagPipelineOrganizationID)
	if err != nil {
		return
	}

	clusterID, err = cmd.Flags().GetInt32(constants.FlagPipelineClusterID)
	if err != nil {
		return
	}

	return
}

func Enabled(cmd *cobra.Command) bool {
	endpoint, token, orgID, clusterID, err := CommandArgs(cmd)
	if err != nil {
		// TODO: remove this silent error.
		return false
	}

	return validator.Empty(map[string]interface{}{
		constants.FlagPipelineAPIEndpoint:    endpoint,
		constants.FlagPipelineAPIToken:       token,
		constants.FlagPipelineOrganizationID: orgID,
		constants.FlagPipelineClusterID:      clusterID,
	}) != nil
}

// ValidArgs ensures all Pipeline API args are present.
func ValidArgs(endpoint, token string, orgID, clusterID int32) error {
	return validator.NotEmpty(map[string]interface{}{
		constants.FlagPipelineAPIEndpoint:    endpoint,
		constants.FlagPipelineAPIToken:       token,
		constants.FlagPipelineOrganizationID: orgID,
		constants.FlagPipelineClusterID:      clusterID,
	})
}
