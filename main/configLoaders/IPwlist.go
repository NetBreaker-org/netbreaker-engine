package configloaders

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/NetBreaker-org/netbreaker-engine/main/logger"
)

func LoadIPWhitelist(filePath *string) (*[]string, error) {
	//get buffer data of that file
	byte_slice, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	//innitialize and parse IP Whitelist
	var whiteList []string

	whiteList = strings.Split(string(byte_slice), "\n")
	if fmt.Sprintf("%T", whiteList) != "[]string" {
		return nil, errors.New("error_in_reading_whitelist")
	}

	contentptr := &whiteList
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

// stack handler for DeleteIPWhiteList function
// FIFO indexing
var deleteIPWhiteList_FIFO [][]string

// Using this lease variable so that this function can't be called more than once at a time, reducing whitelist file / config object corrruption
// 1 when a goroutine of DeleteIPWhiteList, 0 elsewise
var deleteIPWhiteListLease int = 0

func DeleteIPWhiteList(ip_a []string) {
	deleteIPWhiteList_FIFO = append(deleteIPWhiteList_FIFO, ip_a)
	clear_DeleteIPWhiteList()
}

func clear_DeleteIPWhiteList() {
	if deleteIPWhiteListLease == 0 {
		go deleteIPWhiteList_(deleteIPWhiteList_FIFO[0])
	} else {
		var sending_err error = errors.New("Error : Couldn't start DeleteIPWhiteList routine")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
		return
	}
}

// To delete entry in IP Whitelist config object and to the file simultaneously
// to be called with go prefix, for concurrency
func deleteIPWhiteList_(ip_a []string) {
	//grab lease variable
	deleteIPWhiteListLease = 1

	if ConfigObjectIndices.IPWhiteList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in DeleteIPWhiteList call")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
	} else {
		wlistFile, err := os.Open(ConfigObjectIndices.IPWhiteListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
		}
		defer wlistFile.Close()

		wlistDir := filepath.Dir(ConfigObjectIndices.IPWhiteListFLPTH)
		tmpfile, err := os.CreateTemp(wlistDir, "temp-*.txt")
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
		}
		tempName := tmpfile.Name()
		defer os.Remove(tempName)

		tmpScanner := bufio.NewScanner(wlistFile)
		tmpWriter := bufio.NewWriter(tmpfile)

		for tmpScanner.Scan() {
			line := tmpScanner.Text()

			for i := range ip_a {
				//Big O notation : O(len(ip_a))
				if (*ConfigObjectIndices.IPWhiteList)[i] == ip_a[i] {
					(*ConfigObjectIndices.IPWhiteList)[i] = (*ConfigObjectIndices.IPWhiteList)[len(*ConfigObjectIndices.IPWhiteList)-1]
					(*ConfigObjectIndices.IPWhiteList) = (*ConfigObjectIndices.IPWhiteList)[:len(*ConfigObjectIndices.IPWhiteList)-1]

				}
				//Big O notation : O(len(ip_a)^2)
				//Big Ω*omega_symbol* notation : Ω*omega_symbol*(len(ip_a))
				if line != ip_a[i] {
					_, err := fmt.Fprintln(tmpWriter, line)
					if err != nil {
						deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist). Wrote till : %dth index (%s). Closed file.", i, ip_a[i])
						logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
						tmpfile.Close()
					}
				}
			}
		}

		if tmpScanner.Err() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			wlistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			tmpfile.Close()
			os.Remove(tempName)
			return
		}
		if tmpWriter.Flush() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			wlistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			tmpfile.Close()
			os.Remove(tempName)
			return
		}
		if tmpfile.Close() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			wlistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			os.Remove(tempName)
			return
		}

		wlistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
		err_2 := os.Rename(tempName, ConfigObjectIndices.IPWhiteListFLPTH)
		if err_2 != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err_2)
			os.Remove(tempName)
		}

		//release lease variable & clear FIFO list's 0 indexed value
		deleteIPWhiteListLease = 0

		deleteIPWhiteList_FIFO = slices.Delete(deleteIPWhiteList_FIFO, 0, 1) //there might be a better way to do this
	}
}
