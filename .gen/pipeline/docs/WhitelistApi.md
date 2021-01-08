# \WhitelistApi

All URIs are relative to *http://localhost:9090*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateWhitelists**](WhitelistApi.md#CreateWhitelists) | **Post** /api/v1/orgs/{orgId}/clusters/{id}/whitelists | Create Whitelisted deployment
[**DeleteWhitelist**](WhitelistApi.md#DeleteWhitelist) | **Delete** /api/v1/orgs/{orgId}/clusters/{id}/whitelists/{name} | Delete Whitelisted deployment
[**ListWhitelists**](WhitelistApi.md#ListWhitelists) | **Get** /api/v1/orgs/{orgId}/clusters/{id}/whitelists | List Whitelisted deployments



## CreateWhitelists

> CreateWhitelists(ctx, orgId, id, releaseWhiteListItem)

Create Whitelisted deployment

Create Whitelisted deployment

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**releaseWhiteListItem** | [**ReleaseWhiteListItem**](ReleaseWhiteListItem.md)|  | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteWhitelist

> DeleteWhitelist(ctx, orgId, id, name)

Delete Whitelisted deployment

Delete Whitelisted deployment

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 
**name** | **string**| Selected whitelist identification | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWhitelists

> []ReleaseWhiteListItem ListWhitelists(ctx, orgId, id)

List Whitelisted deployments

List Whitelisted deployments

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**orgId** | **int32**| Organization identifier | 
**id** | **int32**| Cluster identifier | 

### Return type

[**[]ReleaseWhiteListItem**](ReleaseWhiteListItem.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

