package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ApplicationDeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//    application := vars["application"]
	//    appEnv := vars["environment"]
	cluster, ok := vars["cluster"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		str := `{"status": "error", "description": "Please specify a target cluster"}`
		w.Write([]byte(str))
		LogWarning.Println("Cluster was not specified in request")
		return
	} else {
		LogInfo.Println("Setting cluster ", cluster)
	}

	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	actionsOutput, err := RunActions("/v1/deploy", m)

	w.Header().Set("Content-Type", "application/json")
	outputJSON, _ := json.MarshalIndent(actionsOutput, "", " ")
	w.Write(outputJSON)
	//w.Header().Set("Content-Type", "application/text")

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	// //w.Write(payload)
	// LogInfo.Println(actionsOutput)

}
