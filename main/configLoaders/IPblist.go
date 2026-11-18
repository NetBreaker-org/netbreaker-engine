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
		return 1
	} else {
		blistFile, err := os.Open(ConfigObjectIndices.IPBlackListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			return 1
		}
		defer blistFile.Close()

		blistDir := filepath.Dir(ConfigObjectIndices.IPBlackListFLPTH)
		tmpfile, err := os.CreateTemp(blistDir, "temp-*.txt")
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
			return 1
		}
		tempName := tmpfile.Name()
		defer os.Remove(tempName)

		tmpScanner := bufio.NewScanner(blistFile)
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
		blistFile.Close() //required because windows, our god father doesn't allow replacing files when is open by some other software. fun fact : my dih is gonna port this software for windows

		os.Rename(tempName, blistDir)
		return 1
	}
}
