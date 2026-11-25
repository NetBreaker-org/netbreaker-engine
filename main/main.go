package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	configloaders "github.com/NetBreaker-org/netbreaker-engine/main/configLoaders"
	"github.com/NetBreaker-org/netbreaker-engine/main/configLoaders/configInterfaces"
	"github.com/NetBreaker-org/netbreaker-engine/main/generals"
	"github.com/NetBreaker-org/netbreaker-engine/main/logger"
)

func main() {
	//Load config file main_config.json
	configloaders.ConfigObjectIndices.ServerConfigFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/master_config.json")
	MainConfig, err := configloaders.MainConfig(&configloaders.ConfigObjectIndices.ServerConfigFLPTH)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Load instances file : instances.json
	configloaders.ConfigObjectIndices.InstanceListFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json")
	InstancesFile, err := configloaders.InstancesConfig(&configloaders.ConfigObjectIndices.InstanceListFLPTH)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Index config objects
	configloaders.ConfigObjectIndices.ServerConfig = MainConfig
	configloaders.ConfigObjectIndices.InstanceList = InstancesFile

	//setup config for each instance
	for i := range *configloaders.ConfigObjectIndices.InstanceList {
		var instanceConfig configInterfaces.InstanceConfigObjectPointers = configInterfaces.InstanceConfigObjectPointers{}

		//Load config file IP-Black List
		instanceConfig.InstanceIPBlackListFLPTH = filepath.Join((*configloaders.ConfigObjectIndices.InstanceList)[i].Dir, "ipBlackList.txt")
		contentBList, err := configloaders.LoadIPBlacklist(&instanceConfig.InstanceIPBlackListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			fmt.Println(err)
			return
		}

		instanceConfig.InstanceIPBlackList = contentBList

		//Load config file IP-White List
		instanceConfig.InstanceIPWhiteListFLPTH = filepath.Join((*configloaders.ConfigObjectIndices.InstanceList)[i].Dir, "ipWhiteList.txt")
		contentWList, err := configloaders.LoadIPWhitelist(&instanceConfig.InstanceIPWhiteListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			fmt.Println(err)
			return
		}

		instanceConfig.InstanceIPWhiteList = contentWList

		//Load IPMap
		instanceConfig.InstanceAvailableHostnamesFLPTH = filepath.Join((*configloaders.ConfigObjectIndices.InstanceList)[i].Dir, "availableHostnamesListFile.txt")
		instanceConfig.InstanceIPMapFLPTH = filepath.Join((*configloaders.ConfigObjectIndices.InstanceList)[i].Dir, "IPMap.json")
		contentAvailableHostnameArray, err := configloaders.AvailableHostnameReader(&instanceConfig.InstanceIPMapFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			fmt.Println(err)
			return
		}

		instanceConfig.InstanceIPMap = contentAvailableHostnameArray

		// Load hostnameMapping
		instanceConfig.InstanceHostnameMappingFLPTH = filepath.Join((*configloaders.ConfigObjectIndices.InstanceList)[i].Dir, "hostnameMapping.json")
		contentHostnameMapping, err := configloaders.HostnameMapping(&instanceConfig.InstanceHostnameMappingFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			fmt.Println(err)
			return
		}

		instanceConfig.InstanceHostnameMapping = contentHostnameMapping

		configloaders.InstanceConfigObjectIndices = append(configloaders.InstanceConfigObjectIndices, configInterfaces.InstanceConfigs{Metadata: (*configloaders.ConfigObjectIndices.InstanceList)[i], Config: instanceConfig})
	}

	//Start Service
	serverAddr := ":" + MainConfig.Port

	configloaders.DeleteIPBlackList([]string{"234.435.2.5"})

	server := http.NewServeMux()
	fmt.Printf("Started server on %s...\n", serverAddr)
	err_listen := http.ListenAndServe(serverAddr, server)
	if err_listen != nil {
		fmt.Printf("#Err_Loc_2 Server failed to start: %v", err_listen)
	}
}
