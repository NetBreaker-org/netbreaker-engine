package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/kaushal937/NetBreaker/main/generals"
)

func resolveLogType(typeLog *string) string {
	switch *typeLog {
	case "error":
		return generals.EnvPath.SystemPath + "/var/log/netbreaker/error.log"
	case "reqLog":
		return generals.EnvPath.SystemPath + "/var/log/netbreaker/reqs.log"
	case "actions":
		return generals.EnvPath.SystemPath + "/var/log/netbreaker/actions.log"
	default:
		return generals.EnvPath.SystemPath + "/var/log/netbreaker/error.log"
	}
}

func WriteLog(typeOfLog *string, err *error) {
	logFile, err_2 := os.OpenFile(resolveLogType(typeOfLog), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err_2 != nil {
		return
	}

	_, writeTimeError := fmt.Fprintf(logFile, "%v  |  ", time.Now())
	if writeTimeError != nil {
		return
	}
	_, writeError := fmt.Fprintln(logFile, *err)
	if writeError != nil {
		return
	}
	logFile.WriteString("\n")

	logFile.Close()
}
