package main

import (
	"fmt"
	"sync"
)



func opsWorksChannelInstanceName(fileName string, c chan string){

	opsWorksStructGrp:=InitServiceGroup()
	fmt.Println("Aws Call for ec2")
	l:=opsWorksStructGrp.GetInstanceIdNameMap()

	ConvertToJsonFile(CreateStructMapForMap(l),fileName)
	fmt.Println("Successfuly wrote ec2 mapping to file.Total Entries =", len(l))
	c<-fileName
}


func opsWorksChannelElbStackName(fileName string, c chan string){
	fmt.Println("AWS call for elb")
	opsWorksStructGrp:=InitServiceGroup()
	l:=opsWorksStructGrp.GetELbStackNameMap()

	ConvertToJsonFile(CreateStructMapForMap(l),fileName)
	fmt.Println("Successfuly wrote elb mapping to file.Total Entries =",len(l))
	c<-fileName
}

func OpsWorksInstanceIdNameMapper(fileName string, c chan string){
	//time.Sleep(180*time.Second)
	opsWorksStructGrp:=InitServiceGroup()
	fmt.Println("Aws Call for opsworks")

	instMap:=opsWorksStructGrp.GetOpsworksInstanceIdNameMap()
	ConvertToJsonFile(CreateStructMapForMap(instMap),fileName)

	fmt.Println("Successfuly wrote opsworks mapping to file. Total Entries =",len(instMap))

	c<-fileName

}


func Syncher(ec2MapFile,elbMapFile,opsMapFile string){
	for{

		chan1:=make(chan string)
		chan2:=make(chan string)
		chan3:=make(chan string)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			opsWorksChannelInstanceName(ec2MapFile,chan1)
			wg.Done()
		}()
		go func() {
			opsWorksChannelElbStackName(elbMapFile,chan2)
			wg.Done()
		}()
		go func() {
			OpsWorksInstanceIdNameMapper(opsMapFile,chan3)
			wg.Done()
		}()
		wg.Wait()
	}

}

