package server

import (
	"fmt"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

/*GetDockerTag fetches the most recent tag of the resource defined in api yaml*/
func GetDockerTag(api string, m map[string]interface{}) (map[string]string, int) {
	/*
		Fetches the most recent docker tag of the deployment
		Returns and error packet for any error
	*/
	var app, lookupresource, lookupnamespace, resource string
	var code int
	response := make(map[string]string)

	app = m["application"].(string)
	apiActions, err := GetAPIActions(api, app, m)
	ErrorHandler(err)

	for _, actionsitem := range apiActions {
		resource = actionsitem["resource"].(string)
		lookupresource = actionsitem["lookupresource"].(string)
		lookupnamespace = actionsitem["lookupresourcenamespace"].(string)
	}

	clientset := GetClientSet()
	client := clientset.AppsV1beta1().Deployments(lookupnamespace)
	dep, err := client.Get(lookupresource, metav1.GetOptions{})

	if err != nil {
		starttime := strings.Replace(time.Now().String(), "T", " ", 1)[:20]
		jobid := fmt.Sprintf("%s-%s", resource, starttime)
		response["JobId"] = jobid
		response["StartTime"] = starttime
		response["JobName"] = resource
		response["LastProbeTime"] = starttime
		response["JobStatus"] = "Failed"
		if kerr.IsNotFound(err) {
			code = 404
		} else if kerr.IsUnauthorized(err) {
			code = 401
		}
		response["Logs"] = err.Error()
	} else {
		response["Tag"] = strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")[1]
		LogInfo.Printf("Tag for %s -> %s", resource, response)
		code = 0
	}
	return response, code
}

/*GetLogs retrieves the logs of the jobresource given in the request*/
func GetLogs(job *batchv1.Job) string {
	var logs []byte
	clientset := GetClientSet()
	listoptions := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
	}
	pods, err := clientset.CoreV1().Pods(job.Namespace).List(listoptions)
	ErrorHandler(err)

	for _, po := range pods.Items {
		LogInfo.Println(po.Name, po.GetCreationTimestamp())

		// Match the correct Pod for the Job by verifying the creation timestamp
		if po.GetCreationTimestamp() == job.GetCreationTimestamp() {
			logs, err = clientset.CoreV1().Pods(job.Namespace).GetLogs(po.Name, &v1.PodLogOptions{}).Do().Raw()
			ErrorHandler(err)
			// LogInfo.Println(string(logs))
		}
	}
	return string(logs)
}

/*GetJobNamespace retrieves the namespace for the k8s Job from YAML*/
func GetJobNamespace(vars map[string]string) string {
	var namespace string
	m := make(map[string]interface{})

	for k, v := range vars {
		m[k] = v
	}

	ret, err := GetAPIActions("/v1/kubejob", vars["application"], m)
	ErrorHandler(err)

	for _, actionsItem := range ret {
		_, ok := actionsItem["namespace"]
		if ok {
			namespace = actionsItem["namespace"].(string)
		}
	}
	return namespace
}

/*GetClientSet returns a clientset object to make API calls*/
func GetClientSet() *kube.Clientset {

	kubeconfig, err := config.GetString("/kube/config", "")
	ErrorHandler(err)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// config, err := rest.InClusterConfig()
	ErrorHandler(err)

	clientset, err := kube.NewForConfig(config)
	ErrorHandler(err)
	return clientset
}
