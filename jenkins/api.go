package jenkins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/abousselmi/go-jenkins-exporter/config"
	"github.com/sirupsen/logrus"
)

// Jenkins job actions struct
type jActions struct {
	Class                 string `json:"_class"`
	QueuingDurationMillis string `json:"queuingDurationMillis"`
	TotalDurationMillis   string `json:"totalDurationMillis"`
	SkipCount             string `json:"skipCount"`
	FailCount             string `json:"failCount"`
	TotalCount            string `json:"totalCount"`
	PassCount             string `json:"passCount"`
}

// Jenkins job statuses struct
type jStatus struct {
	Class     string `json:"_class"`
	Actions   []jActions
	Duration  int `json:"duration"`
	Number    int `json:"number"`
	Timestamp int `json:"timestamp"`
}

// Jenkins job struct
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

// Jenkins API response struct
type jenkinsResponse struct {
	Class string `json:"_class"`
	Jobs  []job  `json:"jobs"`
}

func GetData() {
	// Init a map whose keys are strings and whose values are themselves
	// stored as empty interface values
	var apiurl string = "http://"
	if config.Global.SSLOn {
		apiurl = "https://"
	}
	apiurl += strconv.Itoa(config.Global.JenkinsAPIPort) + config.Global.JenkinsAPIHost + config.Global.JenkinsAPIPath
	logrus.Debug("Getting data from Jenkins API")
	var jResp jenkinsResponse
	resp := request(
		config.Global.JenkinsWithCreds,
		config.Global.JenkinsUsername,
		config.Global.JenkinsPassword,
		config.Global.JenkinsToken,
		apiurl, createQuery(),
		config.Global.JenkinsAPITimeout)
	// Decode to json jenkins reply
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &jResp)
	if err != nil {
		logrus.Error("An error has occured while decoding JSON")
		panic("An error has occured while decoding JSON")
	}

	/* Debugging */
	// Decode the response into the result map
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	// logrus.Debug("Jenkins API returned: %v", jResp)
	for _, job := range jResp.Jobs {
		logrus.Debug("Job name: %v", job.FullName)
	}
	// res := pretty.Pretty(body)
	// logrus.Info(string(res))
}

func request(withcreds bool, username, password, token, apiurl, query string, timeout time.Duration) *http.Response {
	// Init an http client
	httpClient := &http.Client{Timeout: timeout * time.Second}
	// Init a http request, set basic auth and Do the request
	req, err := http.NewRequest("GET", url.QueryEscape(apiurl+query), nil)
	// Some tests
	if withcreds {
		if password != "" {
			req.SetBasicAuth(username, password)
		}
		if token != "" {
			req.SetBasicAuth(username, token)
		}
	}
	// Make the request
	resp, err := httpClient.Do(req)
	// Panic if an error occurs
	if err != nil {
		logrus.Error("An error has occured when getting: %s", apiurl)
		panic("An error has occured when trying to reach jenkins")
	}
	defer resp.Body.Close()
	// Return the Jenskins response
	return resp
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
