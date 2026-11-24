package configloaders

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/NetBreaker-org/netbreaker-engine/main/configLoaders/configInterfaces"
)

func MainConfig(filePath *string) (*configInterfaces.ServerConfig, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//initialize and parse the config data for general stuff
	var mainConfig configInterfaces.ServerConfig

	err_parsing := json.Unmarshal([]byte(string(byte_slice)), &mainConfig)
	if err_parsing != nil {
		return nil, err_parsing
	}

	contentptr := &mainConfig
	fmt.Println("...Succesfully initialized master config file  ✔")
	return contentptr, nil
}

// load instances.json
func InstancesConfig(filePath *string) (*[]configInterfaces.InstancesJson, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//initialize and parse the config data for general stuff
	var instanceConfig []configInterfaces.InstancesJson

	err_parsing := json.Unmarshal([]byte(string(byte_slice)), &instanceConfig)
	if err_parsing != nil {
		return nil, err_parsing
	}

	contentptr := &instanceConfig
	fmt.Println("...Succesfully initialized list of all proxy instances  ✔")
	return contentptr, nil
}
