package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	configloaders "github.com/NetBreaker-org/netbreaker-engine/main/configLoaders"
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
	configloaders.ConfigObjectIndices.ServerConfigFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/instances.json")
	InstancesFile, err := configloaders.InstancesConfig(&configloaders.ConfigObjectIndices.ServerConfigFLPTH)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Load config file IP-Black List
	// configloaders.ConfigObjectIndices.IPBlackListFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/ipBlackList.txt")
	// contentBList, err := configloaders.LoadIPBlacklist(&configloaders.ConfigObjectIndices.IPBlackListFLPTH)
	// if err != nil {
	// 	logger.WriteLog(&logger.ErrorTypes.Error, &err)
	// 	fmt.Println(err)
	// 	return
	// }

	//Load config file IP-White List
	// configloaders.ConfigObjectIndices.IPWhiteListFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/ipWhiteList.txt")
	// contentWList, err := configloaders.LoadIPWhitelist(&configloaders.ConfigObjectIndices.IPWhiteListFLPTH)
	// if err != nil {
	// 	logger.WriteLog(&logger.ErrorTypes.Error, &err)
	// 	fmt.Println(err)
	// 	return
	// }

	//Load IPMap
	// configloaders.ConfigObjectIndices.AvailableHostnamesFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/parsed/availableHostnamesListFile.txt")
	// configloaders.ConfigObjectIndices.IPMapFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/IPMap.json")
	// contentAvailableHostnameArray, err := configloaders.AvailableHostnameReader(&configloaders.ConfigObjectIndices.IPMapFLPTH)
	// if err != nil {
	// 	logger.WriteLog(&logger.ErrorTypes.Error, &err)
	// 	fmt.Println(err)
	// 	return
	// }

	//Load hostnameMapping
	// configloaders.ConfigObjectIndices.HostnameMappingFLPTH = filepath.Join(generals.EnvPath.SystemPath, "/etc/netbreaker/parsed/hostnameMapping.json")
	// contentHostnameMapping, err := configloaders.HostnameMapping(&configloaders.ConfigObjectIndices.HostnameMappingFLPTH)
	// if err != nil {
	// 	logger.WriteLog(&logger.ErrorTypes.Error, &err)
	// 	fmt.Println(err)
	// 	return
	// }

	//Index config objects
	configloaders.ConfigObjectIndices.ServerConfig = MainConfig
	configloaders.ConfigObjectIndices.InstanceList = InstancesFile
	// configloaders.ConfigObjectIndices.IPBlackList = contentBList
	// configloaders.ConfigObjectIndices.IPWhiteList = contentWList
	// configloaders.ConfigObjectIndices.IPMap = contentAvailableHostnameArray
	// configloaders.ConfigObjectIndices.HostnameMapping = contentHostnameMapping

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
