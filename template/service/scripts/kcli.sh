#!/bin/bash

# 引入全局静态变量
source scripts/env

# 引入全局动态变量
source scripts/variable.sh

# 当前目录必须为 gitops 仓库根路径
function commit() {
  # 变量 CI_BIZ_CODE_BUILD 来自 CICD 系统
  if test "${CI_BIZ_CODE_BUILD}" != "true"; then
    return
  fi

  # fix: https://github.blog/2022-04-12-git-security-vulnerability-announced/
  git config --global --add safe.directory $(pwd)

  cd ${KUBERNETES_YAML_DIRECTORY}
  /bin/kustomize edit set image ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}

  git config user.name ${BUILD_USER}
  git config user.email ${BUILD_USER_EMAIL}

  git add *
  git commit -m "build: set image ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}"

  # 变量 CI_OPS_REPO_URL 来自 CICD 系统
  git remote add gitops ${CI_OPS_REPO_URL}
  git push gitops HEAD:refs/heads/main
}

# 当前目录必须为 gitops 仓库根路径
function apply() {
  cd ${KUBERNETES_YAML_DIRECTORY}
  /bin/kustomize build | /bin/kubectl apply -f -
}

$1
