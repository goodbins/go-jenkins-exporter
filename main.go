package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/abousselmi/go-jenkins-exporter/cmd"
	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/abousselmi/go-jenkins-exporter/handlers"
	"github.com/abousselmi/go-jenkins-exporter/prom"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func init() {
	// Get CLI args
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
	// Setup logging
	setupLogging()
}

// Main
func main() {
	// Print start message
	logrus.Info("Starting go-jenkins-exporter")

	// Add an API router
	r := mux.NewRouter()
	r.HandleFunc("/ping", handlers.Ping).Methods("GET")
	r.Handle(config.Global.MetricsPath, promhttp.Handler())

	// Launch metrics update go routine
	go prom.SetGauges()

	// Listen and serve
	logrus.Info("Listning on " + config.Global.ExporterHost + " port " + strconv.Itoa(config.Global.ExporterPort) + " ...")
	logrus.Fatal(http.ListenAndServe(config.Global.ExporterHost+":"+strconv.Itoa(config.Global.ExporterPort), r))
}

func setupLogging() {
	logrus.SetOutput(os.Stdout) // FIXME: see cmd package
	logrus.SetLevel(config.LogLevels[config.Global.LogLevel])
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}
