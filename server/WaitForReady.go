package server

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	//"github.com/bitly/go-simplejson"
	"io/ioutil"
	"github.com/kubeam/kubeam/common"
)

// ApplicationWaitForReady - Wait for a application to be fully deployed. this is a sync call.
func ApplicationWaitForReady(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]
	appCluster := vars["cluster"]

	clusterList, err := DBGetClusterReservation(redisClient, application, appEnv, appCluster)
	common.LogTrace.Println(clusterList)

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("clusters/%v-%v-clusterlist.yaml", application, appEnv))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Could not find a cluster definition for application"}`
		w.Write([]byte(str))
		return
	}

	common.LogInfo.Println(yamlFile)
	cmdName := "./kubectl"
	cmdArgs := []string{"--namespace", "qbo", "rollout", "status", "deployment"}

	resourceName := fmt.Sprintf("%s-%s-c%s-%s", application, appEnv, appCluster, "app")
	resourceArgs := append(cmdArgs, resourceName)
	common.LogTrace.Println(fmt.Sprintf("Running : %s %s", cmdName, resourceArgs))
	cmdOut, _ := exec.Command(cmdName, resourceArgs...).CombinedOutput()

	w.Header().Set("Content-Type", "application/json")
	if strings.Index(string(cmdOut), "successfully rolled out") == -1 {
		str := `{"status": "error", "description": "timedout during rollout"}`
		w.Write([]byte(str))
		return

	}
	str := fmt.Sprintf(`{"status": "success", "resource": "%v", "description": "Successfully rolled out"}`, resourceName)
	w.Write([]byte(str))

}
