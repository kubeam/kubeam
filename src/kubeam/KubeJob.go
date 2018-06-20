package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	mux "github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*RunJob parses payload to execute the cron command and it's arguments*/
func RunJob(w http.ResponseWriter, r *http.Request) {
	var cronData map[string]string

	m := make(map[string]interface{})
	vars := mux.Vars(r)

	byteData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(byteData, &cronData)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		vars["jobcommand"] = cronData["jobcommand"]
		vars["jobparams"] = cronData["jobparams"]
		for k, v := range vars {
			m[k] = v
		}

		res, errcode := GetDockerTag("/v1/kubejob", m)

		if errcode == 0 {
			m["tag"] = res["Tag"]
			actionsOutput, _ := RunActions("/v1/kubejob", m)
			output, _ := json.MarshalIndent(actionsOutput, "", " ")
			response := string(output)
			w.Write([]byte(response))
		} else {
			response, err := json.Marshal(res)
			ErrorHandler(err)
			w.WriteHeader(errcode)
			w.Write([]byte(response))
		}
	}
}

/*GetJobStatus Get execution status of a Job running on cluster*/
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	var resource, namespace string
	var response []byte

	res := make(map[string]interface{})
	vars := mux.Vars(r)

	clientset := GetClientSet()

	namespace = GetJobNamespace(vars)
	resource = fmt.Sprintf("%s-%s-c%s-job-%s", vars["application"], vars["environment"], vars["cluster"], vars["jobname"])

	LogInfo.Printf("get jobstatus: %s on namepsace: %s", resource, namespace)
	myjob, err := clientset.BatchV1().Jobs(namespace).Get(resource, metav1.GetOptions{})
	ErrorHandler(err)

	if err == nil {
		jobStatus := myjob.Status
		res["JobName"] = resource

		res["StartTime"] = strings.Replace(jobStatus.StartTime.String(), "T", " ", 1)[:20]
		res["JobName"] = resource
		res["JobId"] = fmt.Sprintf("%s-%d", resource, jobStatus.StartTime.Unix())
		if len(jobStatus.Conditions) != 0 {
			res["LastProbeTime"] = strings.Replace(jobStatus.CompletionTime.String(), "T", " ", 1)[:20]
		}

		if jobStatus.Active == 0 &&
			jobStatus.Failed == 0 &&
			jobStatus.Succeeded != 0 {
			res["JobStatus"] = "Completed"
			res["Logs"] = GetLogs(myjob)
		} else if jobStatus.Active == 0 &&
			jobStatus.Failed != 0 &&
			jobStatus.Succeeded == 0 {
			res["JobStatus"] = "Failed"
			res["Logs"] = GetLogs(myjob)
		} else {
			res["JobStatus"] = "Running"
			res["Logs"] = "No Logs"
		}
		res["LastProbeTime"] = strings.Replace(time.Now().String(), "T", " ", 1)[:20]
	} else {
		LogInfo.Println(err, myjob)
		res["JobStatus"] = "Failed"
		res["Logs"] = err.Error()
	}
	response, err = json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

/*DeleteJob deletes a kubernetes job from given app-env-cluster*/
func DeleteJob(w http.ResponseWriter, r *http.Request) {
	graceperiod := int64(0)
	clientset := GetClientSet()
	vars := mux.Vars(r)
	resource := fmt.Sprintf("%s-%s-c%s-job-%s", vars["application"], vars["environment"], vars["cluster"], vars["jobname"])
	namespace := GetJobNamespace(vars)
	LogInfo.Printf("Deleting Job ==> %s", resource)
	err := clientset.BatchV1().Jobs(namespace).Delete(resource, &metav1.DeleteOptions{GracePeriodSeconds: &graceperiod})
	if err != nil {
		LogError.Println(err.Error())
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("Job %s Deleted", resource)))
	}
}
