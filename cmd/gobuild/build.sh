#!/bin/sh

# SPDX-FileCopyrightText: 2015-2024 caixw
#
# SPDX-License-Identifier: MIT

cd `dirname $0`
date=`date -u '+%Y%m%d'`
hash=`git rev-parse HEAD`
go build -ldflags "-X main.metadata=${date}.${hash}" -v -o ./gobuild ./
