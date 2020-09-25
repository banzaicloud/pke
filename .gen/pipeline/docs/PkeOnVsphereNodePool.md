# PkeOnVsphereNodePool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Roles** | **[]string** |  | 
**Labels** | **map[string]string** |  | [optional] 
**Size** | **int32** |  | [optional] 
**Vcpu** | **int32** | Number of VCPUs to attach to each node. | 
**Ram** | **int32** | MiBs of RAM to attach to each node. | 
**Template** | **string** | Name of VM template available on vSphere to clone as the base of nodes. Overrides default value from the main cluster secret. | [optional] 
**AdminUsername** | **string** | Name of admin user to deploy the generated SSH public key for. No key will be deployed if omitted. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


