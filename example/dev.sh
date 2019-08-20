#!/bin/sh
# @project: dwlog
# @file: dev.sh
# @author: dungw3b
# @date: 2019-08-12

./build.sh
if [ $? -eq 0 ]; then
  ./bin/dwlog-linux-amd64 -c conf/dwlog.json
fi
