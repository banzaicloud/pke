# EksNodePoolAllOf

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Autoscaling** | [**NodePoolAutoScaling**](NodePoolAutoScaling.md) |  | [optional] 
**VolumeSize** | **int32** | Size of the EBS volume in GBs of the nodes in the pool. | [optional] 
**InstanceType** | **string** | Machine instance type. | 
**Image** | **string** | Instance AMI. | [optional] 
**SpotPrice** | **string** | The upper limit price for the requested spot instance. If this field is left empty or 0 passed in on-demand instances used instead of spot instances. | [optional] 
**SubnetId** | **string** |  | [optional] 
**SecurityGroups** | **[]string** | List of additional custom security groups for all nodes in the pool. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


