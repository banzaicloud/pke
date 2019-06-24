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

package pipeline

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/banzaicloud/pke/.gen/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/banzaicloud/pke/cmd/pke/app/util/transport"
	"github.com/banzaicloud/pke/cmd/pke/app/util/validator"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// Client initializes Pipeline API client
func Client(out io.Writer, endpoint, token string, insecure bool) *pipeline.APIClient {
	config := pipeline.NewConfiguration()
	config.BasePath = endpoint
	config.UserAgent = "pke/1.0.0/go"

	httpClient := http.Client{
		Timeout: 24 * time.Hour,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &httpClient)

	config.HTTPClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	// Since transport.NewRetryTransport is added, this timeout will affect only the cumulated retry calls.
	tl := transport.NewLogger(out, config.HTTPClient.Transport)
	config.HTTPClient.Transport = transport.NewRetryTransport(tl)

	return pipeline.NewAPIClient(config)
}

// CommandArgs extracts args needed for Pipeline API client.
func CommandArgs(cmd *cobra.Command) (endpoint, token string, insecure bool, orgID, clusterID int32, err error) {
	endpoint, err = cmd.Flags().GetString(constants.FlagPipelineAPIEndpoint)
	if err != nil {
		return
	}

	token, err = cmd.Flags().GetString(constants.FlagPipelineAPIToken)
	if err != nil {
		return
	}

	insecure, err = cmd.Flags().GetBool(constants.FlagPipelineAPIInsecure)
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
	endpoint, token, _, orgID, clusterID, err := CommandArgs(cmd)
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
func ValidArgs(endpoint, token string, insecure bool, orgID, clusterID int32) error {
	return validator.NotEmpty(map[string]interface{}{
		constants.FlagPipelineAPIEndpoint:    endpoint,
		constants.FlagPipelineAPIToken:       token,
		constants.FlagPipelineOrganizationID: orgID,
		constants.FlagPipelineClusterID:      clusterID,
	})
}

func NodeJoinArgs(out io.Writer, cmd *cobra.Command) (apiServerHostPort, kubeadmToken, caCertHash string, err error) {
	if !Enabled(cmd) {
		return
	}
	endpoint, token, insecure, orgID, clusterID, err := CommandArgs(cmd)
	if err != nil {
		return
	}

	// Pipeline client.
	c := Client(out, endpoint, token, insecure)

	var b pipeline.GetClusterBootstrapResponse
	b, _, err = c.ClustersApi.GetClusterBootstrap(context.Background(), orgID, clusterID)
	if err != nil {
		return
	}
	apiServerHostPort = b.MasterAddress
	kubeadmToken = b.Token
	caCertHash = b.DiscoveryTokenCaCertHash
	return
}
