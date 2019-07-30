package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
type JenkinsResponse struct {
	Class string `json:"_class"`
	Jobs  []job  `json:"jobs"`
}

var jenkinsFolderClasses = []string{
	"com.cloudbees.hudson.plugins.folder.Folder",
	"jenkins.branch.OrganizationFolder",
	"org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject"}

var jobsList []job                 // List of discovered jobs
var jobFloderLinks []string        // List of job folders
var jobFolderVisitedLinks []string // List of visited/explored folders

func GetData() *[]job {
	logrus.Debug("Get data from jenkins..")
	walkAndGetJobs(getJenkinsApiUrl())
	logrus.Debug("Data retrieved successfully")
	fmt.Println("jobsList ", jobsList)
	fmt.Println("jobFolderLinks ", jobFloderLinks)
	return &jobsList
}

// First url is the API's
func walkAndGetJobs(url string) {
	logrus.Debug("Walking ", url)
	jobs := requestJson(url + "api/json" + createQuery())
	jobFolderVisitedLinks = append(jobFolderVisitedLinks, url)
	updateJobsAndFolders(jobs, &jobsList, &jobFloderLinks)
	for _, fL := range jobFloderLinks {
		if !isVisited(&fL) {
			walkAndGetJobs(fL)
		}
	}
}

func updateJobsAndFolders(reply, jL *[]job, jF *[]string) {
	for _, j := range *reply {
		if isJobsFolder(&j.Class) {
			*jF = append(*jF, j.URL)
			continue
		}
		*jL = append(*jL, j)
	}
}

func isJobsFolder(class *string) bool {
	for _, c := range jenkinsFolderClasses {
		if *class == c {
			return true
		}
	}
	return false
}

func isVisited(link *string) bool {
	for _, j := range jobFolderVisitedLinks {
		if *link == j {
			return true
		}
	}
	return false
}

func requestJson(url string) *[]job {
	// Init a map whose keys are strings and whose values are themselves
	// stored as empty interface values
	var jResp JenkinsResponse
	resp := request(url)
	defer resp.Body.Close()
	// Decode to json the jenkins reply
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	err = json.Unmarshal(body, &jResp)
	if err != nil {
		logrus.Error("An error has occured while decoding JSON")
		panic("An error has occured while decoding JSON")
	}
	return &jResp.Jobs
}

func request(apiurl string) *http.Response {
	// Init an http client
	httpClient := &http.Client{Timeout: config.Global.JenkinsAPITimeout * time.Second}
	// Init a http request, set basic auth and Do the request
	req, err := http.NewRequest("GET", apiurl, nil)
	// Test if credentials are used
	if config.Global.JenkinsWithCreds {
		if config.Global.JenkinsPassword != "" {
			req.SetBasicAuth(config.Global.JenkinsUsername, config.Global.JenkinsPassword)
		}
		if config.Global.JenkinsToken != "" {
			req.SetBasicAuth(config.Global.JenkinsUsername, config.Global.JenkinsToken)
		}
	}
	// Make the request
	resp, err := httpClient.Do(req)
	// Panic if an error occurs
	if err != nil {
		logrus.Error("An error has occured when getting: ", apiurl)
		panic("An error has occured when trying to reach jenkins")
	}
	// Return the Jenskins response
	return resp
}

func getJenkinsApiUrl() string {
	var apiurl string = "http://"
	if config.Global.SSLOn {
		apiurl = "https://"
	}
	apiurl += config.Global.JenkinsAPIHostPort + "/"
	return apiurl
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
