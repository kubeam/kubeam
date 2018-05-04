package main

import (
	"io"
	"log"
	//"gopkg.in/gcfg.v1"
)

var (
	LogTrace   *log.Logger
	LogDebug   *log.Logger
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
)

func InitLogger(
	traceHandle io.Writer,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	LogTrace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogDebug = log.New(debugHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogInfo = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogWarning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogError = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
