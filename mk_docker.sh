bash build.sh kubeam
docker-compose up --build
# docker run -d -p 7000:3306 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=kubeam mariadb:10
# docker run -d -p 6000:6379 redis
# docker run -d -p 8443:443 kubeam

# curl -L -o kubectl.linux https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl

# Test health-check
#curl -I http://localhost:8081/health-check

## Test health-check ssl port
# curl -k -I https://admin:123456@localhost:8443/health-check

## Test create endpoint
# curl -k -I -X PUSH https://admin:123456@localhost:8443/v1/create/sample/qa/main/300

# docker run -p 8443:443 kubeam
# curl -k  -X GET https://admin:123456@localhost:8443/health-check
