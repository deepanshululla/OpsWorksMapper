package main


type MapperStruct struct {
	elbMap map[string]map[string]string
	opsMap map[string]map[string]string
	ec2Map map[string]map[string]string
}


func main()  {
	ec2MapFile:="$GOPATH/src/github.com/influxdata/telegraf/instanceIdNameMap.json"
	elbMapFile:="$GOPATH/src/github.com/influxdata/telegraf/elbClusterNameMap.json"
	opsMapFile:="$GOPATH/src/github.com/influxdata/telegraf/opsWorksInstanceIdNameMap.json"
	Syncher(ec2MapFile,elbMapFile,opsMapFile)
}


