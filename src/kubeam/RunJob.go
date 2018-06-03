package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

/*RunJob extracts command line parameters and creates a Job object to execute
the command. Returns response of execution back*/
func RunJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var data map[string]string
	var response string

	m := make(map[string]interface{})
	vars := mux.Vars(r)

	byteData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		LogError.Println(err.Error())
		response = fmt.Sprintf("Incorrect JSON format: %s", err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")

		cmd := data["cmd"]
		params := data["params"]
		for k, v := range vars {
			m[k] = v
		}
		m["jobcommand"] = cmd
		m["jobparams"] = params

		actionsOutput, _ := RunActions("/v1/runjob", m)
		outputjson, _ := json.MarshalIndent(actionsOutput, "", " ")
		w.Write(outputjson)
	}
	w.Write([]byte(response))
}

/*Security check to disable running multiple operations on command line
Not integrated yet. Would disable running multiple complex commands*/
func commandSanityCheck(cmd, params string) bool {
	LogInfo.Println(cmd, params)
	re, _ := regexp.Compile("[\\|;&$><`\\!]")

	if cmd != "" {
		matchCmd := re.FindString(cmd)
		matchParams := re.FindString(params)

		if len(matchCmd) != 0 || len(matchParams) != 0 {
			return false
		}
	}
	return true
}

func err(err error) {
	if err != nil {
		LogError.Println(err.Error())
	}
}
