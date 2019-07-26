package main

import (
	"net/http"
	"strconv"

	"github.com/abousselmi/go-jenkins-exporter/cmd"
	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/abousselmi/go-jenkins-exporter/handlers"
	"github.com/abousselmi/go-jenkins-exporter/prom"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()
}

func serve() {
	// Print start message
	logrus.Info("Starting go-jenkins-exporter")

	// Launch metrics update go routine
	go prom.SetGauges()

	// Handle routes: / /ping /metrics
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Go Jenkins Exporter</title></head>
		<body>
		<h1>Go Jenkins Exporter</h1>
		<p><a href="` + config.Global.MetricsPath + `">Metrics</a></p>
		</body></html>`))
	})
	http.HandleFunc("/ping", handlers.Ping)
	http.Handle(config.Global.MetricsPath, promhttp.Handler())

	// Listen and serve
	logrus.Info("Listning on " + config.Global.ExporterHost + " port " + strconv.Itoa(config.Global.ExporterPort) + " ...")
	logrus.Fatal(http.ListenAndServe(config.Global.ExporterHost+":"+strconv.Itoa(config.Global.ExporterPort), nil))
}
