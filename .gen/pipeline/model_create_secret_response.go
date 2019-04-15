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

type CreateSecretResponse struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Id string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
	Version int32 `json:"version,omitempty"`
}
