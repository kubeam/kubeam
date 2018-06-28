package server

import (
	//"reflect"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

/*UpdateResources updates the kubernetes resources using the new configuration
from the template list*/
func UpdateResources(vars map[string]string, templateList []string) {
	BuildResources(vars, templateList, true)
}

/*CreateResources creates a list of kubernetes resources specified in the
template list*/
func CreateResources(vars map[string]string, templateList []string) {
	BuildResources(vars, templateList, false)
}

/*BuildResources get Status of specified resources. Iterating over each type of app*/
func BuildResources(vars map[string]string, templateList []string, isReplace bool) {

	var kubeAction string
	if isReplace {
		kubeAction = "replace"
	} else {
		kubeAction = "create"
	}
	LogInfo.Println("CreateResources is set to =", kubeAction)
	// convert vars to something compatible with render_template
	m := make(map[string]interface{})
	for k, v := range vars {
		m[k] = v
	}

	cmdName := "./kubectl"

	for _, templateFile := range templateList {
		rendered := []byte(RenderTemplate(templateFile, m))

		tmpfile, err := ioutil.TempFile("tmp/", fmt.Sprintf("%s.rendered.", path.Base(templateFile)))
		if err != nil {
			LogInfo.Println(err)
		}
		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(rendered); err != nil {
			LogError.Println(err)
		}

		cmdArgs := []string{kubeAction, "-f", tmpfile.Name()}
		LogDebug.Println("Running: ", cmdName, " ", fmt.Sprintln(strings.Join(cmdArgs, " ")))
		cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
		if err != nil {
			LogError.Println(fmt.Sprintf("Running kubectl %s", cmdArgs))
			//
			// Resource might not exist. Lest try creating it.
			if isReplace {
				cmdArgs := []string{"create", "-f", tmpfile.Name()}
				LogDebug.Println("Running: ", cmdName, " ", fmt.Sprintln(strings.Join(cmdArgs, " ")))
				cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
				if err != nil {
					LogError.Println(fmt.Sprintf("2nd Try running kubectl %s", cmdArgs))
				} else {
					LogDebug.Println(string(cmdOut))
				}
			}
		} else {
			LogDebug.Println(string(cmdOut))
		}
		if err := tmpfile.Close(); err != nil {
			LogError.Println(err)
		}
		time.Sleep(2000 * time.Millisecond)
	}

}
