#!/bin/bash

source scripts/env

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/binaries.sh protoc-gen-go"
  exit 0;
fi

# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-aarch_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-osx-aarch_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-osx-x86_64.zip
function protoc() {
  GOHOSTOS=$(go env GOHOSTOS)
  GOARCH=$(go env GOARCH)
  GOPATH=$(go env GOPATH)

  TARGET_OS=$GOHOSTOS
  TARGET_ARCH=$GOARCH

  if test "$GOHOSTOS" == "darwin"; then
    TARGET_OS="osx"
  fi

  if test "$GOARCH" == "arm64"; then
    TARGET_ARCH="aarch_64"
  elif test "$GOARCH" == "amd64"; then
    TARGET_ARCH="x86_64"
  fi

  echo "https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip"

  cd /tmp
  curl -L -O "https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip"

  unzip protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip
  mv bin/protoc "$GOPATH/bin/"
  mv include/google /usr/local/include/

  rm -f protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip
  rm -f readme.txt
  rmdir bin/
  rmdir include/
}

function protoc-gen-go-grpc() {
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
}

function protoc-gen-go() {
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
}

function protoc-gen-grpc-gateway() {
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
}

function protoc-gen-openapiv2() {
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
}

FILE=$(which $1)
if test -f "$FILE"; then
  echo "the binary already exists at: "$FILE""
else
  $1
  echo "download complete, this will place binaries in your \$GOBIN, make sure that your \$GOBIN is in your \$PATH."
fi
