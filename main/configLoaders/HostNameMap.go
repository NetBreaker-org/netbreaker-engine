package configloaders

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kaushal937/NetBreaker/main/configLoaders/configInterfaces"
	"github.com/kaushal937/NetBreaker/main/logger"
)

//feature scope : create backupfile before running hostNameMapFilesBuilder()

// read the raw hostnameMapping.json, with no meta data.
// misc data is stored in stats.json
func HostnameMapping(filePath *string) (*[]configInterfaces.StatusMapping, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading Hostname Mapping. Attempting to load again", *filePath)
			logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)

			nestedErr := hostNameMapFilesBuilder(ConfigObjectIndices.IPMapFLPTH)
			if nestedErr != nil {
				var sending_errNested error = fmt.Errorf("Error : Couldn't read file : %s while loading Hostname Mapping. Not attempting again (nestedErr)", *filePath)
				logger.WriteLog(&logger.ErrorTypes.Error, &sending_errNested)
				return nil, err
			}
			return HostnameMapping(filePath)
		}
		var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading Hostname Mapping. Not attempting again", *filePath)
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return nil, err
	}

	//initialize and parse the config data for proxied servers
	var serverStatusConfig []configInterfaces.StatusMapping

	err_parsing := json.Unmarshal([]byte(string(byte_slice)), &serverStatusConfig)
	if err_parsing != nil {
		return nil, err_parsing
	}

	contentptr := &serverStatusConfig
	fmt.Println("...Succesfully initialized Hostname-based routing mapping  ✔")
	return contentptr, nil
}

// read the list of available hostnames
func AvailableHostnameReader(filePath *string) (*[]string, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading available hostnames list. Attempting to load again", *filePath)
			logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)

			nestedErr := hostNameMapFilesBuilder(ConfigObjectIndices.IPMapFLPTH)
			if nestedErr != nil {
				var sending_errNested error = fmt.Errorf("Error : Couldn't read file : %s while loading available hostnames list. Not attempting again", *filePath)
				logger.WriteLog(&logger.ErrorTypes.Error, &sending_errNested)
				return nil, err
			}
			return AvailableHostnameReader(filePath)
		}
		var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading available hostnames list. Not attempting again", *filePath)
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return nil, err
	}

	//innitialize and parse the config data for proxied servers
	var availableHostnames []string = strings.Split(string(byte_slice), "\n")

	//create pointer and return
	contentptr := &availableHostnames
	fmt.Println("...Succesfully initialized list of available hostnames  ✔")
	return contentptr, nil
}

func hostNameMapFilesBuilder(filePath string) error {
	//Read the IPMap.json file
	fmt.Println(filePath)
	buffFileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	//create the variable with type of map[string][]string to load the unmarshaled data
	var rawHostnameMapping map[string][]string //fix type mismatch
	//create an array to store the available hostnames, to route
	var availableHostnames []string
	//create an array with type []StatusMapping to temporarily store that data and then write to the hostname mapping file
	var hostnameMapping []configInterfaces.StatusMapping

	//parse the buffered data into json with the above type interface
	err = json.Unmarshal([]byte(string(buffFileData)), &rawHostnameMapping)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	//loop over the map : rawHostnameMapping and prepare data appropriate of []StatusMapping to stream/write to hostname mapping file
	for key, value := range rawHostnameMapping {
		//append all hostnmaes available to availableHostnames array, to be used by proxy engines
		availableHostnames = append(availableHostnames, key)

		//create a variable of []configInterfaces.HostnameServerArrayData to store the data of type []configInterfaces.HostnameServerArrayData to add to valueToAppend -> which will be appended to hostnameMapping
		var hostnameServerArrayDataTemplate []configInterfaces.HostnameServerArrayData
		//loop over value to append to hostnameServerArrayDataTemplate for every instance of the proxied server
		for i := range value {
			hostnameServerArrayDataTemplate = append(hostnameServerArrayDataTemplate, configInterfaces.HostnameServerArrayData{Address: value[i], Status: 1, LoadAssignment: []string{}})
		}

		//prepare valueToAppend of type configInterfaces.StatusMapping, to be appended to hostnameMapping
		var valueToAppend configInterfaces.StatusMapping = configInterfaces.StatusMapping{Hostname: key, StatusData: hostnameServerArrayDataTemplate, MiscellaneousData: configInterfaces.MiscellaneousData{Views: 0, Reqs: 0}}
		hostnameMapping = append(hostnameMapping, valueToAppend)
	} //hostnameMapping is now ready to be streamed

	//--write hostnameMapping to hostnameMappingFile & availableHostnames to availableHostnamesListFile--

	//create hostnameMappingFile & availableHostnamesListFile
	hostnameMappingFile, err := os.Create(ConfigObjectIndices.HostnameMappingFLPTH)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}
	defer hostnameMappingFile.Close()

	availableHostnamesListFile, err := os.Create(ConfigObjectIndices.IPMapFLPTH)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}
	defer hostnameMappingFile.Close()

	//create write streams to the files
	hostnameMappingFileWriter := bufio.NewWriter(hostnameMappingFile)
	availableHostnamesListFileWriter := bufio.NewWriter(availableHostnamesListFile)

	//stream hostnameMapping to file
	_, err_2 := fmt.Fprintln(hostnameMappingFileWriter, hostnameMapping)
	if err_2 != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err_2)
		return err
	}

	_, err_3 := fmt.Fprintln(availableHostnamesListFileWriter, availableHostnames)
	if err_3 != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err_3)
		return err
	}
	return nil
}
