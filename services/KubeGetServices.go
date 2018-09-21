package services

import (
	"fmt"
	"strings"

	"github.com/kubeam/kubeam/common"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

/*KubeGetServices gets details of services for the given filter*/
func KubeGetServices(filter string) (map[string]interface{}, error) {
	resources := map[string]interface{}{}
	// creates the in-cluster config
	currconfig, err := rest.InClusterConfig()
	if err != nil {
		common.LogError.Println(err.Error())
		return resources, err
	}
	clientset, err := kubernetes.NewForConfig(currconfig)
	if err != nil {
		common.LogError.Println("Error getting clientset of kubernetes")
		return resources, err
	}

	name := "default"
	services, err := clientset.Core().Services(name).List(metav1.ListOptions{})
	if err != nil {
		common.LogError.Printf("Get service from kubernetes cluster error:%v", err)
		return resources, err
	}

	for _, service := range services.Items {
		if name == "default" && service.GetName() == "kubernetes" {
			continue
		}
		common.LogInfo.Println("namespace", name, "serviceName:", service.GetName(), "serviceKind:", service.Kind, "serviceLabels:", service.GetLabels(), service.Spec.Ports, "serviceSelector:", service.Spec.Selector)

		serviceEndpoints := GetExternalEndpoints(&service)

		common.LogInfo.Printf("serviceEndpoints : %v ", serviceEndpoints)

		//          // labels.Parser
		//          //set := labels.Set(service.Spec.Selector)
		//          selector_s, err_s := metav1.LabelSelectorAsSelector( service.Spec.Selector )
		//          if err_s != nil {
		//              return resources, fmt.Errorf( "invalid label selector: %v", err_s )
		//          }

	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		common.LogError.Println(err.Error())
		return resources, err
	}
	resourceName := filter
	common.LogInfo.Println("Resource name ", resourceName)
	//isFirst := true
	for _, d := range list.Items {
		// If resource matches our fileter. ==0 Insures we match from the beggining of string
		if strings.Index(d.Name, resourceName) == 0 {
			podout := map[string]interface{}{}

			podout["name"] = d.Name
			podout["replicas"] = d.Spec.Replicas

			//set := labels.Set(d.Spec.Selector)
			selector, err := metav1.LabelSelectorAsSelector(d.Spec.Selector)
			if err != nil {
				return resources, fmt.Errorf("invalid label selector: %v", err)
			}

			//if pods, err := clientset.Core().Pods(name).List(metav1.ListOptions{LabelSelector: set.AsSelector()}); err != nil {
			if pods, err := clientset.Core().Pods(name).List(metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
				//common.LogError.Printf("List Pods of service[%s] error:%v", service.GetName(), err)
				common.LogError.Printf("Errog getting Pods of deployment [%v] error:%v", d.Name, err)
			} else {
				podList := map[string]interface{}{}
				for _, v := range pods.Items {
					pod := map[string]interface{}{}
					pod["name"] = v.GetName()
					pod["node"] = v.Spec.NodeName
					pod["container"] = map[string]interface{}{
						"name":  v.Spec.Containers[0].Name,
						"image": v.Spec.Containers[0].Image,
					}

					podList[v.GetName()] = pod

					common.LogInfo.Printf("POD : %v\nNode: %v\nContainer Name: %v\nImage :%v\n",
						v.GetName(), v.Spec.NodeName, v.Spec.Containers[0].Name,
						v.Spec.Containers[0].Image)
				}
				podout["pods"] = podList
			}
			resources[d.Name] = podout

		}
	}
	return resources, nil
}
