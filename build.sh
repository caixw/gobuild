#!/bin/sh

# SPDX-License-Identifier: MIT

cd `dirname $0`
date=`date -u '+%Y%m%d'`
hash=`git rev-parse HEAD`
path=github.com/caixw/gobuild/internal/cmd
go build -ldflags "-X ${path}.metadata=${date}.${hash}" -v -o ./cmd/gobuild/gobuild ./cmd/gobuild
