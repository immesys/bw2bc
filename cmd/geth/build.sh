#!/usr/bin/env bash
CMT=`git log --pretty=format:'%h' -n 1`
go build -a -ldflags="-X main.gitCommit=$CMT"
