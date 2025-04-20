package logger

import (
	"github.com/sirupsen/logrus"
	"golang-test1/pkg/config"
	"os"
)

type Logger struct {
	*logrus.Logger
}

var logger = &Logger{}

// SetUpLogger settings
func SetUpLogger() {
	logger = &Logger{logrus.New()}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	if config.AppCfg().Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
}

func GetLogger() *Logger {
	return logger
}
