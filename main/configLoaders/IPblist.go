package configloaders

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

// stack handler for DeleteIPBlackList function
// FIFO indexing
var DeleteIPBlackList_FIFO [][]string

// Using this lease variable so that this function can't be called more than once at a time, reducing blacklist file / config object corrruption
// 1 when a goroutine of DeleteIPBlackList, 0 elsewise
var DeleteIPBlackListLease int = 0

func DeleteIPBlackList(ip_a []string) {
	DeleteIPBlackList_FIFO = append(DeleteIPBlackList_FIFO, ip_a)
	clear_DeleteIPBlackList()
}

func clear_DeleteIPBlackList() {
	if DeleteIPBlackListLease == 0 {
		go DeleteIPBlackList_(DeleteIPBlackList_FIFO[0])
	} else {
		return
	}
}

// To delete entry in IP Blacklist config object and to the file simultaneously
// to be called with go prefix, for concurrency
func DeleteIPBlackList_(ip_a []string) {
	//grab lease variable
	DeleteIPBlackListLease = 1

	if ConfigObjectIndices.IPBlackList == nil {
		var sending_err error = errors.New("Error : Received nil pointer in DeleteIPBlackList call")
		logger.WriteLog(&logger.ErrorTypes.Error, &sending_err)
	} else {
		blistFile, err := os.Open(ConfigObjectIndices.IPBlackListFLPTH)
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
		}
		defer blistFile.Close()

		blistDir := filepath.Dir(ConfigObjectIndices.IPBlackListFLPTH)
		tmpfile, err := os.CreateTemp(blistDir, "temp-*.txt")
		if err != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err)
		}
		tempName := tmpfile.Name()
		defer os.Remove(tempName)

		tmpScanner := bufio.NewScanner(blistFile)
		tmpWriter := bufio.NewWriter(tmpfile)

		for tmpScanner.Scan() {
			line := tmpScanner.Text()

			for i := range ip_a {
				//Big O notation : O(len(ip_a))
				if (*ConfigObjectIndices.IPBlackList)[i] == ip_a[i] {
					(*ConfigObjectIndices.IPBlackList)[i] = (*ConfigObjectIndices.IPBlackList)[len(*ConfigObjectIndices.IPBlackList)-1]
					(*ConfigObjectIndices.IPBlackList) = (*ConfigObjectIndices.IPBlackList)[:len(*ConfigObjectIndices.IPBlackList)-1]

				}
				//Big O notation : O(len(ip_a)^2)
				//Big Ω*omega_symbol* notation : Ω*omega_symbol*(len(ip_a))
				if !strings.Contains(line, ip_a[i]) {
					_, err := fmt.Fprintln(tmpWriter, line)
					if err != nil {
						deleteError := fmt.Errorf("Couldn't  delete IP addresses provided. Wrote till : %dth index (%s). Closed file.", i, ip_a[i])
						logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
						tmpfile.Close()
					}
				}
			}
		}

		if tmpScanner.Err() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			blistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			tmpfile.Close()
			os.Remove(tempName)
			return
		}
		if tmpWriter.Flush() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			blistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			tmpfile.Close()
			os.Remove(tempName)
			return
		}
		if tmpfile.Close() != nil {
			deleteError := fmt.Errorf("Couldn't delete IP addresses provided (whitelist)")
			logger.WriteLog(&logger.ErrorTypes.Error, &deleteError)
			blistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
			os.Remove(tempName)
			return
		}

		blistFile.Close() //required because windows, our god father doesn't allow replacing files when it is open by some other software. fun fact : my dih is gonna port this software for windows. also you can't just build yourself for windows, this NetBreaker uses Unix syscall and UDS
		err_2 := os.Rename(tempName, blistDir)
		if err_2 != nil {
			logger.WriteLog(&logger.ErrorTypes.Error, &err_2)
			os.Remove(tempName)
		}

		//release lease variable & clear FIFO list's 0 indexed value
		DeleteIPBlackListLease = 0

		DeleteIPBlackList_FIFO = slices.Delete(DeleteIPBlackList_FIFO, 0, 1) //there might be a better way to do this
	}
}
