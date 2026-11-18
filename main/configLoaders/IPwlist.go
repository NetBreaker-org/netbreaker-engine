package configloaders

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

// To delete entry in IP Whitelist config object and to the file simultaneously
// to be called with go prefix, for concurrency
func DeleteIPWhiteList(ip_a []string) int {
	if ConfigObjectIndices.IPWhiteList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in DeleteIPWhiteList call")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return 1
	} else {
		wlistFile, err := os.Open(ConfigObjectIndices.IPWhiteListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			return 1
		}
		defer wlistFile.Close()

		wlistDir := filepath.Dir(ConfigObjectIndices.IPBlackListFLPTH)
		tmpfile, err := os.CreateTemp(wlistDir, "temp-*.txt")
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			return 1
		}
		tempName := tmpfile.Name()
		defer os.Remove(tempName)

		tmpScanner := bufio.NewScanner(wlistFile)
		tmpWriter := bufio.NewWriter(tmpfile)

		for tmpScanner.Scan() {
			line := tmpScanner.Text()

			for i := range ip_a {
				if !strings.Contains(line, ip_a[i]) {
					_, err := fmt.Fprintln(tmpWriter, line)
					if err != nil {
						deleteError := fmt.Errorf("Couldn't  delete IP addresses provided. wrote till : %dth index (%s). Closed file.", i, ip_a[i])
						logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
						tmpfile.Close()
						return 1
					}
				}
			}
		}

		if tmpScanner.Err() != nil {
			tmpfile.Close()
			return 1
		}
		if tmpWriter.Flush() != nil {
			tmpfile.Close()
			return 1
		}
		if tmpfile.Close() != nil {
			return 1
		}
		wlistFile.Close() //required because windows, our god father doesn't allow replacing files when is open by some other software. fun fact : my dih is gonna port this software for windows

		os.Rename(tempName, wlistDir)
		return 1
	}
}
