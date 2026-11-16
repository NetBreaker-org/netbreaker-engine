package configloaders

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kaushal937/NetBreaker/main/logger"
)

func LoadIPWhitelist(filePath *string) (*[]string, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//innitialize and parse IP Whitelist
	var blackList []string

	blackList = strings.Split(string(byte_slice), "\n")
	if fmt.Sprintf("%T", blackList) != "[]string" {
		return nil, errors.New("error_in_reading_whitelist")
	}

	contentptr := &blackList
	fmt.Println("...Succesfully initialized IP whitelist  ✔")
	return contentptr, nil
}

// to be called with go prefix, for concurrency
func AppendIPWhiteList(ip_a []string) int {
	if ConfigObjectIndices.IPWhiteList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in AppendIPWhiteList")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return 0
	} else {
		wlistFile, err := os.OpenFile(ConfigObjectIndices.IPWhiteListFLPTH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0
		}

		for i := range ip_a {
			*ConfigObjectIndices.IPWhiteList = append(*ConfigObjectIndices.IPWhiteList, ip_a[i])
			wlistFile.WriteString(ip_a[i] + "\n")
		}
		wlistFile.Close()
		return 1
	}
}
