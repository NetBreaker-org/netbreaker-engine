package configloaders

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kaushal937/NetBreaker/main/configLoaders/configInterfaces"
)

func MainConfig(filePath *string) (*configInterfaces.ServerConfig, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//innitialize and parse the config data for general stuff
	var mainConfig configInterfaces.ServerConfig

	err_parsing := json.Unmarshal([]byte(string(byte_slice)), &mainConfig)
	if err_parsing != nil {
		return nil, err_parsing
	}

	contentptr := &mainConfig
	fmt.Println("...Succesfully initialized main config file  ✔")
	return contentptr, nil
}
