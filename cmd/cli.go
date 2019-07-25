package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Init the CLI
func init() {

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

// RootCommand Run the cobra command line program
func RootCommand() *cobra.Command {
	cobraCmd := cobra.Command{
		Use:     "go-jenkins-exporter",
		Long:    "A simple jenkins exporter for prometheus, written in Go.",
		Run:     run,
		Version: config.CurrentVersion,
	}

	// Define and init flags
	cobraCmd.Flags().BoolVarP(&config.Global.SSLOn, "ssl", "s", false, "Enable TLS (default false)")                                    // Optional
	cobraCmd.Flags().StringVar(&config.Global.JenkinsAPIHost, "jhost", "", "Jenkins host")                                              // Mendatory
	cobraCmd.Flags().IntVar(&config.Global.JenkinsAPIPort, "jport", 8080, "Jenkins port")                                               // Optional
	cobraCmd.Flags().StringVar(&config.Global.JenkinsAPIPath, "path", "/api/json", "Jenkins API path")                                  // Optional
	cobraCmd.Flags().DurationVarP(&config.Global.JenkinsAPITimeout, "timeout", "t", 10*time.Second, "Jenkins API timeout in seconds")   // Optional
	cobraCmd.Flags().StringVar(&config.Global.ExporterHost, "host", "localhost", "Exporter host")                                       // Optional
	cobraCmd.Flags().IntVar(&config.Global.ExporterPort, "port", 5000, "Exporter port")                                                 // Optional
	cobraCmd.Flags().StringVarP(&config.Global.MetricsPath, "metrics", "m", "/metrics", "Path under which to expose metrics")           // Optional
	cobraCmd.Flags().DurationVarP(&config.Global.MetricsUpdateRate, "rate", "r", 1*time.Second, "Set metrics update rate in seconds")   // Optional
	cobraCmd.Flags().StringVarP(&config.Global.LogLevel, "log", "l", "info", "Logging level: info, debug, warn, error, panic or fatal") // Optional
	viper.BindEnv("username", "JENKINS_USERNAME")                                                                                       // Mendatory
	viper.BindEnv("password", "JENKINS_PASSWORD")                                                                                       // Optional/Mendatory
	viper.BindEnv("token", "JENKINS_TOKEN")                                                                                             // Optional/Mendatory
	config.Global.JenkinsUsername = viper.GetString("username")
	config.Global.JenkinsPassword = viper.GetString("password")
	config.Global.JenkinsToken = viper.GetString("token")
	config.Global.JenkinsWithCreds = false
	return &cobraCmd
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
	ok := checkFlags()
	if !ok {
		fmt.Println("Use --help to get more info...")
		os.Exit(1)
	}
	config.SetupLogging()
}

func checkFlags() bool {
	/* Check if mendatory flags are set */
	// Check jenkins address
	if config.Global.JenkinsAPIHost == "" {
		fmt.Println("Jenkins host address is missing !")
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
	if config.Global.JenkinsAPIPort < 1024 || config.Global.ExporterPort < 1024 {
		fmt.Println("Privileged ports are not supported. Choose one bigger than 1024...")
		return false
	}

	// Check if the provided config level is in acceptable config levels
	if _, ok := config.LogLevels[config.Global.LogLevel]; !ok {
		fmt.Println("Accepted log levels are: info, debug, warn, error, panic or fatal")
		return false
	}

	return true
}
