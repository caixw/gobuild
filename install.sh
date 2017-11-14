#!/bin/sh
#
# Copyright 2017 by caixw, All rights reserved.
# Use of this source code is governed by a MIT
# license that can be found in the LICENSE file.

cd `dirname $0`
builddate=`date -u '+%Y%m%d'`
commithash=`git rev-parse HEAD`
go install -ldflags "-X main.buildDate=${builddate}  -X main.commitHash=${commithash}" -v
