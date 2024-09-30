#!/bin/bash

# 引入全局静态变量
source scripts/env

# 引入全局动态变量
source scripts/variable.sh

# 检查是否有 docker 或 podman
if command -v docker &> /dev/null; then
  CONTAINER_ENGINE="docker"
elif command -v podman &> /dev/null; then
  CONTAINER_ENGINE="podman"
else
  echo "Neither Docker nor Podman is installed. Please install Docker or Podman."
  exit 1
fi

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/docker.sh build"
  echo "\t ./scripts/docker.sh push"
fi

function build() {
  $CONTAINER_ENGINE build -t ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION} ./
  echo "Now you can upload image: "${CONTAINER_ENGINE} push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}""
}

function push() {
  $CONTAINER_ENGINE push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}
}

function run() {
  $CONTAINER_ENGINE run -i -t --rm \
      -v $GOPATH/pkg:/go/pkg \
      -v $(pwd):/usr/local/src \
      -w /usr/local/src \
      --network host \
      ccr.ccs.tencentyun.com/grpc-kit/cli:${CLI_VERSION} \
      make run
}

$1
