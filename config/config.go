package config

import "time"

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
}

var Global Config
