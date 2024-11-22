#!/bin/bash -x
git describe --always --tags --long > version.txt
go get gopkg.in/yaml.v3
VERSION=$(cat version.txt) 
GOOS=windows GOARCH=386 go build -v -ldflags "-X main.version=${VERSION}" -o hermes2go.exe