package exporter

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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

// Update prometheus Gauges
func SetGauges() {
	for {
		numberGauge.Set(rand.Float64())
		durationGauge.Set(rand.Float64())
		timestampGauge.Set(rand.Float64())
		queuingDurationMillisGauge.Set(rand.Float64())
		totalDurationMillisGauge.Set(rand.Float64())
		skipCountGauge.Set(rand.Float64())
		failCountGauge.Set(rand.Float64())
		totalCountGauge.Set(rand.Float64())
		passCountGauge.Set(rand.Float64())
		time.Sleep(1 * time.Second)
	}
}

// Create prometheus Gauges
var (
	numberGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "number",
		Help: "Jenkins build number",
	})
	durationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "duration",
		Help: "Jenkins build duration in seconds",
	})
	timestampGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "timestamp",
		Help: "Jenkins build timestamp in unixtime",
	})
	queuingDurationMillisGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "queuingDurationMillis",
		Help: "Jenkins build queuing duration in seconds",
	})
	totalDurationMillisGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "totalDurationMillis",
		Help: "Jenkins build total duration in seconds",
	})
	skipCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "skipCount",
		Help: "Jenkins build skip counts",
	})
	failCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "failCount",
		Help: "Jenkins build fail counts",
	})
	totalCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "totalCount",
		Help: "Jenkins build total counts",
	})
	passCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "passCount",
		Help: "Jenkins build pass counts",
	})
)
