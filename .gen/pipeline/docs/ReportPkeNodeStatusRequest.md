# ReportPkeNodeStatusRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | name of node | [optional] 
**NodePool** | **string** | name of nodepool | [optional] 
**Ip** | **string** | ip address of node (where the other nodes can reach it) | [optional] 
**Message** | **string** | detailed description about the current bootstrapping status (including the cause of the failure) | [optional] 
**Phase** | **string** | the current phase of the bootstrap process | [optional] 
**Final** | **bool** | if this is the final status report, that describes the conclusion of the whole process | [optional] 
**Status** | [**ProcessStatus**](ProcessStatus.md) |  | [optional] 
**Timestamp** | Pointer to [**time.Time**](time.Time.md) | exact time of event | [optional] 
**ProcessId** | **string** | ID of the process registered earlier (register new process if empty) | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


