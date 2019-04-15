/*
 * Pipeline API
 *
 * Pipeline v0.3.0 swagger
 *
 * API version: 0.17.0
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pipeline

type CreateObjectStoreBucketProperties struct {
	Alibaba CreateAmazonObjectStoreBucketProperties `json:"alibaba,omitempty"`
	Azure CreateAzureObjectStoreBucketProperties `json:"azure,omitempty"`
	Google CreateGoogleObjectStoreBucketProperties `json:"google,omitempty"`
	Oracle CreateOracleObjectStoreBucketProperties `json:"oracle,omitempty"`
}
