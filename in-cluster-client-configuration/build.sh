#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -v -x -o in-cluster main.go

docker build -t in-cluster .
