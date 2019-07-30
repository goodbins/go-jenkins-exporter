package config

import (
	"time"
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
}

// Global The Global variable instance
var Global Config

// CurrentVersion version of the software
const CurrentVersion string = "v0.1.5"
