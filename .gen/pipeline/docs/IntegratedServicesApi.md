# \IntegratedServicesApi

All URIs are relative to *http://localhost:9090*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateIntegratedService**](IntegratedServicesApi.md#ActivateIntegratedService) | **Post** /api/v1/orgs/{orgId}/clusters/{id}/services/{serviceName} | Activate an integrated service
[**DeactivateIntegratedService**](IntegratedServicesApi.md#DeactivateIntegratedService) | **Delete** /api/v1/orgs/{orgId}/clusters/{id}/services/{serviceName} | Deactivate an integrated service
[**IntegratedServiceDetails**](IntegratedServicesApi.md#IntegratedServiceDetails) | **Get** /api/v1/orgs/{orgId}/clusters/{id}/services/{serviceName} | Get details of an integrated service
[**ListIntegratedServices**](IntegratedServicesApi.md#ListIntegratedServices) | **Get** /api/v1/orgs/{orgId}/clusters/{id}/services | List enabled integrated services of a cluster
[**UpdateIntegratedService**](IntegratedServicesApi.md#UpdateIntegratedService) | **Put** /api/v1/orgs/{orgId}/clusters/{id}/services/{serviceName} | Update an integrated service



## ActivateIntegratedService

> ActivateIntegratedService(ctx, orgId, id, serviceName, activateIntegratedServiceRequest)

Activate an integrated service

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**serviceName** | **string**| service name | 
**activateIntegratedServiceRequest** | [**ActivateIntegratedServiceRequest**](ActivateIntegratedServiceRequest.md)|  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeactivateIntegratedService

> DeactivateIntegratedService(ctx, orgId, id, serviceName)

Deactivate an integrated service

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**serviceName** | **string**| service name | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## IntegratedServiceDetails

> IntegratedServiceDetails IntegratedServiceDetails(ctx, orgId, id, serviceName)

Get details of an integrated service

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**serviceName** | **string**| service name | 

### Return type

[**IntegratedServiceDetails**](IntegratedServiceDetails.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, 

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListIntegratedServices

> map[string]IntegratedServiceDetails ListIntegratedServices(ctx, orgId, id)

List enabled integrated services of a cluster

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 

### Return type

[**map[string]IntegratedServiceDetails**](IntegratedServiceDetails.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, 

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateIntegratedService

> UpdateIntegratedService(ctx, orgId, id, serviceName, updateIntegratedServiceRequest)

Update an integrated service

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**serviceName** | **string**| service name | 
**updateIntegratedServiceRequest** | [**UpdateIntegratedServiceRequest**](UpdateIntegratedServiceRequest.md)|  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

