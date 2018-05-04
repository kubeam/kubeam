package main

import (
	//"reflect"
	"bytes"
	"encoding/json"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	//"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

////
// get Status of specified resources.
// Iteracting over each type of app
///
func GetResourceStatus(parms map[string]string, resourcePostfix []string) ([]byte, error) {
	application := parms["application"]
	appEnv := parms["environment"]
	cluster := parms["cluster"]

	// creates the in-cluster config
	currconfig, err := rest.InClusterConfig()
	if err != nil {
		LogError.Println(err.Error())
		return []byte("{}"), err
	}
	clientset, err := kubernetes.NewForConfig(currconfig)
	if err != nil {
		LogError.Println("Error getting clientset of kubernetes")
		return []byte("{}"), err
	}
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		LogError.Println(err.Error())
		return []byte("{}"), err
	}
	resourceName := fmt.Sprintf("%v-%v-c%v", application, appEnv, cluster)
	LogInfo.Println("Resource name ", resourceName)
	var output bytes.Buffer
	output.WriteString("{")
	output.WriteString(fmt.Sprintf("\"application\": \"%v\",\n", application))
	output.WriteString(fmt.Sprintf("\"environment\": \"%v\",\n", appEnv))
	output.WriteString(fmt.Sprintf("\"cluster\": \"%v\",\n", cluster))
	isFirst := true
	for _, d := range list.Items {
		// If resource matches our fileter. ==0 Insures we match from the beggining of string
		if strings.Index(d.Name, resourceName) == 0 {
			out := map[string]interface{}{}

			out["name"] = d.Name
			out["replicas"] = d.Spec.Replicas
			//out["strategy"] = d.Spec.strategy
			out["paused"] = d.Spec.Paused

			outputJSON, _ := json.Marshal(out)
			if isFirst == true {
				isFirst = false
			} else {
				output.WriteString(",")
			}
			output.WriteString(fmt.Sprintf("\"%v\" :", d.Name))
			output.WriteString(string(outputJSON))

		}
	}
	output.WriteString("}")

	out := map[string]interface{}{}
	json.Unmarshal(output.Bytes(), &out)
	outputJSON, _ := json.MarshalIndent(out, "", " ")
	LogInfo.Println("Output :", output.String())

	return outputJSON, err
}
