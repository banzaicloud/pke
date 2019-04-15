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

type ProfileListResponse struct {
	Name string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
	Cloud string `json:"cloud,omitempty"`
	// The lifespan of the cluster expressed in minutes after which it is automatically deleted. Zero value means the cluster is never automatically deleted.
	TtlMinutes int32 `json:"ttlMinutes,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}
