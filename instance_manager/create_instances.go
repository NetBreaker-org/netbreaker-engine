package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NetBreaker-org/netbreaker-engine/instance_manager/generals"
)

type SeedFileStruct struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type InstancesJson struct {
	Name   string `json:"instance_identifier"`
	Dir    string `json:"instance_dir"`
	Schema string `json:"schema_used"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// create an instance folder and add entry to instances.json
func CreateInstance(seedFilePath string) int {
	//check if seed file exists
	if fileExists(seedFilePath) == true {
		//read seed file
		byteSlice, err := os.ReadFile(seedFilePath)
		if err != nil {
			return 1
		}

		//variable to store parsed json data from seed file
		var storeSeedFileJson SeedFileStruct

		//unmarshal the data in seed file and load to storeSeedFileJson
		err = json.Unmarshal(byteSlice, &storeSeedFileJson)
		if err != nil {
			fmt.Printf("Error in making instance, error : %v\n", err)
			return 1
		}

		//open instances.json file and load to instances_json variable
		var instance_json []InstancesJson
		instFileByteSlice, err := os.ReadFile(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json"))
		if err != nil {
			fmt.Printf("Error in making instance : %s, error : %v\n", storeSeedFileJson.Name, err)
			return 1
		}

		err = json.Unmarshal(instFileByteSlice, &instance_json)
		if err != nil {
			fmt.Printf("Error in making instance : %s, error : %v\n", storeSeedFileJson.Name, err)
			return 1
		}

		//check if there is already an entry for this instance name
		for i := range instance_json {
			if instance_json[i].Name == storeSeedFileJson.Name {
				fmt.Printf("Error in making instance : %s, error : %v\n", storeSeedFileJson.Name, "An instance with the provided name already exists in the instances.json file")
				return 1
			}
		}

		//make directory for the instance
		//instance folder's file path
		instance_filepath := filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances", storeSeedFileJson.Name)

		//check if instance folder already exists
		instanceExistsStat := fileExists(instance_filepath)
		if instanceExistsStat == true {
			fmt.Printf("Error in making instance : %s, error : instance already exists; if the instance is no longer used, delete the instance first\n", storeSeedFileJson.Name)
			return 1
		}

		err = os.MkdirAll(instance_filepath, 0755)
		if err != nil {
			fmt.Printf("Error in making instance : %s, error : %v\n", storeSeedFileJson.Name, err)
			return 1
		}

		//copy schema to the instance folder
		schemaPath := filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/seed_schemas/", storeSeedFileJson.Schema)
		err = os.CopyFS(instance_filepath, os.DirFS(schemaPath))
		if err != nil {
			fmt.Printf("Error in making instance : %s, error : %v\n", storeSeedFileJson.Name, err)
			return 1
		}

		//add entry to instances.json
		var toAppendToInstances_json InstancesJson
		toAppendToInstances_json.Name = storeSeedFileJson.Name
		toAppendToInstances_json.Schema = storeSeedFileJson.Schema
		toAppendToInstances_json.Name = instance_filepath

		//append to instance_json
		instance_json = append(instance_json, toAppendToInstances_json)

		//write modified instance_json to instance.json file
		instanceJsonMarshaledByteSlice, err := json.Marshal(instance_json)
		os.WriteFile(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json"), instanceJsonMarshaledByteSlice, 0644)

		//confirm message and return 0
		fmt.Printf("Succesfully created an instance named \"%s\", at %s with schema %s\n", storeSeedFileJson.Name, instance_filepath, storeSeedFileJson.Schema)
		return 0
	} else {
		fmt.Printf("Error in making instance, : %s", "Seed file couldn't be found\n")
		return 1
	}
}
