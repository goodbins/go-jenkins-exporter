package config

import (
	"time"
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

// CurrentVersion version of the software
const CurrentVersion string = "v0.1.3"
