package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/*ApplicationCreate ...*/
func ApplicationCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// application := vars["application"]
	// appEnv := vars["environment"]
	cluster, ok := vars["cluster"]
	if !ok {
		LogError.Println("Cluster/Shard parameter is required for creating resources")
	} else {
		LogInfo.Println("Setting cluster ", cluster)
	}

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	actionsOutput, err := RunActions("/v1/create", m)

	/****/

	time.Sleep(2000 * time.Millisecond)
	// payload, _ := GetResourceStatus(vars, []string{
	// 	"adminweb",
	// 	"appweb",
	// 	"admin",
	// 	"app",
	// 	"db",
	// })

	w.Header().Set("Content-Type", "application/json")
	outputJSON, _ := json.MarshalIndent(actionsOutput, "", " ")
	w.Write(outputJSON)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

}
