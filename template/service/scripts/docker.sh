#!/bin/bash

# 引入全局静态变量
source scripts/env

# 引入全局动态变量
source scripts/variable.sh

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/docker.sh build"
  echo "\t ./scripts/docker.sh push"
fi

function build() {
  docker build -t ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION} ./
  echo "Now you can upload image: "docker push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}""
}

function push() {
  docker push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}
}

function run() {
  docker run -i -t --rm \
      -v $GOPATH/pkg:/go/pkg \
      -v $(pwd):/usr/local/src \
      -w /usr/local/src \
      --network host \
      ccr.ccs.tencentyun.com/grpc-kit/cli:${CLI_VERSION} \
      make run
}

$1
