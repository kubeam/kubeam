
## setup Minikube with a local registry
`cd demo;  bash minikube-setup.sh ;cd ..
source demo/set-registry.sh`

## Build KubeAM, docker and push to minikube registry
`bash download-kubectl.sh`
`cp config-sample.yaml config.yaml`
`docker build -f Dockerfile-kubeam.dkr . -t localhost:5000/kubeam`j
`docker push localhost:5000/kubeam`

## Install temporary bootstrap kubeam
`kubectl apply -f demo/kubeam-bootstrap-pod.yaml && kubectl apply -f demo/kubeam-bootstrap-service.yaml`

## Check if is running
`curl -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/health-check`

## Deploy full kubeAM stack
`curl -X POST -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/create/kubeam/minikube/main/latest`
