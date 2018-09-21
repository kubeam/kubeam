package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kubeam/kubeam/common"
	"github.com/kubeam/kubeam/services"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

/*RunJob parses payload to execute the cron command and it's arguments*/
func RunJob(w http.ResponseWriter, r *http.Request) {
	var cronData map[string]string
	var err error

	m := make(map[string]interface{})
	vars := mux.Vars(r)
	byteData, _ := ioutil.ReadAll(r.Body)

	if err = json.Unmarshal(byteData, &cronData); err == nil {
		vars["jobcommand"] = cronData["jobcommand"]
		vars["jobparams"] = cronData["jobparams"]
		for k, v := range vars {
			m[k] = v
		}

		if tag, err := services.GetDockerTag("/v1/kubejob", m); err == nil {
			common.LogDebug.Printf("Found docker tag %s", tag)
			m["tag"] = tag
			if m["namespace"], err = services.GetJobNamespace(vars); err == nil {
				actionsOutput, _ := services.RunJobActions("/v1/kubejob", m)
				output, _ := json.MarshalIndent(actionsOutput, "", " ")
				w.Header().Set("Content-Type", "application/json")
				w.Write(output)
			}
		}
	} else {
		common.LogError.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte(err.Error()))
	}
}

/*DeleteJob deletes a kubernetes job from given app-env-cluster*/
func DeleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	vars := mux.Vars(r)
	resource := fmt.Sprintf("%s-%s-c%s-job-%s", vars["application"],
		vars["environment"], vars["cluster"], vars["jobname"])
	if namespace, err := services.GetJobNamespace(vars); err == nil {
		result := services.DeleteKubernetesJob(resource, namespace)
		w.Write([]byte(result))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

/*GetJobStatus Get execution status of a Job running on cluster*/
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	vars := mux.Vars(r)

	if jobapi, err := services.GetJobAPI(vars); err == nil {
		if resource, ok := jobapi["resource"]; ok {
			if ns, ok := jobapi["namespace"]; ok {

				client := services.GetClientSet().BatchV1().Jobs(ns.(string))
				common.LogInfo.Printf("get jobstatus: %s on namepsace: %s", resource, ns)

				if myjob, err := client.Get(resource.(string), metav1.GetOptions{}); err == nil {
					jobStatus := myjob.Status
					res["JobName"] = resource
					res["StartTime"] = strings.Replace(jobStatus.StartTime.String(), "T", " ", 1)[:20]
					res["JobId"] = fmt.Sprintf("%s-%d", resource, jobStatus.StartTime.Unix())
					res["LastProbeTime"] = strings.Replace(time.Now().String(), "T", " ", 1)[:20]

					if len(jobStatus.Conditions) != 0 {
						res["LastProbeTime"] = strings.Replace(jobStatus.CompletionTime.String(), "T", " ", 1)[:20]
					}

					if jobStatus.Active == 0 && jobStatus.Failed == 0 &&
						jobStatus.Succeeded != 0 {
						res["JobStatus"] = "Completed"
						res["common.Logs"], _ = services.GetLogs(myjob)
					} else if jobStatus.Active == 0 && jobStatus.Failed != 0 &&
						jobStatus.Succeeded == 0 {
						res["JobStatus"] = "Failed"
						res["common.Logs"], _ = services.GetLogs(myjob)
					} else {
						res["JobStatus"] = "Running"
						res["common.Logs"] = "No common.Logs"
					}
					if response, err := json.Marshal(res); err == nil {
						w.Header().Set("Content-Type", "application/json")
						w.Write(response)
					}
				} else {
					common.LogError.Println(err.Error())
					w.Header().Set("Content-Type", "application/text")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			} else {
				common.LogError.Println("412 - Failed to obtain resource")
				w.Header().Set("Content-Type", "application/text")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("412 - Failed to obtain resource"))
				return
			}
		} else {
			common.LogError.Println("412 - Failed to obtain namespace")
			w.Header().Set("Content-Type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("412 - Failed to obtain namespace"))
			return
		}
	} else {
		common.LogError.Println(err.Error())
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
