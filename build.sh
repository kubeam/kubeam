#!/usr/bin/env bash
##
## Simple build script for Go vendoring dependencies
## By Luis E Limon <luise_limon@intuit.com>


if [ "$1" == "" ]; then
    APPLICATIONS="kubeam isakonf"
else
    APPLICATIONS=$*
fi

REGISTRY="your registry"



##
## create  self signed cert. We don't want to provide real information on cert signature
if [ ! -f server.key -o ! -f server.crt ]; then
  openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=selfsigned.com" \
    -keyout server.key  -out server.crt
fi


BUILD_TARGET_DIR="`pwd`/target"
#APP_NAME=$( basename `pwd` | tr '[:upper:]' '[:lower:]' )

echo Building : $APP_NAME

#CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -a -o ${APP_NAME} kubeam/. && \
export GOPATH=$BUILD_TARGET_DIR
mkdir -p $GOPATH/src && \
mkdir -p $GOPATH/bin && \
mkdir -p $GOPATH/pkg && \
GOBIN=$GOPATH/bin && \
( [  -L  $GOPATH/src/kubeam ]  && rm $GOPATH/src/kubeam) 
( [  -L  $GOPATH/src/isakonf ]  && rm $GOPATH/src/isakonf) 

curl -L -o kubectl.linux https://storage.googleapis.com/kubernetes-release/release/v1.7.0/bin/linux/amd64/kubectl

CURR_DIR=`pwd`
#for a in kubeam isakonf; do
for a in $APPLICATIONS; do
    APP_NAME=$a

    cd $CURR_DIR
    echo $GOPATH && \
    ln -vFs `pwd`/src/$APP_NAME $GOPATH/src/$APP_NAME && \
    cd $GOPATH/src/$APP_NAME && \
    echo $GOPATH && \
    go get &&  \
    CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -a -o ${APP_NAME} $APP_NAME/. && \
    echo "Success your binary is in $GOPATH/bin" && \
    ls -al $GOPATH/bin
    if [ $? -ne 0 ]; then
        exit 1
    fi

    ## Run test of isakonf
    if [ "$a" == "isakonf" ]; then
        cd $GOPATH/src/$a/test
        mkdir -p rendered
        $GOPATH/bin/isakonf -v test-templates.yaml
    fi

done
#-installsuffix cgo


