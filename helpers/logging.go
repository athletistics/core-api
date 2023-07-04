// Brief   : Helper functions simplifying logging events
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 04.07.2023

package helpers

import "log"

type LogLevel int

const (
	logNone LogLevel = iota
	logError
	logSecurity
	logWarning
	logInfo
	logDebug
)

func (l LogLevel) String() string {
	switch l {
	case logError:
		return "ERROR"
	case logSecurity:
		return "SECURITY"
	case logWarning:
		return "WARNING"
	case logInfo:
		return "INFO"
	case logDebug:
		return "DEBUG"
	default:
		return "NONE"
	}
}

func LogError(msg string, user string) {
	log.Printf("[%v] (user %v): %v\n", logError, user, msg)
}

func LogSecurity(msg string, user string) {
	log.Printf("[%v] (user %v): %v\n", logSecurity, user, msg)
}

func LogWarning(msg string, user string) {
	log.Printf("[%v] (user %v): %v\n", logWarning, user, msg)
}

func LogInfo(msg string, user string) {
	log.Printf("[%v] (user %v): %v\n", logInfo, user, msg)
}

func LogDebug(msg string, user string) {
	log.Printf("[%v] (user %v): %v\n", logDebug, user, msg)
}

func LogFatal(err error) {
	log.Fatalf("[%v] (system): %v\n", logError, err)
}
