# \GoogleApi

All URIs are relative to *http://localhost:9090*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListProjects**](GoogleApi.md#ListProjects) | **Get** /api/v1/orgs/{orgId}/cloud/google/projects | Retrieves projects visible for the user identified by the secret id



## ListProjects

> GoogleProjects ListProjects(ctx, orgId, secretId)

Retrieves projects visible for the user identified by the secret id

Retrieves projects visible by the user represented by the secretId header from the google cloud

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**secretId** | **string**| Secret identification. | 

### Return type

[**GoogleProjects**](GoogleProjects.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, 

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

