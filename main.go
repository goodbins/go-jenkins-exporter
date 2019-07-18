package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/abousselmi/go-jenkins-exporter/exporter"
	"github.com/abousselmi/go-jenkins-exporter/handlers"
	"github.com/abousselmi/go-jenkins-exporter/jenkins"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Init logging, prometheus collectors, etc.
func init() {
	//init logging
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

// Main
func main() {
	// Print start message
	logrus.Debug("Starting go jenkins exporter")

	// Add an API router
	r := mux.NewRouter()
	r.HandleFunc("/ping", handlers.Ping).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	// Update go routine
	logrus.Debug("Launching metrics update loop")
	go exporter.SetGauges()

	// Get data from Jenkins

	go jenkins.GetData()

	logrus.Info("Listning on " + config.Global.ExporterHost + " port " + strconv.Itoa(config.Global.ExporterPort) + " ...")
	logrus.Fatal(http.ListenAndServe(config.Global.ExporterHost+":"+strconv.Itoa(config.Global.ExporterPort), r))
}
