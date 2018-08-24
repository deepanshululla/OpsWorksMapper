package main

import (
	"fmt"
	"sync"
)



func opsWorksChannelInstanceName(fileName string){
	opsWorksStructGrp:=InitServiceGroup()
	fmt.Println("Aws Call for ec2")
	l:=opsWorksStructGrp.GetInstanceIdNameMap()
	ConvertToJsonFile(CreateStructMapForMap(l),fileName)
	fmt.Println("Successfuly wrote ec2 mapping to file.Total Entries =", len(l))
}


func opsWorksChannelElbStackName(fileName string){
	fmt.Println("AWS call for elb")
	opsWorksStructGrp:=InitServiceGroup()
	l:=opsWorksStructGrp.GetELbStackNameMap()
	ConvertToJsonFile(CreateStructMapForMap(l),fileName)
	fmt.Println("Successfuly wrote elb mapping to file.Total Entries =",len(l))
}

func OpsWorksInstanceIdNameMapper(fileName string){
	opsWorksStructGrp:=InitServiceGroup()
	fmt.Println("Aws Call for opsworks")
	instMap:=opsWorksStructGrp.GetOpsworksInstanceIdNameMap()
	ConvertToJsonFile(CreateStructMapForMap(instMap),fileName)
	fmt.Println("Successfuly wrote opsworks mapping to file. Total Entries =",len(instMap))
}


func Syncher(ec2MapFile,elbMapFile,opsMapFile string){
	for{
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			opsWorksChannelInstanceName(ec2MapFile)
			wg.Done()
		}()
		go func() {
			opsWorksChannelElbStackName(elbMapFile)
			wg.Done()
		}()
		go func() {
			OpsWorksInstanceIdNameMapper(opsMapFile)
			wg.Done()
		}()
		wg.Wait()
	}
}

