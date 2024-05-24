#!/bin/bash

# 生成的镜像相关
# DOCKER_IMAGE_FROM=scratch
RELEASE_VERSION=$(cat VERSION)
if test -z "${DOCKER_IMAGE_VERSION}"; then
  if test -z "${BUILD_ID}"; then
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}
  else
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
  fi
fi
if test -z "${CI_REGISTRY_IMAGE}"; then
  if test -z "${CI_REGISTRY_HOSTNAME}"; then
    CI_REGISTRY_HOSTNAME="docker.io"
  fi
  if test -z "${CI_REGISTRY_NAMESPACE}"; then
    CI_REGISTRY_NAMESPACE=${PRODUCT_CODE}
  fi
  CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}
fi

# Kubernetes 相关变量
if test -z "${KUBERNETES_NAMESPACE}"; then
  KUBERNETES_NAMESPACE=default
fi

# 启用构建相关的用户
if test -z "${BUILD_USER}"; then
  BUILD_USER=${USER}
fi
if test -z "${BUILD_USER_EMAIL}"; then
  BUILD_USER_EMAIL=${BUILD_USER}@$(hostname)
fi

# 部署编译环境相关
if test -z "${DEPLOY_ENV}"; then
  DEPLOY_ENV=dev
fi
if test -z "${BUILD_ENV}"; then
  BUILD_ENV=local
fi

if test -z "${CI_BIZ_GROUP_APPID}"; then
  CI_BIZ_GROUP_APPID=${PRODUCT_CODE}
fi

# 如果存在以下各对应环境的文件，则覆盖以上所设置的同名变量
if test -f "scripts/env-${DEPLOY_ENV}-${BUILD_ENV}"; then
  source scripts/env-${DEPLOY_ENV}-${BUILD_ENV}
fi
