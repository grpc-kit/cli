#!/bin/bash

# 确保GOPATH变量有设置
if test -z "${GOPATH}"; then
  echo "Please set the environment variable GOPATH before running make"
  exit 1
fi
