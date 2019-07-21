package prom

import (
	"math/rand"
	"time"

	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/abousselmi/go-jenkins-exporter/jenkins"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func init() {
	//register prometheus metric collectors
	prometheus.MustRegister(numberGauge)
	prometheus.MustRegister(durationGauge)
	prometheus.MustRegister(timestampGauge)
	prometheus.MustRegister(queuingDurationMillisGauge)
	prometheus.MustRegister(totalDurationMillisGauge)
	prometheus.MustRegister(skipCountGauge)
	prometheus.MustRegister(failCountGauge)
	prometheus.MustRegister(totalCountGauge)
	prometheus.MustRegister(passCountGauge)
	//init prometheus gauges
	numberGauge.Set(0.0)
	durationGauge.Set(0.0)
	timestampGauge.Set(0.0)
	queuingDurationMillisGauge.Set(0.0)
	totalDurationMillisGauge.Set(0.0)
	skipCountGauge.Set(0.0)
	failCountGauge.Set(0.0)
	totalCountGauge.Set(0.0)
	passCountGauge.Set(0.0)
}

// Get data from Jenkins and update prometheus metrics
func SetGauges() {
	logrus.Debug("Launching metrics update loop: updating rate is set to ", config.Global.MetricsUpdateRate)
	for {
		var jResp jenkins.JenkinsResponse = jenkins.GetData()
		for _, job := range jResp.Jobs {

			jobname = job.FullName

			numberGauge.Set(rand.Float64())
			durationGauge.Set(rand.Float64())
			timestampGauge.Set(rand.Float64())
			queuingDurationMillisGauge.Set(rand.Float64())
			totalDurationMillisGauge.Set(rand.Float64())
			skipCountGauge.Set(rand.Float64())
			failCountGauge.Set(rand.Float64())
			totalCountGauge.Set(rand.Float64())
			passCountGauge.Set(rand.Float64())
		}
		time.Sleep(config.Global.MetricsUpdateRate)
	}
}

var jobname string

// Create prometheus Gauges
var (
	numberGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "number",
		Help: "Jenkins build number",
	})
	durationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "duration",
		Help: "Jenkins build duration in seconds",
	})
	timestampGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "timestamp",
		Help: "Jenkins build timestamp in unixtime",
	})
	queuingDurationMillisGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "queuingDurationMillis",
		Help: "Jenkins build queuing duration in seconds",
	})
	totalDurationMillisGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "totalDurationMillis",
		Help: "Jenkins build total duration in seconds",
	})
	skipCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "skipCount",
		Help: "Jenkins build skip counts",
	})
	failCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "failCount",
		Help: "Jenkins build fail counts",
	})
	totalCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "totalCount",
		Help: "Jenkins build total counts",
	})
	passCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jenkins_job_" + "passCount",
		Help: "Jenkins build pass counts",
	})
)

var prometheusMetrics map[string]prometheus.Gauge
