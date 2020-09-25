# EnableArkRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cloud** | **string** |  | 
**BucketName** | **string** |  | 
**Schedule** | **string** |  | 
**Ttl** | **string** |  | 
**SecretId** | **string** |  | 
**Location** | **string** |  | [optional] 
**UseClusterSecret** | **bool** | relevant only in case of Amazon clusters. By default set to false in which case you must add snapshot permissions to your node instance role. Should you set to true Pipeline will deploy your cluster secret to the cluster. | [optional] 
**ServiceAccountRoleARN** | **string** | relevant only in case of Amazon clusters. This a third option to give permissions for volume snapshots to Velero, besides the default NodeInstance role or cluster secret deployment. | [optional] 
**StorageAccount** | **string** | required only case of Azure | [optional] 
**ResourceGroup** | **string** | required only case of Azure | [optional] 
**Labels** | [**Labels**](Labels.md) |  | [optional] 
**Options** | [**BackupOptions**](BackupOptions.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


