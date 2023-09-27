package log

import (
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/mefit/mefit-server/utils/config"
	// "gitlab.com/mefit/mefit-server/utils/log/internal/slack"
)

var logger *logrus.Logger

func Logger() *logrus.Logger {
	return logger
}

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	debugMode := config.Config().GetBool("DEBUG")
	logger.SetLevel(logrus.DebugLevel)
	if !debugMode {
		logger.SetLevel(logrus.InfoLevel)
		logger.Formatter = new(joonix.FluentdFormatter)
	}
	logger.Infof("Debug mode is %v", debugMode)
}
