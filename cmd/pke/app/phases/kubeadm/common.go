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

package kubeadm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/pkg/errors"
)

const (
	urlAzureAZ = "http://169.254.169.254/metadata/instance?api-version=2018-10-01"
)

func WriteKubeadmAzureConfig(out io.Writer, filename, cloudProvider, tenantID, subnetName, securityGroupName, vnetName, vnetResourceGroup, vmType, loadBalancerSku, routeTableName string) error {
	if cloudProvider == constants.CloudProviderAzure {
		if http.DefaultClient.Timeout < 10*time.Second {
			http.DefaultClient.Timeout = 10 * time.Second
		}

		req, err := http.NewRequest(http.MethodGet, urlAzureAZ, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Metadata", "true")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("failed to get azure availability zone. http status code: %d", resp.StatusCode))
		}
		defer func() { _ = resp.Body.Close() }()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		type metadata struct {
			Compute struct {
				AZEnvironment     string `json:"azEnvironment"`
				Location          string `json:"location"`
				ResourceGroupName string `json:"resourceGroupName"`
				SubscriptionId    string `json:"subscriptionId"`
			} `json:"compute"`
		}
		var r metadata
		if err := json.Unmarshal(b, &r); err != nil {
			return errors.Wrap(err, "failed to parse response")
		}

		if vmType == "" {
			vmType = "standard"
		}
		if loadBalancerSku == "" {
			loadBalancerSku = "basic"
		}

		conf := `{
    "cloud":"{{ .Cloud }}",
    "tenantId": "{{ .TenantId }}",
    "subscriptionId": "{{ .SubscriptionId }}",
    "resourceGroup": "{{ .ResourceGroup }}",
    "location": "{{ .Location }}",
    "subnetName": "{{ .SubnetName }}",
    "securityGroupName": "{{ .SecurityGroupName }}",
    "vnetName": "{{ .VNetName }}",
    "vnetResourceGroup": "{{ .VNetResourceGroup }}",
    "vmType": "{{ .VMType }}",
    "loadBalancerSku": "{{ .LoadBalancerSku }}",
    "routeTableName": "{{ .RouteTableName }}",
    "cloudProviderBackoff": false,
    "useManagedIdentityExtension": true,
    "useInstanceMetadata": true
}`

		tmpl, err := template.New("azure-config").Parse(conf)
		if err != nil {
			return err
		}

		w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
		if err != nil {
			return err
		}
		defer func() { _ = w.Close() }()

		type data struct {
			Cloud             string
			TenantId          string
			SubscriptionId    string
			ResourceGroup     string
			Location          string
			SubnetName        string
			SecurityGroupName string
			VNetName          string
			VNetResourceGroup string
			VMType            string
			LoadBalancerSku   string
			RouteTableName    string
		}

		d := data{
			Cloud:             r.Compute.AZEnvironment,
			TenantId:          tenantID,
			SubscriptionId:    r.Compute.SubscriptionId,
			ResourceGroup:     r.Compute.ResourceGroupName,
			Location:          r.Compute.Location,
			SubnetName:        subnetName,
			SecurityGroupName: securityGroupName,
			VNetName:          vnetName,
			VNetResourceGroup: vnetResourceGroup,
			VMType:            vmType,
			LoadBalancerSku:   loadBalancerSku,
			RouteTableName:    routeTableName,
		}

		return tmpl.Execute(w, d)
	}

	return nil
}
