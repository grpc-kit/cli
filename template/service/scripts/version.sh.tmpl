#!/bin/bash

source scripts/env

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/version.sh prefix"
  echo "\t ./scripts/version.sh release"
  echo "\t ./scripts/version.sh update"
  exit 0;
fi

function prefix() {
  TEMP=$(grep "version: \".*\"" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml)
  PREFIX_VERSION=$(echo -n $TEMP | awk -F"\"" '{ print $2 }')
  echo $PREFIX_VERSION
}

function release() {
  TEMP=$(cat VERSION)
  RELEASE_VERSION=$TEMP

  if test -z $RELEASE_VERSION; then
    RELEASE_VERSION=$(git describe --tags --dirty --always 2>/dev/null)
  fi

  echo $RELEASE_VERSION
}

function update() {
  GOHOSTOS=$(go env GOHOSTOS)

  PREFIX_VERSION=$(prefix)
  RELEASE_VERSION=$(release)

  if test $PREFIX_VERSION == $RELEASE_VERSION; then
    return
  fi

  if test ${GOHOSTOS} = "darwin"; then
    sed -i "" "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml
  else
    # fix run in container
    # sed: couldn't open temporary file sed1DDoX9: Permission denied
    #sed -i "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml
    cp api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml /tmp/microservice.openapiv2.yaml
    sed -i "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" /tmp/microservice.openapiv2.yaml
    mv /tmp/microservice.openapiv2.yaml api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml > /dev/null 2>&1
  fi
}

$1
