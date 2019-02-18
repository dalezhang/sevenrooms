#!/bin/bash

set -e -x
cd $PROJECT_DIR

echo "List was in the current directory."
ls -lat
go env
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  $PROJECT_DIR/dist/main
retVal=$?
if [ ! $? -eq 0 ]; then
    echo "Go Build Error"
fi