package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type structMap struct {
	dictionary map[string]map[string]string
}

func CreateStructMapForMap(d map[string]map[string]string) structMap{
	m:=structMap{d}
	return m;
}


func ConvertToJsonFile(sm structMap, fileName string){
	b, err := json.MarshalIndent(sm.dictionary, "", "  ")
	if (err!=nil){
		fmt.Println(err)
	}
	ioutil.WriteFile(fileName,[]byte(b),0666);
}

func ReadJsonMap(fileName string) map[string]map[string]string{
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var hashMap map[string]map[string]string
	json.Unmarshal(byteValue, &hashMap)

	return hashMap
}