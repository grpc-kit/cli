#!/bin/bash

# 引入全局静态变量
source scripts/env

# 引入全局动态变量
source scripts/variable.sh

# 需生成的模版类型
if test -z "${TEMPLATES}"; then
  TEMPLATES=dockerfile
fi

GOHOSTOS=$(go env GOHOSTOS)

# 解决 linux 与 darwin 的 sed 存在的差异
if test ${GOHOSTOS} = "darwin"; then
  SED="sed -i ''"
else
  SED="sed -i"
fi

# 编译与部署环境
if test -z "${BUILD_ENV}"; then
  BUILD_ENV=local
fi
if test -z "${DEPLOY_ENV}"; then
  DEPLOY_ENV=dev
fi
if test -z "${KUBERNETES_NAMESPACE}"; then
  if test -n "${CI_BIZ_GROUP_APPID}"; then
    KUBERNETES_NAMESPACE=biz-${DEPLOY_ENV}-${CI_BIZ_GROUP_APPID}
  fi
fi

function clean() {
  rm -rf Dockerfile
  rm -rf deploy/systemd/
  rm -rf deploy/supervisor/
  rm -rf deploy/kubernetes/${DEPLOY_ENV}
}

function systemd() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/systemd
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
  fi

  cp -rf scripts/templates/systemd/microservice.service ${TEMPLATE_PATH}/${APPNAME}.service

  eval "$SED" "s#_SERVICE_CODE_#${SERVICE_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_PRODUCT_CODE_#${PRODUCT_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_SHORT_NAME_#${SHORT_NAME}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_API_VERSION_#${API_VERSION}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${APPNAME}.service
}

function supervisor() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/supervisor
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
  fi

  cp -rf scripts/templates/supervisor/microservice.conf ${TEMPLATE_PATH}/${APPNAME}.conf

  eval "$SED" "s#_SERVICE_CODE_#${SERVICE_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_PRODUCT_CODE_#${PRODUCT_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_SHORT_NAME_#${SHORT_NAME}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_API_VERSION_#${API_VERSION}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
}

function kubernetes() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/kubernetes/${DEPLOY_ENV}
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
  fi

  cp -rf scripts/templates/kubernetes/* ${TEMPLATE_PATH}
  mv ${TEMPLATE_PATH}/workloads/deployments/microservice.yaml ${TEMPLATE_PATH}/workloads/deployments/${APPNAME}.yaml

  if test -f config/app-${DEPLOY_ENV}-${BUILD_ENV}.yaml; then
    cp -a config/app-${DEPLOY_ENV}-${BUILD_ENV}.yaml ${TEMPLATE_PATH}/config/configmap/app.yaml
  fi

  for FILE in kustomization.yaml service/ingresses.yaml service/services.yaml workloads/deployments/${APPNAME}.yaml
  do
    eval "${SED}" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_CI_REGISTRY_IMAGE_#${CI_REGISTRY_IMAGE}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_DOCKER_IMAGE_VERSION_#${DOCKER_IMAGE_VERSION}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_KUBERNETES_NAMESPACE_#${KUBERNETES_NAMESPACE}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_KUBERNETES_LABEL_PREFIX_#${KUBERNETES_LABEL_PREFIX}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_KUBERNETES_PM2_UUID_#${KUBERNETES_PM2_UUID}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_KUBERNETES_CLUSTER_DOMAIN_#${KUBERNETES_CLUSTER_DOMAIN}#g" ${TEMPLATE_PATH}/${FILE}
  done
}

function dockerfile() {
  # 如未设置父镜像，默认为 scratch
  if test -z ${DOCKER_IMAGE_FROM}; then
    DOCKER_IMAGE_FROM=scratch
  fi

  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)
  fi

  cp scripts/templates/Dockerfile ${TEMPLATE_PATH}

  if test ${GOHOSTOS} = "darwin"; then
    sed -i "" "s#_DOCKER_IMAGE_FROM_#${DOCKER_IMAGE_FROM}#g" ${TEMPLATE_PATH}/Dockerfile
  else
    sed -i "s#_DOCKER_IMAGE_FROM_#${DOCKER_IMAGE_FROM}#g" ${TEMPLATE_PATH}/Dockerfile
  fi
}

# 不做清理
#clean

# 避免运行无意义的指令
if test "${TEMPLATES}" = "dockerfile"; then
  dockerfile
elif test "${TEMPLATES}" = "kubernetes"; then
  kubernetes
elif test "${TEMPLATES}" = "systemd"; then
  systemd
elif test "${TEMPLATES}" = "supervisor"; then
  supervisor
fi
