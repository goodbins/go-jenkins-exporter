package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SetupLogging setup the logging properties
func SetupLogging() {
	logrus.SetOutput(os.Stdout)
	if Global.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}
