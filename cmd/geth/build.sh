#!/usr/bin/env bash
CMT=`git log --pretty=format:'%H' -n 1`
godep go build -a -ldflags="-X main.gitCommit=$CMT" -o bw2bc
