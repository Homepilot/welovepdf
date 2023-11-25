package utils

import "log"

func LogFatalAndPanic(msg string, err error) {
	log.Fatalf("%s : %s", msg, err.Error())
	panic(err)
}
