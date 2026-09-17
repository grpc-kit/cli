#!/bin/bash

source scripts/env

function gen_protoc() {
  # 生成 *.pb.go 文件
  protoc \
    -I ./ \
    -I /usr/local/include/ \
    -I "${GOPATH}"/src \
    -I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/ \
    -I "${GOPATH}"/src/github.com/googleapis/googleapis/ \
    --go_opt paths=source_relative \
    --go_out ./ \
    --go-grpc_opt paths=source_relative \
    --go-grpc_opt require_unimplemented_servers=false \
    --go-grpc_out ./ \
    ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/*.proto

  # 生成 *.pb.gw.go 与 swagger 接口文档
  protoc \
    -I ./ \
    -I /usr/local/include/ \
    -I "${GOPATH}"/src \
    -I "${GOPATH}"/src/github.com/googleapis/googleapis/ \
    -I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/ \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt grpc_api_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.gateway.yaml \
    --grpc-gateway_out ./ \
    --openapiv2_opt disable_default_errors=true \
    --openapiv2_opt disable_service_tags=true \
    --openapiv2_opt grpc_api_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.gateway.yaml \
    --openapiv2_opt openapi_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml \
    --openapiv2_out=json_names_for_fields=false:./ \
    ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.proto

  # 移动生成的 microservice.swagger.json 文件
  if test -f ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.swagger.json; then
    mv ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.swagger.json ./public/openapi/
  fi
}

function gen_ent() {
  # 因 "go generate" 会继承 "GOOS" 与 "GOARCH" 变量，如在交叉编译环境则可能无法运行
  GOOS=${GOHOSTOS} GOARCH=${GOHOSTARCH} go generate ./modeler/ent/
}

if test -z $1; then
  gen_protoc
else
  $1
fi
