package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

/*GetDockerTag fetches the most recent tag of the resource defined in api yaml*/
func GetDockerTag(api string, m map[string]interface{}) (string, error) {
	/*
		Fetches the most recent docker tag of the qbo deployment
		Returns and error packet for any error
	*/
	var app, lookupresource, lookupnamespace string

	app = m["application"].(string)
	apiActions, err := GetAPIActions(api, app, m)
	ErrorHandler(err)

	for _, actionsitem := range apiActions {
		lookupresource = actionsitem["lookupresource"].(string)
		lookupnamespace = actionsitem["lookupresourcenamespace"].(string)
	}

	clientset := GetClientSet()
	client := clientset.AppsV1beta1().Deployments(lookupnamespace)
	dep, err := client.Get(lookupresource, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	tag := strings.Split(dep.Spec.Template.Spec.Containers[0].Image, ":")[1]
	return tag, nil
}

/*GetLogs retrieves the logs of the jobresource given in the request*/
func GetLogs(job *batchv1.Job) (string, error) {
	var logs []byte
	var err error
	clientset := GetClientSet()
	client := clientset.CoreV1().Pods(job.Namespace)

	if pods, err := client.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
	}); err == nil {
		for _, po := range pods.Items {
			// Match the correct Pod for the Job by verifying the creation timestamp
			if po.GetCreationTimestamp() == job.GetCreationTimestamp() {
				LogDebug.Printf("Found job pod: %s %s", po.Name, po.GetCreationTimestamp())
				logs, err = client.GetLogs(po.Name, &v1.PodLogOptions{}).Do().Raw()
				ErrorHandler(err)
				return string(logs), nil
			}
		}
	}
	return "", err
}

/*GetJobNamespace retrieves the namespace for the k8s Job from YAML*/
func GetJobNamespace(vars map[string]string) (string, error) {
	var namespace string
	m := make(map[string]interface{})

	for k, v := range vars {
		m[k] = v
	}
	ret, err := GetAPIActions("/v1/kubejob", vars["application"], m)
	if err != nil {
		return "", err
	}
	for _, actionsItem := range ret {
		_, ok := actionsItem["namespace"]
		if ok {
			namespace = actionsItem["namespace"].(string)
		}
	}
	return namespace, nil
}

/*RunJobActions ...*/
func RunJobActions(api string, vars map[string]interface{}) (map[string]interface{}, error) {
	app := vars["application"].(string)
	env := vars["environment"].(string)
	cluster := vars["cluster"].(string)
	actionsRetVal := make(map[string]interface{})
	actionsOutput := make(map[string][]interface{})

	ret, err := GetAPIActions(api, app, vars)
	if err != nil {
		LogDebug.Println("Error in call to GetApiActionsItem", api, app, vars)
		return nil, err
	}

	for _, actionsItem := range ret {
		LogDebug.Println("actionItem :", actionsItem)

		doAction := true
		if _, envDefined := actionsItem["environment"]; envDefined {
			doAction = IsClusterInList(env, cluster, actionsItem["environment"])
			LogDebug.Println("envDefined doAction :", doAction)
		}
		if doAction {
			LogDebug.Println("doAction :", doAction)
			var tempfileName string
			var tmpfile *os.File
			var currActionOutput actionOutput

			if _, ok := actionsItem["file"]; ok {
				LogDebug.Println("Creating rendered yaml ", fmt.Sprintf("applications/%v/%v", app, actionsItem["file"].(string)))
				rendered := []byte(RenderTemplate(fmt.Sprintf("applications/%v/%v", app, actionsItem["file"].(string)), vars))
				LogDebug.Println("Creating temp file", fmt.Sprintf("%s.rendered.", path.Base(actionsItem["file"].(string))))

				tmpfile, err = ioutil.TempFile("tmp/", fmt.Sprintf("%s.rendered.", path.Base(actionsItem["file"].(string))))
				if err != nil {
					LogInfo.Println(err)
				} else {
					tempfileName = tmpfile.Name()
					LogDebug.Println("Temp file name is: ", tempfileName)
				}
				defer os.Remove(tmpfile.Name())

				if _, err := tmpfile.Write(rendered); err != nil {
					LogError.Println(err)
				}
			}
			currActionOutput.Type = actionsItem["type"].(string)
			currActionOutput.Resource = actionsItem["resource"].(string)
			currActionOutput.Name = actionsItem["name"].(string)
			currActionOutput.Action = actionsItem["action"].(string)
			if jobobj, err := decodeKubernetesJobYAML(tempfileName); err == nil {
				currActionOutput.Log = createJob(jobobj, vars["namespace"].(string))
			} else {
				currActionOutput.Log = err.Error()
			}
			actionsOutput["actions"] = append(actionsOutput["actions"], currActionOutput)
		} else {
			LogInfo.Println("No action to do for this cluster")
		}
	}

	actionsRetVal["actions"] = actionsOutput["actions"]
	return actionsRetVal, nil
}

func decodeKubernetesJobYAML(tempfileName string) (*batchv1.Job, error) {
	LogDebug.Println("Decoding rendered temp file to kubejob object")
	yamlData, err := ioutil.ReadFile(tempfileName)
	if err == nil {
		decode := scheme.Codecs.UniversalDeserializer().Decode
		kubeobj, _, err := decode(yamlData, nil, nil)
		if err == nil {
			return kubeobj.(*batchv1.Job), nil
		}
		LogError.Println(err.Error())
		return nil, err
	}
	LogError.Println(err.Error())
	return nil, err
}

func createJob(kubeobj *batchv1.Job, namespace string) string {
	deleteKubernetesJob(kubeobj.GetName(), namespace)
	return createKubernetesJob(kubeobj, namespace)
}

func createKubernetesJob(kubeobj *batchv1.Job, namespace string) string {
	LogDebug.Println("Creating : ", kubeobj.GetName())
	clientset := GetClientSet()
	client := clientset.BatchV1().Jobs(namespace)

	if _, err := client.Create(kubeobj); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Created: %s", kubeobj.GetName())
}

func deleteKubernetesJob(resource, namespace string) string {
	graceperiod := int64(0)
	clientset := GetClientSet()
	client := clientset.BatchV1().Jobs(namespace)
	LogDebug.Println("Deleting : ", resource)
	deletePolicy := metav1.DeletePropagationForeground
	if err := client.Delete(resource, &metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &graceperiod,
	}); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Deleted: %s", resource)
}

// GetJobAPI retrieves all keypairs api for the k8s Job from YAML
func GetJobAPI(vars map[string]string) (map[string]interface{}, error) {
	LogDebug.Println("Fetching Job API actions")

	m := make(map[string]interface{})

	for k, v := range vars {
		m[k] = v
	}

	ret, err := GetAPIActions("/v1/kubejob", vars["application"], m)
	ErrorHandler(err)

	for _, actionsItem := range ret {
		return actionsItem, nil
	}
	return make(map[string]interface{}), fmt.Errorf("No definition for /v1/kubejob application %v", vars["application"])
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
