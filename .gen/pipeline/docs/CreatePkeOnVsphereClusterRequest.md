# CreatePkeOnVsphereClusterRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**SecretId** | **string** |  | [optional] 
**SecretName** | **string** |  | [optional] 
**SshSecretId** | **string** |  | [optional] 
**ScaleOptions** | [**ScaleOptions**](ScaleOptions.md) |  | [optional] 
**Type** | **string** |  | 
**Kubernetes** | [**CreatePkeClusterKubernetes**](CreatePKEClusterKubernetes.md) |  | 
**Proxy** | [**PkeClusterHttpProxy**](PKEClusterHTTPProxy.md) |  | [optional] 
**StorageSecretId** | **string** | Secret ID used to setup VSphere storage classes. Overrides the default settings in main cluster secret. | [optional] 
**StorageSecretName** | **string** | Secret name used to setup VSphere storage classes. Overrides default value from the main cluster secret. | [optional] 
**Folder** | **string** | Folder to create nodes in. Overrides default value from the main cluster secret. | [optional] 
**Datastore** | **string** | Name of datastore or datastore cluster to place VM disks on. Overrides default value from the main cluster secret. | [optional] 
**ResourcePool** | **string** | Virtual machines will be created in this resource pool. Overrides default value from the main cluster secret. | [optional] 
**Nodepools** | [**[]PkeOnVsphereNodePool**](PKEOnVsphereNodePool.md) |  | [optional] 
**LoadBalancerIPRange** | **string** | IPv4 range to allocate addresses for LoadBalancer Services (MetalLB) | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


