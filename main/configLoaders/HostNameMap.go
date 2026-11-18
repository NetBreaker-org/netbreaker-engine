package configloaders

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kaushal937/NetBreaker/main/configLoaders/configInterfaces"
)

// read the raw IPMap.json, with no meta data.
// misc data is stored in stats.json
func HostnameMapping(filePath *string) (*[]configInterfaces.StatusMapping, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//innitialize and parse the config data for proxied servers
	var serverStatusConfig []configInterfaces.StatusMapping

	err_parsing := json.Unmarshal([]byte(string(byte_slice)), &serverStatusConfig)
	if err_parsing != nil {
		return nil, err_parsing
	}

	contentptr := &serverStatusConfig
	fmt.Println("...Succesfully initialized Hostname-based routing mapping✔")
	return contentptr, nil
}
