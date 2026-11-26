package configloaders

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NetBreaker-org/netbreaker-engine/main/configLoaders/configInterfaces"
	"github.com/NetBreaker-org/netbreaker-engine/main/logger"
)

//feature scope : create backupfile before running hostNameMapFilesBuilder()

// read the raw hostnameMapping.json, with no meta data.
// misc data is stored in stats.json
func HostnameMapping(filePath *string, indexInInstanceConfig int) (*[]configInterfaces.StatusMapping, error) {
	// if len(InstanceConfigObjectIndices) == 0 {
	// 	fmt.Println("stupid retard goofy motherfucker")
	// 	return nil, errors.New("stupid retard goofy motherfucker")
	// }

	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading Hostname Mapping. Attempting to load again", *filePath)
			logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)

			nestedErr := hostNameMapFilesBuilder(filepath.Join((*ConfigObjectIndices.InstanceList)[indexInInstanceConfig].Dir, "IPMap.json"), indexInInstanceConfig)
			if nestedErr != nil {
				var sending_errNested error = fmt.Errorf("Error : Couldn't read file : %s while loading Hostname Mapping. Not attempting again (nestedErr)", *filePath)
				logger.WriteLog(&logger.ErrorTypes.Error, &sending_errNested)
				return nil, err
			}
			return HostnameMapping(filePath, indexInInstanceConfig)
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
func AvailableHostnameReader(filePath *string, indexInInstanceConfig int) (*[]string, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var sending_err error = fmt.Errorf("Error : Couldn't read file : %s while loading available hostnames list. Attempting to load again", *filePath)
			logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)

			nestedErr := hostNameMapFilesBuilder(InstanceConfigObjectIndices[indexInInstanceConfig].Config.InstanceIPMapFLPTH, indexInInstanceConfig)
			if nestedErr != nil {
				var sending_errNested error = fmt.Errorf("Error : Couldn't read file : %s while loading available hostnames list. Not attempting again", *filePath)
				logger.WriteLog(&logger.ErrorTypes.Error, &sending_errNested)
				return nil, err
			}
			return AvailableHostnameReader(filePath, indexInInstanceConfig)
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

func hostNameMapFilesBuilder(filePath string, indexInInstanceConfig int) error {
	//Read the IPMap.json file
	fmt.Println(filePath)
	buffFileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	//create the variable with type of map[string][]string to load the unmarshaled data
	type fileData struct {
		Hostname string   `json:"hostname"`
		Pointers []string `json:"pointers"`
	}
	var loadedFileData []fileData
	//create an array to store the available hostnames, to route
	var availableHostnames []string
	//create an array with type []StatusMapping to temporarily store that data and then write to the hostname mapping file
	var hostnameMapping []configInterfaces.StatusMapping

	//parse the buffered data into json with the above type interface
	err = json.Unmarshal(buffFileData, &loadedFileData)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	//transfer loadedFileData to rawHostnameMapping map
	//loop over the loadedFileData :  and prepare data appropriate of []StatusMapping to stream/write to hostname mapping file
	for i := range loadedFileData {
		//append all hostnmaes available to availableHostnames array, to be used by proxy engines
		availableHostnames = append(availableHostnames, loadedFileData[i].Hostname)

		//create a variable of []configInterfaces.HostnameServerArrayData to store the data of type []configInterfaces.HostnameServerArrayData to add to valueToAppend -> which will be appended to hostnameMapping
		var hostnameServerArrayDataTemplate []configInterfaces.HostnameServerArrayData
		//loop over value to append to hostnameServerArrayDataTemplate for every instance of the proxied server
		for k := range loadedFileData[i].Pointers {
			hostnameServerArrayDataTemplate = append(hostnameServerArrayDataTemplate, configInterfaces.HostnameServerArrayData{Address: loadedFileData[i].Pointers[k], Status: 1, LoadAssignment: []string{}})
		}

		//prepare valueToAppend of type configInterfaces.StatusMapping, to be appended to hostnameMapping
		var valueToAppend configInterfaces.StatusMapping = configInterfaces.StatusMapping{Hostname: loadedFileData[i].Hostname, StatusData: hostnameServerArrayDataTemplate, MiscellaneousData: configInterfaces.MiscellaneousData{Views: 0, Reqs: 0}}
		hostnameMapping = append(hostnameMapping, valueToAppend)
	} //hostnameMapping is now ready to be streamed

	//--write hostnameMapping to hostnameMappingFile & availableHostnames to availableHostnamesListFile--

	//create hostnameMappingFile & availableHostnamesListFile
	hostnameMappingFile, err := os.Create(filepath.Join((*ConfigObjectIndices.InstanceList)[indexInInstanceConfig].Dir, "parsed", "hostnameMapping.json"))
	fmt.Println(filepath.Join((*ConfigObjectIndices.InstanceList)[indexInInstanceConfig].Dir, "parsed", "hostnameMapping.json"))
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}
	defer hostnameMappingFile.Close()

	availableHostnamesListFile, err := os.Create(filepath.Join((*ConfigObjectIndices.InstanceList)[indexInInstanceConfig].Dir, "parsed", "availableHostnamesListFile.txt"))
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}
	defer availableHostnamesListFile.Close()

	//create write streams to the files
	hostnameMappingFileWriter := bufio.NewWriter(hostnameMappingFile)
	availableHostnamesListFileWriter := bufio.NewWriter(availableHostnamesListFile)

	//stream hostnameMapping to file
	toWriteHostnameMapping, err := json.Marshal(hostnameMapping)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	//write to hostnameMapping file
	_, err_2 := fmt.Fprintln(hostnameMappingFileWriter, string(toWriteHostnameMapping))
	if err_2 != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err_2)
		return err
	}

	err = hostnameMappingFileWriter.Flush()
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	toWriteAvailableHostnames, err := json.Marshal(availableHostnames)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	_, err_3 := fmt.Fprintln(availableHostnamesListFileWriter, string(toWriteAvailableHostnames))
	if err_3 != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err_3)
		return err
	}

	err = availableHostnamesListFileWriter.Flush()
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		return err
	}

	return nil
}

//becoz of some shit the json isnt beint streamed to the files
