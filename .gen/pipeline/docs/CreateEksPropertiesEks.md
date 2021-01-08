# CreateEksPropertiesEks

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthConfig** | [**EksAuthConfig**](EKSAuthConfig.md) |  | [optional] 
**Version** | **string** |  | [optional] 
**EncryptionConfig** | [**[]EksEncryptionConfig**](EKSEncryptionConfig.md) | List of encryption config objects to define the encryption providers and their corresponding resources to encrypt. More information can be found at https://docs.aws.amazon.com/eks/latest/userguide/create-cluster.html and https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateCluster.html. | [optional] 
**LogTypes** | **[]string** |  | [optional] 
**NodePools** | [**map[string]EksNodePool**](EKSNodePool.md) |  | 
**Vpc** | [**EksVpc**](EKSVpc.md) |  | [optional] 
**RouteTableId** | **string** | Id of the RouteTable of the VPC to be used by subnets. This is used only when subnets are created into existing VPC. | [optional] 
**Subnets** | [**[]EksSubnet**](EKSSubnet.md) | Subnets for EKS master and worker nodes. All worker nodes will be launched in the same subnet (the first subnet in the list - which may not coincide with first subnet in the cluster create request payload as the deserialization may change the order) unless a subnet is specified for the workers that belong to a node pool at node pool level. | [optional] 
**Iam** | [**EksIam**](EKSIam.md) |  | [optional] 
**ApiServerAccessPoints** | **[]string** | List of access point types for the API server; public and private are the only valid values | [optional] [default to ["public"]]
**Tags** | **map[string]string** | User defined tags to be added to created AWS resources. Empty keys and values are not permitted. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


