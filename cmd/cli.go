package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goodbins/go-jenkins-exporter/config"
	"github.com/goodbins/go-jenkins-exporter/exporter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the cobraCmd.
func Execute() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

// RootCommand Run the cobra command line program
func RootCommand() *cobra.Command {
	cobraCmd := cobra.Command{
		Use: "go-jenkins-exporter",
		Long: `A simple jenkins exporter for prometheus, written in Go.

Note: To setup jenkins credentials, use these environment variables:
JENKINS_USERNAME, JENKINS_PASSWORD and/or JENKINS_TOKEN
If they are not set, we assume no credentials.`,
		Run:     run,
		Version: config.CurrentVersion,
	}

	// Define and init flags
	cobraCmd.Flags().BoolVarP(&config.Global.SSLOn, "ssl", "s", false, "Enable TLS (default false)")                                  // Optional
	cobraCmd.Flags().StringVarP(&config.Global.JenkinsAPIHostPort, "jenkins", "j", "", "Jenkins API host:port pair")                  // Mendatory
	cobraCmd.Flags().StringVarP(&config.Global.JenkinsAPIPath, "path", "a", "/api/json", "Jenkins API path")                          // Optional
	cobraCmd.Flags().DurationVarP(&config.Global.JenkinsAPITimeout, "timeout", "t", 10*time.Second, "Jenkins API timeout in seconds") // Optional
	cobraCmd.Flags().StringVarP(&config.Global.ExporterHostPort, "listen", "l", "localhost:5000", "Exporter host:port pair")          // Optional
	cobraCmd.Flags().StringVarP(&config.Global.MetricsPath, "metrics", "m", "/metrics", "Path under which to expose metrics")         // Optional
	cobraCmd.Flags().DurationVarP(&config.Global.MetricsUpdateRate, "rate", "r", 1*time.Second, "Set metrics update rate in seconds") // Optional
	cobraCmd.Flags().BoolVarP(&config.Global.Verbose, "verbose", "v", false, "Enable verbosity")                                      // Optional
	cobraCmd.Flags().StringVar(&config.Global.LogLevel, "log", "info", "Log level, one of: info, debug, warn, error, fatal")          // Optional
	viper.BindEnv("username", "JENKINS_USERNAME")                                                                                     // Optional/Mendatory
	viper.BindEnv("password", "JENKINS_PASSWORD")                                                                                     // Optional/Mendatory
	viper.BindEnv("token", "JENKINS_TOKEN")                                                                                           // Optional/Mendatory
	config.Global.JenkinsUsername = viper.GetString("username")
	config.Global.JenkinsPassword = viper.GetString("password")
	config.Global.JenkinsToken = viper.GetString("token")
	config.Global.JenkinsWithCreds = false
	return &cobraCmd
}

func run(cmd *cobra.Command, args []string) {
	ok := checkFlags()
	if !ok {
		fmt.Println("Use --help to get more info...")
		os.Exit(1)
	}
	config.SetupLogging()
	exporter.Serve()
}

func checkFlags() bool {
	/* Check if mendatory flags are set */
	// Check jenkins address
	if config.Global.JenkinsAPIHostPort == "" {
		fmt.Println("Jenkins host:port address is missing !")
		return false
	}

	// Check if jenkins credentials are ok
	if config.Global.JenkinsUsername == "" && (config.Global.JenkinsPassword != "" || config.Global.JenkinsToken != "") {
		fmt.Println("You provided an empty username !")
		return false
	} else if config.Global.JenkinsUsername != "" && config.Global.JenkinsPassword == "" && config.Global.JenkinsToken == "" {
		fmt.Println("You need to provide either a password or a token !")
		return false
	} else if config.Global.JenkinsUsername+config.Global.JenkinsPassword+config.Global.JenkinsToken == "" {
		fmt.Println("Connecting to jenkins without credentials !")
	} else {
		config.Global.JenkinsWithCreds = true
	}

	/* Check other flags */
	// Check if ports are not privileged
	jPort, _ := strconv.Atoi(strings.Split(config.Global.JenkinsAPIHostPort, ":")[1])
	ePort, _ := strconv.Atoi(strings.Split(config.Global.ExporterHostPort, ":")[1])
	if jPort < 1024 || ePort < 1024 {
		fmt.Println("Privileged ports are not supported. Choose one bigger than 1024...")
		return false
	}

	return true
}
