#!/bin/sh

# TODO; kaniko 镜像仅支持 /bin/sh 解析器

# 引入全局静态变量
source scripts/env

# 引入全局动态变量
source scripts/variable.sh

function build() {
  /kaniko/executor --dockerfile ./Dockerfile --context ./ --destination ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION} --log-format text --log-timestamp
}

build
