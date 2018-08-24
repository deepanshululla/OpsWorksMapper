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
		//os.Exit(1)
	}
	ioutil.WriteFile(fileName,[]byte(b),0666);
}

func ReadJsonMap(fileName string) map[string]map[string]string{
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Successfully Opened file")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var hashMap map[string]map[string]string
	json.Unmarshal(byteValue, &hashMap)

	return hashMap
}