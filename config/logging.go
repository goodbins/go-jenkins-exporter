package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LogLevels Map of the logrus logging levels
var LogLevels map[string]logrus.Level

func init() {
	// Prepare log levels map
	LogLevels = make(map[string]logrus.Level)
	LogLevels["info"] = logrus.InfoLevel
	LogLevels["debug"] = logrus.DebugLevel
	LogLevels["warn"] = logrus.WarnLevel
	LogLevels["error"] = logrus.ErrorLevel
	LogLevels["panic"] = logrus.PanicLevel
	LogLevels["fatal"] = logrus.FatalLevel
}

// SetupLogging setup the logging properties
func SetupLogging() {
	logrus.SetOutput(os.Stdout) // FIXME: see cmd package
	logrus.SetLevel(LogLevels[Global.LogLevel])
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}
