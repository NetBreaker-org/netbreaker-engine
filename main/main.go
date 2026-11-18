package main

import (
	"fmt"
	"net/http"

	configloaders "github.com/kaushal937/NetBreaker/main/configLoaders"
	"github.com/kaushal937/NetBreaker/main/generals"
	"github.com/kaushal937/NetBreaker/main/logger"
)

func main() {
	//Load config file main_config.json
	mainConfigPath := generals.EnvPath.SystemPath + "/etc/netbreaker/main_config.json"
	MainConfig, err := configloaders.MainConfig(&mainConfigPath)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Load config file IP-Black List
	IPblistPath := generals.EnvPath.SystemPath + "/etc/netbreaker/ipBlackList.txt"
	contentBList, err := configloaders.LoadIPBlacklist(&IPblistPath)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Load config file IP-White List
	IPwlistPath := generals.EnvPath.SystemPath + "/etc/netbreaker/ipWhiteList.txt"
	contentWList, err := configloaders.LoadIPWhitelist(&IPwlistPath)
	if err != nil {
		logger.WriteLog(&logger.ErrorTypes.Error, &err)
		fmt.Println(err)
		return
	}

	//Index config objects
	configloaders.ConfigObjectIndices.ServerConfig = MainConfig
	configloaders.ConfigObjectIndices.ServerConfigFLPTH = mainConfigPath

	configloaders.ConfigObjectIndices.IPBlackList = contentBList
	configloaders.ConfigObjectIndices.IPBlackListFLPTH = IPblistPath

	configloaders.ConfigObjectIndices.IPWhiteList = contentWList
	configloaders.ConfigObjectIndices.IPWhiteListFLPTH = IPwlistPath
	fmt.Println(*configloaders.ConfigObjectIndices.IPWhiteList)

	//Start Service
	serverAddr := ":" + MainConfig.Port

	configloaders.DeleteIPWhiteList([]string{"33.245.235.53"})

	server := http.NewServeMux()
	fmt.Printf("Started server on %s...\n", serverAddr)
	err_listen := http.ListenAndServe(serverAddr, server)
	if err_listen != nil {
		fmt.Printf("#Err_Loc_2 Server failed to start: %v", err_listen)
	}
}
