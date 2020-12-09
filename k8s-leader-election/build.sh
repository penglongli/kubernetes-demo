#!/bin/bash

set -eo pipefail

# build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -v -x  -o leaderelection main.go

# build image
docker build -t leaderelection:1.0 .