package log

import "2019_2_Next_Level/pkg/logger"

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
