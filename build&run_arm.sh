#!/bin/bash

cwd=$(dirname $(pwd))

#env GOPATH=$cwd go get "github.com/stianeikeland/go-rpio"

if [ "$#" -ne 1 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
 
echo compiling $1
env GOOS=linux GOARCH=arm go build $1

if [ $? -ne 0 ]; then
    exit $?
fi

echo copying...
scp main root@172.22.11.8:/tmp/watchmen

#echo executing...
#ssh root@172.22.11.8 /tmp/watchmen

