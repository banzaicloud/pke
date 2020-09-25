# \ProcessesApi

All URIs are relative to *http://localhost:9090*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelProcess**](ProcessesApi.md#CancelProcess) | **Post** /api/v1/orgs/{orgId}/processes/{id}/cancel | Cancel a process in Pipeline
[**GetProcess**](ProcessesApi.md#GetProcess) | **Get** /api/v1/orgs/{orgId}/processes/{id} | Get a process in Pipeline
[**ListProcesses**](ProcessesApi.md#ListProcesses) | **Get** /api/v1/orgs/{orgId}/processes | List processes in Pipeline



## CancelProcess

> CancelProcess(ctx, orgId, id)

Cancel a process in Pipeline

Cancel a process in Pipeline

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **string**| Process id | 

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


## GetProcess

> Process GetProcess(ctx, orgId, id)

Get a process in Pipeline

Get a process in Pipeline

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **string**| Process id | 

### Return type

[**Process**](Process.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, 

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListProcesses

> []Process ListProcesses(ctx, orgId, optional)

List processes in Pipeline

List processes in Pipeline

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
 **optional** | ***ListProcessesOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ListProcessesOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **type_** | **optional.String**| Type of processes to query | 
 **resourceId** | **optional.String**| The id of the resource to list processes for | 
 **parentId** | **optional.String**| The id of the parent process | 
 **status** | [**optional.Interface of ProcessStatus**](.md)| The status of processes to query | 

### Return type

[**[]Process**](Process.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, 

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

