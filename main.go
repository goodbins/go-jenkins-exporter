package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abousselmi/go-jenkins-exporter/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
)

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

// Get data using jenkins API
// To request data, we need the following information
//		host: jenkins host
//		port: jenkins port
//		username: jenkins username
//		password: jenkins password
//
var apiURL string = "http://127.0.0.1:8080/api/json/"
var jenkinsAPITimeout time.Duration = 10
var jenkinsUsername string = "admin"
var jenkinsPassword string = "09269d8d3892403299b61ad47795fbbe"
var jenkinsAPIToken string = "113c32fc7f2833993ec339268f87f5a664"

// Jenkins exporter host and port
var myHost = "127.0.0.1"
var myPort = "5000"

// Types and Structs
type jActions struct {
	Class                 string `json:"_class"`
	QueuingDurationMillis string `json:"queuingDurationMillis"`
	TotalDurationMillis   string `json:"totalDurationMillis"`
	SkipCount             string `json:"skipCount"`
	FailCount             string `json:"failCount"`
	TotalCount            string `json:"totalCount"`
	PassCount             string `json:"passCount"`
}

type jStatus struct {
	Class     string `json:"_class"`
	Actions   []jActions
	Duration  int `json:"duration"`
	Number    int `json:"number"`
	Timestamp int `json:"timestamp"`
}

type job struct {
	Class                 string  `json:"_class"`
	FullName              string  `json:"fullName"`
	URL                   string  `json:"url"`
	LastBuild             jStatus `json:"lastBuild"`
	LastCompletedBuild    jStatus `json:"lastCompletedBuild"`
	LastFailedBuild       jStatus `json:"lastFailedBuild"`
	LastStableBuild       jStatus `json:"lastStableBuild"`
	LastSuccessfulBuild   jStatus `json:"lastSuccessfulBuild"`
	LastUnstableBuild     jStatus `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild jStatus `json:"lastUnsuccessfulBuild"`
}

type jenkinsResponse struct {
	Class string `json:"_class"`
	Jobs  []job  `json:"jobs"`
}

// Init logging, prometheus collectors, etc.
func init() {
	//init logging
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
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
	go setGauges()

	// Create the
	go getData(apiURL, createQuery(), jenkinsAPITimeout, jenkinsUsername, jenkinsPassword)

	logrus.Info("Listning on " + myHost + " port " + myPort + " ...")
	logrus.Fatal(http.ListenAndServe(myHost+":"+myPort, r))
}

func getData(myurl, query string, timeout time.Duration, username, password string) map[string]interface{} {
	// Init a map whose keys are strings and whose values are themselves
	// stored as empty interface values
	var jResp jenkinsResponse
	// Init an http client
	httpClient := &http.Client{Timeout: timeout * time.Second}
	// Init a http request, set basic auth and Do the request
	req, err := http.NewRequest("GET", url.QueryEscape(myurl+query), nil)

	//var bearer = "Bearer " + jenkinsAPIToken

	req.SetBasicAuth(username, jenkinsAPIToken)
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Bearer", jenkinsAPIToken)
	resp, err := httpClient.Do(req)
	// Panic if an error occurs
	if err != nil {
		logrus.Error("An error has occured when getting: %s", apiURL)
		panic("An error has occured when trying to reach jenkins")
	}
	defer resp.Body.Close()

	// Decode: the better version
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &jResp)

	/*// Decode the response into the result map
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)*/
	logrus.Debug("Jenkins API returned: %v", jResp)

	for _, job := range jResp.Jobs {
		logrus.Debug("Job name: %v", job.FullName)
	}

	res := pretty.Pretty(body)
	logrus.Info(string(res))

	/*
			var jobs map[string]interface{}

			for job := range result["jobs"] {

			}
		            for job in result['jobs']:
		                if job['_class'] == 'com.cloudbees.hudson.plugins.folder.Folder' or \
		                   job['_class'] == 'jenkins.branch.OrganizationFolder' or \
		                   job['_class'] == 'org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject':
		                    jobs += parsejobs(job['url'] + '/api/json')
		                else:
							jobs.append(job)
	*/
	// Hand back the result
	return nil
}

// Update prometheus Gauges
func setGauges() {
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

func createQuery() string {

	var jobStatuses = []string{
		"lastBuild",
		"lastCompletedBuild",
		"lastFailedBuild",
		"lastStableBuild",
		"lastSuccessfulBuild",
		"lastUnstableBuild",
		"lastUnsuccessfulBuild",
	}

	var jobStatusProperties string = `[
		fullName,
		number,
		timestamp,
		duration,
		actions[
			queuingDurationMillis,
			totalDurationMillis,
			skipCount,
			failCount,
			totalCount,
			passCount]]`

	var query string
	for _, s := range jobStatuses {
		query += "," + s + jobStatusProperties
	}
	return strings.ReplaceAll(strings.ReplaceAll(
		fmt.Sprintf("?tree=jobs[fullName,url%s]", query),
		"\n", ""),
		"\t", "")
}
