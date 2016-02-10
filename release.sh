#!/bin/sh

export GOOS=linux GOARCH=amd64
echo "Building for ${GOOS}-${GOARCH}"
go build
file ./git-remind 
tar cvzf "git-remind.${GOOS}-${GOARCH}.tar.gz" git-remind config.json

export GOOS=darwin GOARCH=amd64
echo "Building for ${GOOS}-${GOARCH}"
go build
file ./git-remind 
tar cvzf "git-remind.${GOOS}-${GOARCH}.tar.gz" git-remind config.json

export GOOS=windows GOARCH=amd64
echo "Building for ${GOOS}-${GOARCH}"
go build
file ./git-remind.exe 
zip "git-remind.${GOOS}-${GOARCH}.zip" git-remind.exe config.json
