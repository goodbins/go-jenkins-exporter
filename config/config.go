package config

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Config Global configuration for the jenkins exporter
type Config struct {
	SSLOn             bool
	JenkinsAPIHost    string
	JenkinsAPIPort    int
	JenkinsAPIPath    string
	JenkinsAPITimeout time.Duration
	JenkinsUsername   string
	JenkinsPassword   string
	JenkinsToken      string
	JenkinsWithCreds  bool
	ExporterHost      string
	ExporterPort      int
	MetricsPath       string
	MetricsUpdateRate time.Duration
	LogLevel          string
}

// Global The Global variable instance
var Global Config

// LogLevels Map of the logrus logging levels
var LogLevels map[string]logrus.Level
