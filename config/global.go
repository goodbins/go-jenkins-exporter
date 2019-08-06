package config

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Config Global configuration for the jenkins exporter
type Config struct {
	SSLOn              bool
	JenkinsAPIHostPort string
	JenkinsAPIPath     string
	JenkinsAPITimeout  time.Duration
	JenkinsUsername    string
	JenkinsPassword    string
	JenkinsToken       string
	JenkinsWithCreds   bool
	ExporterHostPort   string
	MetricsPath        string
	MetricsUpdateRate  time.Duration
	Verbose            bool
	LogLevel           string
}

// Global The Global variable instance
var Global Config

// CurrentVersion version of the software
const CurrentVersion string = "v0.2.1"

// Logrus log levels
var LogrusLevels = map[string]logrus.Level{
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}
