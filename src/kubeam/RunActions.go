package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type ApiList struct {
	Description string
	Application map[string][]map[string]interface{}
}

func Keys(m map[string][]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func IsClusterInList(environment, cluster, m interface{}) bool {
	ms := make(map[string]interface{})

	for k, v := range m.(map[interface{}]interface{}) {
		if k == environment {
			ms[k.(string)] = v

			for _, cv := range v.([]interface{}) {
				if cv.(string) == cluster {
					return true
				}

			}

		}
	}
	return false
}

/*
* We load Api definitions from file every time work needs to be done.
* A future enhancement will be to preload this so that we don't need to keep ready from disk
 */
func GetApiActions(api string, application string, m map[string]interface{}) ([]map[string]interface{}, error) {
	var myapi ApiList

	fmt.Printf("Load Api for application %v\n", application)
	// yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%v.yaml", application))
	// if err != nil {
	// 	return nil, errors.New(fmt.Sprintf("Could not find a cluster definition for api %v ", application))
	// }
	rendered := []byte(render_template(fmt.Sprintf("applications/%v/api.yaml", application), m))
	err := yaml.Unmarshal(rendered, &myapi)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse %v", application))
	}

	for kubeamApi, actionsMap := range myapi.Application {
		if strings.ToLower(kubeamApi) == api {
			return actionsMap, nil
		}

	}
	return nil, errors.New(fmt.Sprintf("Api %v not found for application %v", api, application))
}

func ExecCmd(cmdName string, cmdArgs []string) ([]byte, error) {
	LogDebug.Println("Running: ", cmdName, fmt.Sprintln(strings.Join(cmdArgs, " ")))
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if cmdOut != nil {
		LogInfo.Println("CmdOut: ", string(cmdOut))
	}
	if err != nil {
		LogError.Println("Running: ", cmdName, fmt.Sprintln(strings.Join(cmdArgs, " ")))
		LogError.Println(err.Error())
		return nil, err
	} else {
		return cmdOut, err
	}

}

type actionOutput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Log      string `json:"Log"`
}

func RunActions(api string, vars map[string]interface{}) (map[string]interface{}, error) {
	app := vars["application"].(string)
	env := vars["environment"].(string)
	cluster := vars["cluster"].(string)

	ret, err := GetApiActions(api, app, vars)
	if err != nil {
		return nil, err
	}

	CmdOutProcessor := func(x []byte) string {
		var ret string

		ret = strings.Replace(string(x), "\"", "", -1)
		ret = strings.Replace(string(ret), "\n", "", -1)
		return ret
	}

	actionsRetVal := make(map[string]interface{})
	actionsOutput := make(map[string][]interface{})

	for _, actionsItem := range ret {
		var namespace []string

		_, ok := actionsItem["namespace"]
		if ok {
			namespace = []string{"--namespace", actionsItem["namespace"].(string)}

		}
		doAction := false
		_, envDefined := actionsItem["environment"]
		if envDefined {
			doAction = IsClusterInList(env, cluster, actionsItem["environment"])
		} else {
			doAction = true
		}
		if doAction {
			cmdName := "./kubectl"
			var tempfileName string
			var tmpfile *os.File
			_, ok := actionsItem["file"]
			if ok {
				rendered := []byte(render_template(fmt.Sprintf("applications/%v/%v", app, actionsItem["file"].(string)), vars))

				tmpfile, err = ioutil.TempFile("tmp/", fmt.Sprintf("%s.rendered.", path.Base(actionsItem["file"].(string))))
				if err != nil {
					LogInfo.Println(err)
				} else {
					tempfileName = tmpfile.Name()
				}
				//defer os.Remove(tmpfile.Name()) // clean up

				if _, err := tmpfile.Write(rendered); err != nil {
					LogError.Println(err)
				}

			}
			var currActionOutput actionOutput
			currActionOutput.Type = actionsItem["type"].(string)
			currActionOutput.Resource = actionsItem["resource"].(string)
			currActionOutput.Name = actionsItem["name"].(string)
			currActionOutput.Action = actionsItem["action"].(string)

			if actionsItem["action"] == "create" || actionsItem["action"] == "replace" || actionsItem["action"] == "apply" {
				fmt.Println(tempfileName)
				cmdArgs := []string{actionsItem["action"].(string), "-f", tempfileName}
				cmdOut, _ := ExecCmd(cmdName, append(namespace, cmdArgs...))
				currActionOutput.Log = CmdOutProcessor(cmdOut)

			} else if actionsItem["action"] == "recreate" {
				var cmdArgs []string

				cmdArgs = []string{"delete", fmt.Sprintf("%v/%v", actionsItem["type"], actionsItem["resource"])}
				cmdOut, _ := ExecCmd(cmdName, append(namespace, cmdArgs...))
				currActionOutput.Log = CmdOutProcessor(cmdOut)

				time.Sleep(3000 * time.Millisecond)
				cmdArgs = []string{"create", "-f", tempfileName}
				cmdOut, _ = ExecCmd(cmdName, append(namespace, cmdArgs...))
				currActionOutput.Log = fmt.Sprintf("%v, %v", currActionOutput.Log, CmdOutProcessor(cmdOut))

			}
			if tmpfile != nil {
				if err := tmpfile.Close(); err != nil {
					LogWarning.Println(os.Stderr, err)
				}
			}
			actionsOutput["actions"] = append(actionsOutput["actions"], currActionOutput)
		}
	}

	actionsRetVal["actions"] = actionsOutput["actions"]

	return actionsRetVal, nil

}
