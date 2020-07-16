package logger

import (
	"log"

	"github.com/dryairship/IITKBucks/config"
)

type LogMessageType int

const (
	RareError LogMessageType = iota
	CommonError
	MajorEvent
	MinorEvent
)

var loggingRareErrors bool
var loggingCommonErrors bool
var loggingMajorEvents bool
var loggingMinorEvents bool

func init() {
	level := config.LOG_LEVEL
	if level > 0 {
		loggingRareErrors = true
	}
	if level > 1 {
		loggingMajorEvents = true
	}
	if level > 2 {
		loggingCommonErrors = true
	}
	if level > 3 {
		loggingMinorEvents = true
	}
}

func shouldLog(messageType LogMessageType) bool {
	switch messageType {
	case RareError:
		if !loggingRareErrors {
			return false
		}
	case CommonError:
		if !loggingCommonErrors {
			return false
		}
	case MajorEvent:
		if !loggingMajorEvents {
			return false
		}
	case MinorEvent:
		if !loggingMinorEvents {
			return false
		}
	}
	return true
}

func Println(messageType LogMessageType, v ...interface{}) {
	if shouldLog(messageType) {
		log.Println(v...)
	}
}

func Printf(messageType LogMessageType, format string, v ...interface{}) {
	if shouldLog(messageType) {
		log.Printf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}
