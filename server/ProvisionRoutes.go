package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/kubeam/kubeam/common"
	"github.com/kubeam/kubeam/services"
)

/*ApplicationStatus fetches the status of kubernetes application*/
func ApplicationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	payload, _ := GetResourceStatus(vars, []string{
		"admin",
		"app",
		"db",
		"appweb",
		"adminweb",
	})

	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/text")
	w.Write(payload)
}

/*SelfProvision ...*/
func SelfProvision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	appEnv, ok := vars["environment"]
	if !ok {
		vars["environment"] = "latest"
		appEnv = vars["environment"]
	} else {
		common.LogInfo.Println("Setting environment", appEnv)
	}

	tag, ok := vars["tag"]
	if !ok {
		vars["tag"] = "latest"
		tag = vars["tag"]
	} else {
		common.LogInfo.Println("Setting tag", tag)
	}

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	cmdName := "./kubectl"

	UpdateResources(vars, []string{
		"applications/kubeam/kubeam-deployment.yaml",
		//"applications/kubeam/kubeam-service.yaml",
		"applications/kubeam/kubeam-redis-deployment.yaml",
		"applications/kubeam/kubeam-redis-service.yaml",
	})

	// Due to redis using persistant storage We should not use replace configuration.
	// BUG/FIX: Until we have detection of what is already running. We issue a Create (if already exists just silently failes
	CreateResources(vars, []string{
		"applications/kubeam/kubeam-service.yaml",
		"applications/kubeam/kubeam-redis-deployment.yaml",
		"applications/kubeam/kubeam-redis-service.yaml",
	})

	time.Sleep(2000 * time.Millisecond)

	cmdArgs := []string{"get", "deployment", fmt.Sprintf("%s-%s", appEnv, "kubeam")}
	common.LogTrace.Println(fmt.Sprintf("Running : %s %s", cmdName, cmdArgs))
	cmdOut, err := exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		// this error is not critical
		common.LogWarning.Println("Error running kubectl to get status")
	}
	payload := cmdOut
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/text")
	w.Write(payload)
}

/*ApplicationProvision is a wrapper to deploy applications to kubernetes*/
func ApplicationProvision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]
	cluster, ok := vars["cluster"]
	if !ok {
		//ttl, err := time.ParseDuration (  ttl )
		//if err != nil {
		ttl := time.Duration(900 * time.Second)
		//}
		ttl = time.Duration(900 * time.Second)
		clusterNumber, err := DBClientFindAndReserve(redisClient, application, appEnv, ttl)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			str := `{"status": "error", "description": "Unable to select cluster for specified environment, No free slots?"}`
			w.Write([]byte(str))
			return
		}

		common.LogInfo.Println("Cluster Number", clusterNumber)
		vars["cluster"] = clusterNumber
		cluster = vars["cluster"]
	} else {
		common.LogInfo.Println("Setting cluster ", cluster)
	}

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	actionsOutput, err := services.RunActions("/v1/provision", m)

	actionsOutput["cluster"] = cluster

	w.Header().Set("Content-Type", "application/json")
	outputJSON, _ := json.MarshalIndent(actionsOutput, "", " ")
	w.Write(outputJSON)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	common.LogInfo.Println(actionsOutput)
}

/*ApplicationDelete is a wrapper to delete kuberbetes deployment*/
func ApplicationDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}
	actionsOutput, err := services.RunActions("/v1/delete", m)

	w.Header().Set("Content-Type", "application/json")
	outputJSON, _ := json.MarshalIndent(actionsOutput, "", " ")
	w.Write(outputJSON)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	common.LogInfo.Println(actionsOutput)
}

/*ApplicationUpdate ...*/
func ApplicationUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]
	cluster := vars["cluster"]
	build := vars["build"]
	ttl, ok := vars["ttl"]
	if !ok {
		var err error
		ttl, err = common.Config.GetString("application/default/ttl", "600")
		if err != nil {
			common.LogInfo.Println(err)
		}
	}

	json := simplejson.New()
	json.Set("application", application)
	json.Set("environment", appEnv)
	json.Set("cluster", cluster)
	json.Set("tag", build)
	json.Set("ttl", ttl)

	payload, err := json.MarshalJSON()
	if err != nil {
		common.LogInfo.Println(err)
	}

	//w.Header().Set("Content-Type", "application/text")
	w.Write(payload)
}
