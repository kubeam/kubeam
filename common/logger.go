package common

import (
	"io"
	"log"
)

var (
	//LogTrace ...
	LogTrace *log.Logger
	//LogDebug ...
	LogDebug *log.Logger
	//LogInfo ...
	LogInfo *log.Logger
	//LogWarning ...
	LogWarning *log.Logger
	//LogError ...
	LogError *log.Logger
)

/*InitLogger initializes different logging handlers*/
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
