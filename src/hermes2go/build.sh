#!/bin/bash -x
git describe --always --tags --long > version.txt
go get gopkg.in/yaml.v2 
VERSION=$(cat version.txt) 
go build -v -ldflags "-X main.version=${VERSION}"