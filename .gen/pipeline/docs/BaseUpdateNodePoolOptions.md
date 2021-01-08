# BaseUpdateNodePoolOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaxSurge** | **int32** | Maximum number of extra nodes that can be created during the update. | [optional] [default to 0]
**MaxBatchSize** | **int32** | Maximum number of nodes that can be updated simultaneously. | [optional] [default to 2]
**MaxUnavailable** | **int32** | Maximum number of nodes that can be unavailable during the update. | [optional] [default to 0]
**Drain** | [**UpdateNodePoolDrainOptions**](UpdateNodePoolDrainOptions.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


