#!/bin/sh

# SPDX-License-Identifier: MIT

cd `dirname $0`
builddate=`date -u '+%Y%m%d'`
commithash=`git rev-parse HEAD`
path=github.com/caixw/gobuild/internal/cmd
go build -ldflags "-X ${path}.buildDate=${builddate}  -X ${path}.commitHash=${commithash}" -v -o ./cmd/gobuild/gobuild ./cmd/gobuild
