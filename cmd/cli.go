package cmd

import (
	"time"

	"github.com/abousselmi/go-jenkins-exporter/config"
)

// Will be read from CMD package
var sslOn bool = false
var apiHost string = "127.0.0.1"
var apiPort int = 8080
var apiPath string = "/api/json/"
var jenkinsAPITimeout time.Duration = 10
var jenkinsUsername string = "admin"
var jenkinsPassword string = "09269d8d3892403299b61ad47795fbbe"
var jenkinsAPIToken string = "113c32fc7f2833993ec339268f87f5a664"
var jenkinsWithCreds bool = true
var exporterHost = "127.0.0.1"
var exporterPort = 5000

func init() {

	config.Global.SSLOn = sslOn
	config.Global.JenkinsAPIHost = apiHost
	config.Global.JenkinsAPIPort = apiPort
	config.Global.JenkinsAPIPath = apiPath
	config.Global.JenkinsAPITimeout = jenkinsAPITimeout
	config.Global.JenkinsUsername = jenkinsUsername
	config.Global.JenkinsPassword = jenkinsPassword
	config.Global.JenkinsToken = jenkinsAPIToken
	config.Global.JenkinsWithCreds = jenkinsWithCreds
	config.Global.ExporterHost = exporterHost
	config.Global.ExporterPort = exporterPort

}
