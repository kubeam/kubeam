package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"

	"github.com/gorilla/mux"
)

// ["https://admin:@localhost:8443/v1/event/schedulecron", "POST", "{'cmd': 'echo', 'params': 'hello world'}"]

func ParseCronDetails(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	var cronData map[string]string
	// var actionsOutput map[string]interface{}
	var response string

	m := make(map[string]interface{})
	vars := mux.Vars(r)

	byteData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(byteData, &cronData)
	if err != nil {
		LogInfo.Println(err.Error())
	} else {
		cmd := cronData["cmd"]
		params := cronData["params"]
		// if commandSanityCheck(cmd, params) == false {
		// 	statusCode = 400
		// 	response = "Too Many Commands/Arguments"
		// } else {

		for k, v := range vars {
			m[k] = v
		}
		// send the cron command along with parameters to RunAction
		m["jobcommand"] = cmd
		m["jobparams"] = params

		actionsOutput, _ := RunActions("/v1/runjob", m)
		w.Header().Set("Content-Type", "application/json")
		outputJSON, _ := json.MarshalIndent(actionsOutput, "", " ")
		w.WriteHeader(statusCode)
		w.Write(outputJSON)
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

// security check to disable running multiple operations on command line
func commandSanityCheck(cmd, params string) bool {
	LogInfo.Println(cmd, params)
	re, err := regexp.Compile("[\\|;&$><`\\!]")
	if err != nil {
		LogInfo.Println(err.Error())
		return false
	}

	if cmd != "" {
		matchCmd := re.FindString(cmd)
		matchParams := re.FindString(params)
		LogInfo.Println("SUSPICOUS CHARACTES:")
		LogInfo.Println(matchCmd, reflect.TypeOf(matchCmd), len(matchCmd))
		LogInfo.Println(matchParams, reflect.TypeOf(matchCmd), len(matchParams))

		if len(matchCmd) != 0 || len(matchParams) != 0 {
			return false
		}
	}
	return true
}

// if !ok {
// 	w.Header().Set("Content-Type", "application/json")
// 	str := `{"status": "error", "description": "Please specify a target cluster"}`
// 	w.Write([]byte(str))
// 	LogWarning.Println("Cluster was not specified in request")
// 	return
// } else {
// 	LogInfo.Println("Setting cluster ", cluster)
// }
