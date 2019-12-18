package log

import (
	"2019_2_Next_Level/pkg/logger"
	"fmt"
)

var log logger.ILog

func SetLogger(logger logger.ILog) {
	log = logger
}

func Log() logger.ILog {
	if log == nil {
		log = &logger.Log{}
	}
	return log
}

func GetLogString(login string, params ...interface{}) string {
	return fmt.Sprintf("Login: %s, status: ", login)+fmt.Sprintln(params...)
}
