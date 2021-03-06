/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments. 
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pipeline
// Error Generic error object.
type Error struct {
	// A URI reference [RFC3986] that identifies the problem type.
	Type string `json:"type,omitempty"`
	// A short, human-readable summary of the problem type.
	Title string `json:"title,omitempty"`
	// The HTTP status code ([RFC7231], Section 6) generated by the origin server for this occurrence of the problem.
	Status int32 `json:"status,omitempty"`
	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`
	// A URI reference that identifies the specific occurrence of the problem.
	Instance string `json:"instance,omitempty"`
	// HTTP status code. Deprecated: use status instead.
	Code int32 `json:"code,omitempty"`
	// Error message. Deprecated: use detail instead.
	Message string `json:"message,omitempty"`
	// Error message. Deprecated: use title instead.
	Error string `json:"error,omitempty"`
}
