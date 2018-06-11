
1)   bash minikube-setup.sh
2)   source set-registry.sh
3)   bash download-kubectl.sh
4)   cp config-sample.yaml config.yaml
5)   docker build -f Dockerfile-kubeam.dkr . -t localhost:5000/kubeam
6)   docker push localhost:5000/kubeam
7)   curl -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/health-check
8)   curl -X POST -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/create/kubeam/minikube/main/latest