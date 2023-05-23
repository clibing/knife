#########################################################################
# File Name: build.sh
# Author: clibing
# mail: wmsjhappy@@gmail.com
# Created Time: 二  5/23 15:30:32 2023
#########################################################################
#!/bin/bash

# cross build on macos error, will build on docker linux/amd64.

echo 'run this shell, need execute \"cd /go/src/github.com/clibing/knife/ && make\"'

docker run --rm -it -v `pwd`/:/go/src/github.com/clibing/knife/ clibing/golang:1.20 bash


