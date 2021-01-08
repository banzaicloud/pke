# EksNodePool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceType** | **string** |  | 
**SpotPrice** | **string** |  | 
**Autoscaling** | **bool** |  | [optional] 
**Count** | **int32** |  | [optional] 
**MinCount** | **int32** |  | 
**MaxCount** | **int32** |  | 
**Labels** | **map[string]string** |  | [optional] 
**VolumeSize** | **int32** | Size of the EBS volume in GBs of the nodes in the pool. | [optional] 
**Image** | **string** |  | [optional] 
**Subnet** | [**EksSubnet**](EKSSubnet.md) |  | [optional] 
**SecurityGroups** | **[]string** | List of additional custom security groups for all nodes in the pool. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


