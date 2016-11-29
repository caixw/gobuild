#!/bin/sh

cd `dirname $0`

go build -ldflags "-X main.buildDate=`date -u '+%Y%m%d'`  -X main.commitHash=`git rev-parse HEAD`" -v
