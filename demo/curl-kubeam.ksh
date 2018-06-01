#curl -k $(minikube service kubeamservice --format "https://{{.IP}}:{{.Port}}" --url)/health-check
curl -k -X GET $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/health-check
curl -k -X POST $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/create/kubeam/ci/main/latest
