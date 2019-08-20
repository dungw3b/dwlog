#!/usr/bin/env bash
# @project: dwlog
# @file: build.sh
# @author: dungw3b
# @date: 2019-08-12

GOPATH=/srv/dwlog
GOBIN=$GOPATH/bin

rm -f bin/*
env GOPATH=$GOPATH go clean

platforms=("linux/amd64")
 
for platform in "${platforms[@]}"
do
  split=(${platform//\// })
  GOOS=${split[0]}
  GOARCH=${split[1]}
  output=dwlog'-'$GOOS'-'$GOARCH
  if [ $GOOS = "windows" ]; then
    output+=".exe"
  fi

  env GOPATH=$GOPATH GOOS=$GOOS GOARCH=$GOARCH go build -v -o $GOBIN/$output server.go
  if [ $? -ne 0 ]; then
    echo -e "\033[31mBuild error! Aborting the script execution...\033[0m"
    exit 1
  fi
done

