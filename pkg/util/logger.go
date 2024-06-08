package util

import (
	"fmt"
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type Logger struct{}

func (l *Logger) Init(logFileName string) *Logger {

	cfg := prettyconsole.NewConfig()

	cfg.OutputPaths = []string{
		"stdout",
	}

	disableLogFile, parsedBoolError := strconv.ParseBool(os.Getenv("DISABLE_LOG_FILE"))

	if parsedBoolError != nil {
		disableLogFile = false
	}

	if !disableLogFile {

		logFile := os.Getenv("LOG_FILE")

		if logFile == "" {
			_ = os.Mkdir("./logs", os.ModePerm)
			logFile = fmt.Sprintf("./logs/%s.log", logFileName)
		}

		//fmt.Println(fmt.Sprintf("logfile set to: %s", logFile))

		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}

	logger := zap.Must(cfg.Build())

	zap.ReplaceGlobals(logger)

	return l
}
