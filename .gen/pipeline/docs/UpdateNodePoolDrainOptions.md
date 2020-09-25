# UpdateNodePoolDrainOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timeout** | **int32** | How long should drain wait for pod eviction (in seconds) | [optional] [default to 0]
**FailOnError** | **bool** | Whether the process should fail if draining fails/times out. | [optional] [default to false]
**PodSelector** | **string** | Only evict those pods that matches this selector. | [optional] [default to ]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


