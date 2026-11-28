package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/NetBreaker-org/netbreaker-engine/main/generals"
	"github.com/NetBreaker-org/netbreaker-engine/main/logger"
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
		toAppendToInstances_json.Dir = instance_filepath

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

// delete the instance folder and remove entry from instances.json
func DeleteInstance(instanceNameToDelete string) int {
	//check if entry exists in instances.json file
	var instance_json []InstancesJson
	instanceFileByteSlice, err := os.ReadFile(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json"))
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Printf("instances.json file not found. Exiting process.")
		fmt.Println(err)
		return 1
	}

	err = json.Unmarshal(instanceFileByteSlice, &instance_json)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return 1
	}

	//run for loop on the instance_json, find the entry to delete and get the array of that obj in instance_json and delete it
	indexNumberOfInstanceToDeleteInInstanceJson := 0
	for i := range instance_json {
		if instance_json[i].Name == instanceNameToDelete {
			indexNumberOfInstanceToDeleteInInstanceJson = i
			break
		} else {
			toSendErr := fmt.Errorf("No instance named \"%s\" found in instances.json. Exiting process", instanceNameToDelete)
			fmt.Printf("No instance named \"%s\" found in instances.json. Exiting process\n", instanceNameToDelete)
			logger.WriteLog(&logger.ErrorTypes.Error, &toSendErr)
			return 1
		}
	}

	//check if folder of that instance exists
	if fileExists(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances", instanceNameToDelete)) != true {
		toSendErr := fmt.Errorf("No instance named \"%s\" found in instances folder. Exiting process", instanceNameToDelete)
		fmt.Printf("No instance named \"%s\" found in instances folder. Exiting process\n", instanceNameToDelete)
		logger.WriteLog(&logger.ErrorTypes.Error, &toSendErr)
		return 1
	}

	//now since both things exist, ask for confirmation and delete instance entry and instance folder
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Are you sure you want to delete this instance?\nType \"I am deleting %s\" : ", instanceNameToDelete)
	if scanner.Scan() {
		inText := scanner.Text()
		if inText == "I am deleting "+instanceNameToDelete {
			fmt.Printf("Deleting Instance %s.\n", instanceNameToDelete)

			//delete folder
			fmt.Printf("Attempting to delete instance folder...\n")
			err = os.RemoveAll(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances", instanceNameToDelete))
			if err != nil {
				toSendErr := fmt.Errorf("Error in deleting instance folder : \"%s\". Exiting process with code 1", instanceNameToDelete)
				fmt.Printf("I'm facing a hard time deleting the instance folder...\nExiting process with code 1...\n")
				logger.WriteLog(&logger.ErrorTypes.Error, &toSendErr)
				return 1
			}
			fmt.Printf("Successfully deleted instance folder\n")

			//delete entry in instances.json
			fmt.Printf("Attempting to delete instance entry from instances.json...\n")

			instance_json = slices.Delete(instance_json, indexNumberOfInstanceToDeleteInInstanceJson, indexNumberOfInstanceToDeleteInInstanceJson+1)
			toWriteInstanceJson, err := json.Marshal(&instance_json)
			if err != nil {
				toSendErr := fmt.Errorf("Error in marshaling instances.json data (to write). Exiting process with code 1")
				fmt.Printf("Deleting instance entry from instances.json is causing me some trouble...\nExiting process with code 1...\n")
				logger.WriteLog(&logger.ErrorTypes.Error, &toSendErr)
				return 1
			}

			err = os.WriteFile(filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json"), toWriteInstanceJson, 0644)
			if err != nil {
				toSendErr := fmt.Errorf("Error in deleting instance folder : \"%s\". Exiting process with code 1", instanceNameToDelete)
				fmt.Printf("Deleting instance entry from instances.json is causing some trouble...\nExiting process with code 1...\n")
				logger.WriteLog(&logger.ErrorTypes.Error, &toSendErr)
				return 1
			}
			fmt.Printf("Successfully removed instance entry from instances.json\n")
			fmt.Printf("Successfully deleted the instance. Exiting process with code 0\n")
			return 0
		} else {
			fmt.Printf("I'm sorry... I can't do this sin...\nExiting process with code 1\n")
			return 1
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error occured in reading os.Stdin\n")
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return 1
	}
	return 0
}
