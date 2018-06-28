#!/usr/bin/env bash
##
## Simple build script for Go vendoring dependencies
## By Luis E Limon <luise_limon@intuit.com>

BUILD_TARGET_DIR="`pwd`/target"
APP_NAME=$( basename `pwd` | tr '[:upper:]' '[:lower:]' )

echo Building : $APP_NAME

GOPATH=$BUILD_TARGET_DIR
mkdir -p $GOPATH/src && \
mkdir -p $GOPATH/bin && \
GOBIN=$GOPATH/bin && \
( [  -L  $GOPATH/src/$APP_NAME ]  && rm $GOPATH/src/$APP_NAME ) 

echo $GOPATH && \
ln -vFs `pwd`/src $GOPATH/src/$APP_NAME && \
cd $GOPATH/src/$APP_NAME && \
echo $GOPATH && \
go get &&  \
go build -a -installsuffix cgo -o $APP_NAME . && \
echo "Success your binary is in $GOPATH/bin" && \
ls -al $GOPATH/bin



