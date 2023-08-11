#!/bin/sh

# SPDX-License-Identifier: MIT

cd `dirname $0`
date=`date -u '+%Y%m%d'`
hash=`git rev-parse HEAD`
go build -ldflags "-X main.metadata=${date}.${hash}" -v -o ./gobuild ./
