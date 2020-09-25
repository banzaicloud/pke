# AmazonAutoScalingGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Image** | **string** |  | 
**VolumeSize** | **int32** | Size of root EBS volume to attach to the nodes in GBs. Zero means that the size is determined automatically. | [optional] 
**Zones** | **[]string** |  | 
**InstanceType** | **string** |  | 
**LaunchConfigurationName** | **string** |  | 
**LaunchTemplate** | [**map[string]interface{}**](.md) |  | [optional] 
**VpcID** | **string** |  | 
**SecurityGroupID** | **string** |  | 
**Subnets** | **[]string** |  | 
**Tags** | **map[string]map[string]interface{}** |  | 
**SpotPrice** | **string** |  | 
**Size** | [**AmazonAutoScalingGroupSize**](AmazonAutoScalingGroup_size.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


