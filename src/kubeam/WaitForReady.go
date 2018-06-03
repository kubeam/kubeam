package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

/*ApplicationWaitForReady ...*/
func ApplicationWaitForReady(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	application := vars["application"]
	appEnv := vars["environment"]
	appCluster := vars["cluster"]

	clusterList, err := DBGetClusterReservation(redisClient, application, appEnv, appCluster)
	LogTrace.Println(clusterList)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Unable to select cluster for specified environment"}`
		w.Write([]byte(str))
	} else {
		cmdName := "./kubectl"
		cmdArgs := []string{"rollout", "status", "deployment"}

		resourceName := fmt.Sprintf("%s-%s-c%s-%s", application, appEnv, appCluster, "app")
		resourceArgs := append(cmdArgs, resourceName)
		LogTrace.Printf("Running : %s %s", cmdName, resourceArgs)
		cmdOut, _ := exec.Command(cmdName, resourceArgs...).CombinedOutput()

		w.Header().Set("Content-Type", "application/json")
		var str string
		if strings.Index(string(cmdOut), "successfully rolled out") >= 0 {
			str = fmt.Sprintf(`{"status": "success", "resource": "%v", "description": "Successfully rolled out"}`, resourceName)
		} else {
			str = `{"status": "error", "description": "timedout during rollout"}`
		}
		w.Write([]byte(str))
	}
}
