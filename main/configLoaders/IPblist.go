package configloaders

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kaushal937/NetBreaker/main/logger"
)

func LoadIPBlacklist(filePath *string) (*[]string, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//innitialize and parse IP Blacklist
	var blackList []string

	blackList = strings.Split(string(byte_slice), "\n")
	if fmt.Sprintf("%T", blackList) != "[]string" {
		return nil, errors.New("error_in_reading_blacklist")
	}

	contentptr := &blackList
	fmt.Println("...Succesfully initialized IP blacklist  ✔")
	return contentptr, nil
}

// To append to IP Blacklist config object and to the file simultaneously
// to be called with go prefix, for concurrency
func AppendIPBlackList(ip_a []string) int {
	if ConfigObjectIndices.IPBlackList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in AppendIPBlackList call")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return 1
	} else {
		blistFile, err := os.OpenFile(ConfigObjectIndices.IPBlackListFLPTH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		for i := range ip_a {
			*ConfigObjectIndices.IPBlackList = append(*ConfigObjectIndices.IPBlackList, ip_a[i])
			blistFile.WriteString(ip_a[i] + "\n")
		}
		blistFile.Close()
		return 0
	}
}

// To delete entry in IP Blacklist config object and to the file simultaneously
// to be called with go prefix, for concurrency
func DeleteIPBlackList(ip_a []string) int {
	if ConfigObjectIndices.IPBlackList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in DeleteIPBlackList call")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return 0
	} else {
		blistFile, err := os.OpenFile(ConfigObjectIndices.IPBlackListFLPTH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0
		}
		//write new logic here
		blistFile.Close()
		return 1
	}
}
