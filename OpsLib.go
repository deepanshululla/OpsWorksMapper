package main

import (
	"github.com/aws/aws-sdk-go/service/opsworks"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)



type OpsWorksStruct struct {
	service *opsworks.OpsWorks
}

type OpsWorksStructGrp struct {
	structGroup []OpsWorksStruct
}
type opsClient interface {
	GetStackIdNameMap() map[string]map[string]string
	GetInstanceIdNameMap() map[string]map[string]string
	GetELbStackNameMap() map[string]map[string]string
}

var regions=[]string{"ap-south-1", "eu-west-3","eu-west-2",
	"eu-west-1", "ap-northeast-2", "ap-northeast-1",
	"sa-east-1", "ca-central-1", "ap-southeast-1",
	"ap-southeast-2", "eu-central-1", "us-east-1",
	"us-east-2" ,"us-west-1" ,"us-west-2"}


func InitService(region string) OpsWorksStruct {
	opsworksRegionMap:=map[string]string{
		"us-east-2": "us-east-2",
		"us-east-1": "us-east-1",
		"us-west-1": "us-east-1",
		"us-west-2": "us-east-1",
		"eu-west-1": "us-east-1",
		"eu-west-2": "eu-west-2",
		"ap-south-1": "ap-south-1",
		"ap-northeast-1": "us-east-1",
		"ap-northeast-2": "us-east-1",
		"ap-southeast-2": "us-east-1",
	}
	// Load session from shared config
	value,ok:=opsworksRegionMap[region]
	if ok{
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(value),
		}))
		opsworksServ:=opsworks.New(sess)
		return OpsWorksStruct{opsworksServ}
	} else {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
		opsworksServ:=opsworks.New(sess)
		return OpsWorksStruct{opsworksServ}
	}
}

func InitServiceGroup() OpsWorksStructGrp{
	opsGrp:=[]OpsWorksStruct{}
	for _,region:=range regions{
		opsGrp=append(opsGrp,InitService(region))
	}
	return OpsWorksStructGrp{opsGrp}
}

/*
Mapping between Opsworks Stack id and Cluster Name.
 */
func (opsWorksStruct OpsWorksStruct) GetStackIdNameMap() map[string]map[string]string{
	opsworksServ:=opsWorksStruct.service
	stackMap := map[string]map[string]string{}
	resp,err:=opsworksServ.DescribeStacks(nil)
	if err!=nil{
		fmt.Println("Error",err)
		//os.Exit(1)
	}
	for _,stack:=range resp.Stacks{
		stackId:=*stack.StackId
		stackMap[stackId]=map[string]string{}
		stackMap[stackId]["Region"]=*stack.Region
		stackMap[stackId]["Name"]=*stack.Name
		stackMap[stackId]["VpcId"]=*stack.VpcId
	}
	return stackMap
}

func (opsWorksStructGrp OpsWorksStructGrp) GetStackIdNameMap() map[string]map[string]string{
	stackMap := map[string]map[string]string{}
	for _,opsStruct:=range opsWorksStructGrp.structGroup{
		curStack:=opsStruct.GetStackIdNameMap()
		for key,value:= range curStack{
			stackMap[key]=value
		}
	}
	return stackMap
}
/*
Mapping between EC2 Instance Ids and Opswork Instance Ids as well as mapping between
EC2 instance ids and cluster name
 */
func (opsWorksStruct OpsWorksStruct) GetInstanceIdNameMap() map[string]map[string]string{
	opsworksServ:=opsWorksStruct.service
	instMap := map[string]map[string]string{}
	stackMap:=opsWorksStruct.GetStackIdNameMap()
	for stackId,stackValue := range stackMap{
		stackName:=stackValue["Name"]
		stackRegion:=stackValue["Region"]
		inp:=&opsworks.DescribeInstancesInput{
			StackId:&stackId,
		}
		resp,err:=opsworksServ.DescribeInstances(inp)
		if err!=nil{
			fmt.Println("Error",err)
			//os.Exit(1)
		}
		for _,instancePointer:=range resp.Instances{
			instance:=*instancePointer
			if instance.Ec2InstanceId!=nil{
				ec2InstanceId:=*instance.Ec2InstanceId
				instMap[ec2InstanceId]=map[string]string{}
				instMap[ec2InstanceId]["Hostname"]=*instance.Hostname
				instMap[ec2InstanceId]["OpsworksId"]=*instance.InstanceId
				instMap[ec2InstanceId]["StackId"]=*instance.StackId
				instMap[ec2InstanceId]["StackName"]=stackName
				instMap[ec2InstanceId]["Region"]=stackRegion
				instMap[ec2InstanceId]["InstanceType"]=*instance.InstanceType
			}
		}
	}
	return instMap
}

func (opsWorksStructGrp OpsWorksStructGrp) GetInstanceIdNameMap() map[string]map[string]string{
	instanceMap := map[string]map[string]string{}
	for _,opsStruct:=range opsWorksStructGrp.structGroup{
		curInstance:=opsStruct.GetInstanceIdNameMap()
		for key,value:= range curInstance{
			instanceMap[key]=value
		}
	}
	return instanceMap
}
/*
Mapping Between ELB and Clusters
 */
func (opsWorksStruct OpsWorksStruct) GetELbStackNameMap() map[string]map[string]string{
	opsworksServ:=opsWorksStruct.service
	elbMap := map[string]map[string]string{}
	stackMap:=opsWorksStruct.GetStackIdNameMap()
	for stackId,stackValue := range stackMap{
		stackName:=stackValue["Name"]
		inp:=&opsworks.DescribeElasticLoadBalancersInput{
			StackId:&stackId,
		}
		resp,err:=opsworksServ.DescribeElasticLoadBalancers(inp)
		if err!=nil{
			fmt.Println("Error",err)
			//os.Exit(1)
		}
		for _,elb:=range resp.ElasticLoadBalancers{
			lbName:=*elb.ElasticLoadBalancerName
			elbMap[lbName]=map[string]string{}
			elbMap[lbName]["StackName"]=stackName
			elbMap[lbName]["Region"]=*elb.Region
			elbMap[lbName]["StackId"]=*elb.StackId
			elbMap[lbName]["DNSName"]=*elb.DnsName
		}
	}
	return elbMap
}


func (opsWorksStructGrp OpsWorksStructGrp) GetELbStackNameMap() map[string]map[string]string{
	elbMap := map[string]map[string]string{}
	for _,opsStruct:=range opsWorksStructGrp.structGroup{
		curInstance:=opsStruct.GetELbStackNameMap()
		for key,value:= range curInstance{
			elbMap[key]=value
		}
	}
	return elbMap
}


/*
Mapping between opsworks Instance Ids and  Instance names
 */
func (opsWorksStruct OpsWorksStruct) GetOpsworksInstanceIdNameMap() map[string]map[string]string{
	opsworksServ:=opsWorksStruct.service
	instMap := map[string]map[string]string{}
	stackMap:=opsWorksStruct.GetStackIdNameMap()
	for stackId,stackValue := range stackMap{
		stackName:=stackValue["Name"]
		stackRegion:=stackValue["Region"]
		inp:=&opsworks.DescribeInstancesInput{
			StackId:&stackId,
		}
		resp,err:=opsworksServ.DescribeInstances(inp)
		if err!=nil{
			fmt.Println("Error From Opsworks:",err)
			//os.Exit(1)
		}
		for _,instancePointer:=range resp.Instances{
			instance:=*instancePointer
			if instance.InstanceId!=nil{
				opsInstanceId:=*instance.InstanceId
				instMap[opsInstanceId]=map[string]string{}
				instMap[opsInstanceId]["Hostname"]=*instance.Hostname
				//instMap[opsInstanceId]["EC2InstanceId"]=*instance.Ec2InstanceId
				instMap[opsInstanceId]["StackId"]=*instance.StackId
				instMap[opsInstanceId]["StackName"]=stackName
				instMap[opsInstanceId]["Region"]=stackRegion
				instMap[opsInstanceId]["InstanceType"]=*instance.InstanceType
			}

		}
	}
	return instMap
}

func (opsWorksStructGrp OpsWorksStructGrp) GetOpsworksInstanceIdNameMap() map[string]map[string]string{
	instanceMap := map[string]map[string]string{}
	for _,opsStruct:=range opsWorksStructGrp.structGroup{
		curInstance:=opsStruct.GetOpsworksInstanceIdNameMap()
		for key,value:= range curInstance{
			instanceMap[key]=value
		}
	}
	return instanceMap
}