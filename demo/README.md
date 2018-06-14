
## Setup minikube with a local registry
`bash demo/minikube-setup.sh`  

## set minikube as your registry.
`source demo/set-registry.sh`  

## Build KubeAM, docker and push to minikube registry
`bash build.sh kubeam`  
`bash download-kubectl.sh`  
`cp config-sample.yaml config.yaml`  
`docker build -f Dockerfile-kubeam.dkr . -t localhost:5000/kubeam`  
`docker push localhost:5000/kubeam`  

## Install temporary bootstrap kubeam
`kubectl apply -f demo/kubeam-bootstrap-pod.yaml && kubectl apply -f demo/kubeam-bootstrap-service.yaml`

## Check if is running
`curl -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/health-check`

## Deploy full kubeAM stack
`curl -X POST -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/create/kubeam/minikube/main/latest`
## Sample app work flow sample

### Build Sample app
cd ~/kubeam-demo/springboot-jsp  
mvn package  
docker build . -t localhost:5000/springboot-jsp:001  
docker push localhost:5000/springboot-jsp:001  

### Deploy springboot-jsp sample app (Once is build and docker available
`curl -X POST -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/create/springboot-jsp/dev/main/001`  

### Open deployed service
minikube --namespace dev-springboot-jsp service dev-springboot-jsp-service

### build new version and deploy
(In springboot-jsp repo)
mvn package  
docker build . -t localhost:5000/springboot-jsp:002
docker push localhost:5000/springboot-jsp:002
`curl -X POST -k $(minikube service kubeamservice --format "https://admin:123456@{{.IP}}:{{.Port}}" --url)/v1/deploy/springboot-jsp/dev/main/002`  

### Test your application
`curl $(minikube --namespace dev-springboot-jsp service dev-springboot-jsp-service --url)`  



