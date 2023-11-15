# CHANGELOG-0.3

| 名称         | 说明                     |
|------------|------------------------|
| Added      | 添加新功能                  |
| Changed    | 功能的变更                  |
| Deprecated | 未来会删除                  |
| Removed    | 之前为Deprecated状态，此版本被移除 |
| Fixed      | 功能的修复                  |
| Security   | 有关安全问题的修复              |

## [Unreleased]

## [0.3.4] - 2023-11-15

### Added

#### 新增客户端 health 验证示例函数

生成服务的模版地址 `cmd/client/health.go`。

```shell
# go run cmd/client/health.go

grpc_health_v1 check ok
grpc health private check ok
```

该示例探测以下两个方法：

1. 标准的 [grpc_health_v1](https://pkg.go.dev/google.golang.org/grpc/health/grpc_health_v1) Check 与 Watch 方法；
2. 自定义的 `HealthCheck` 方法：

```text
rpc HealthCheck(grpc_kit.api.known.status.v1.HealthCheckRequest) returns (grpc_kit.api.known.status.v1.HealthCheckResponse) {}
```

#### 对类库 google.golang.org/grpc 版本不在锁定

原先的版本锁定为 v1.38.0 因为被 etcd 依赖，此次调整更新为最新 v1.59.0 版本，移除 go.mod 中 replace 语句。

```shell
replace google.golang.org/grpc => google.golang.org/grpc v1.38.0
```

### Fixed

#### 修复 rpc_grpc_status_code 状态码获取不对

```text
rpc_server_duration_milliseconds_bucket{rpc_grpc_status_code="0",rpc_method="Demo", ... ,rpc_system="grpc",le="0"} 1
rpc_server_duration_milliseconds_bucket{rpc_grpc_status_code="0",rpc_method="Demo", ... ,rpc_system="grpc",le="10"} 2
```

这个为 [otelgrpc 在 v0.45.0 版本](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/9d4eb7e7706038b07d33f83f76afbe13f53d171d/instrumentation/google.golang.org/grpc/otelgrpc/interceptor.go#L371C69-L371C69) 之前存在 BUG 需升级版本，`statusCode` 未成功赋值。

```text
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.0
```

在新版本中已经修复。

## [0.3.3] - 2023-11-09

### Added

#### 配置对接兼容 S3 接口的公有云对象存储

1. 基于 minio-go/v7 封装，提供简单：Get、Upload、Delete、Attributes、Iter 等方法；
2. 实现腾讯云 COS、阿里云 OSS的接入，见[配置使用文档](https://grpc-kit.com/docs/spec-cfg/objstore/)。

#### 可观测性重构以支持 OTLP 协议上报数据

1. 支持 遥测数据（链路、指标）通过 OLTP 协议上报数据，见 [配置使用文档](https://grpc-kit.com/docs/spec-cfg/observables/)；
2. 支持 阿里云 [可观测链路 OpenTelemetry 版](https://grpc-kit.com/docs/spec-cfg/observables/#%E5%AF%B9%E6%8E%A5%E9%98%BF%E9%87%8C%E4%BA%91%E6%9C%8D%E5%8A%A1) 的接入；
3. 支持 腾讯云 [应用性能监控](https://grpc-kit.com/docs/spec-cfg/observables/#%E5%AF%B9%E6%8E%A5%E8%85%BE%E8%AE%AF%E4%BA%91%E6%9C%8D%E5%8A%A1) 的接入；
4. 支持 [私有 jaeger 服务](https://grpc-kit.com/docs/spec-cfg/observables/#%E5%AF%B9%E6%8E%A5%E7%A7%81%E6%9C%89-jaeger-%E6%9C%8D%E5%8A%A1) 的接入。

#### 环境文件 env-* 新增变量 DOCKER_IMAGE_FROM

```shell
scripts/env-dev-local

# 基础镜像：构建业务镜像依赖的基础环境
DOCKER_IMAGE_FROM=centos:latest
```

也就是控制 Dockerfile 中 FROM 的镜像来源。

#### 在 jwt 中添加属性 `tenant` 表示租户

示例 token 格式如下：

```shell
{
  "aud": "api-gateway",
  "exp": 1893427200,
  "iat": 1668396542,
  "iss": "https://grpc-kit.com/oauth2",
  "sub": "oneops",
  "email": "oneops@grpc-kit.com",
  "email_verified": true,
  "federated_claims": {
    "connector_id": "local",
    "user_id": "oneops"
  },
  "groups": [
    "admin"
  ],
  "tenant" : "default"
}
```

相比原先新增 `tenant` 表示租户，参考了[第三方文档设计](https://userfront.com/dashboard/jwt)。

#### 支持 http 响应体为空时以 204 状态码返回

当微服务中 rpc 方法定义使用 `google.protobuf.Empty` 类型返回时，处理请求时会在 http gateway 中判断 proto 类型是否为 `*emptypb.Empty` 如则以状态码 204 返回。

```golang
// 该微服务支持的 RPC 方法定义
service OpsaidTest1 {
  rpc HelloNoContent(DemoRequest) returns (google.protobuf.Empty) {}
}
```

### Fixed

#### make manifests 变量问题修复

```shell
scripts/variable.sh: line 45: CI_BIZ_GROUP_APPID: command not found
```

### Changed

#### 升级 go-grpc-middleware 为 v2 版本

1. go-grpc-middleware v1 已被废弃；
2. 更改了服务依赖组件，以兼容 v2 版本。

#### 私有 http handler 实现函数添加错误返回

在代码文件 `handler/private.go` 中函数 `privateHTTPHandle(mux *http.ServeMux) error` 添加错误返回。

同时这里实现的 http 接口默认不支持链路可观测，需用户特殊编码后开启，见[配置使用文档](https://grpc-kit.com/docs/spec-cfg/observables/)。

## [0.3.2] - 2023-05-28

### Added

#### 更改 gitlab runner 为有向无环图 (DAG) 流水线

  1. 每个 job 均使用独立的容器来运行，避免无意义拆分多个 job；
  2. 在根目录默认生成 Dockerfile 文件；
  3. 确定默认镜像相关使用的变量名；

#### 为方便 nginx 配置路由转发，更改 swagger 使用相对地址

  1. 由原先 “/openapi-spec/microservice.swagger.json” 更改为 "./microservice.swagger.json"

     ```shell
     <body>
       <redoc spec-url='./microservice.swagger.json'></redoc>
       <script src="./redoc.standalone.js"> </script>
     </body>
     ```

  2. 在 nginx 中 location 配置

     ```shell
     location /opsaid/test1/v1/openapi-spec/ {
         proxy_pass http://opsaid-test1:10080/openapi-spec/;
         proxy_set_header  Host $http_host;
         proxy_set_header  X-Real-IP  $remote_addr;
         proxy_set_header  X-Real-Port $remote_port;
         proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
     }
     ```

   3. 实现对接口文档的转发

#### 修复 gitlab 的 check-protoc 阶段检测文件错误

  ```shell
  protoc-gen-go-grpc
  ```

#### gitlab runner 更改为 有向无环图 (DAG) 流水线

  1. 每个 job 均使用独立的 容器 来运行，适合用来运行独立的任务，加快速度，避免无意义拆分多个 job；
  2. 让 Dockerfile 默认生成；
  3. 确定默认镜像使用的变量名；

    ```shell
    CI_REGISTRY
    CI_REGISTRY_IMAGE
    CI_REGISTRY_USER
    CI_REGISTRY_PASSWORD
    ```

    进入：https://{gitlab}/-/settings/ci_cd，设置好变量。

#### 新增 jenkins 流水线模版配置

  ```shell
  .jenkins/workflows/Jenkinsfile
  ```

  依赖 k8s 环境，需提前配置好，参考 [Jenkins Pipeline](https://grpc-kit.com/docs/devops/integration/jenkins/)。

#### 统一规范 CICD 变量名

  1. 新增 `scripts/variable.sh` 用于动态变量生成；
  2. 区别 `scripts/env` 用于全局静态变量；
  3. 支持静态配置编译运行时的变量，文件路径 `scripts/env-${DEPLOY_ENV}-${BUILD_ENV}`；

#### 更改本微服务的 proto 为相对路径

  ```shell
  github.com/opsaid/test1/api/opsaid/test1/v1/microservice.proto

  更改为

  api/opsaid/test1/v1/microservice.proto
  ```

  为解决服务在容器内构建时，如果为绝对路径，则代码目录必须存放至 $GOPATH/src/$REPOSITORY 路径下，否则无法运行。

  ```shell
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
      ./api/opsaid/test5/${API_VERSION}/*.proto
  ```

  添加 paths=source_relative 这个的意思是在当前目录生成 *.pb.go 文件，而忽略 proto 文件中的 go_package 路径。

#### 统一 Makefile 与 scripts 中镜像相关的变量

- 对 `make manifests` 自动生成部署清单：

  1. 文件：Dockerfile
  2. 目录：deploy/*

- 移除 `Makefile` 中以下变量

  转移至 `scripts/env` 中做设定，因不直接在 Makefile 文件中使用，简化结构。

  ```shell
  # 构建Docker容器变量
  BUILD_GOOS      ?= $(shell ${GO} env GOOS)
  IMAGE_FROM      ?= scratch
  IMAGE_HOST      ?= hub.docker.com
  IMAGE_NAME      ?= ${IMAGE_HOST}/${NAMESPACE}/${SHORTNAME}
  IMAGE_VERSION   ?= ${RELEASE_VERSION}

  # 部署与运行相关变量
  BUILD_ENV       ?= local
  DEPLOY_ENV      ?= dev
  ```

- 更改 NAMESPACE 为部署使用的空间

  区别于 `PRODUCT_CODE` 表示产品代码或项目代码，而 `NAMESPACE` 表示租户空间，部署含义。

- 改进 `scripts/manifests.sh` 后的变量

  1. 移除 Makefile 中的 `NAMESPACE` 变量；

  ```shell
  BIZ_GROUP_APPID=hello
  DEPLOY_ENV=dev
  DEPLOY_ENV=local
  ```

  部署的环境变量，值：dev test prod stress demo staging

- 生成模版时支持自定义路径

  ```shell
  make manifests TEMPLATES=kubernetes     TEMPLATE_PATH=../gitops/deploy/kubernetes/dev/
  ```

  添加以下内容：

  1. scripts/kaniko.sh
  2. 移除 scripts/env 镜像变量
  3. 支持设置全局变量以 env-$DEPLOY_ENV-$BUILD_ENV 文件为准；

### Fixed

#### go embed 存在 .svn 异常

- 问题

  ```shell
  + make lint
  >> precheck environment
  >> generation release version
  >> generation code from proto files
  public/doc/embed.go:9:12: pattern openapi-spec/*: cannot embed directory openapi-spec/.svn: invalid name .svn
  public/doc/embed.go:9:12: pattern openapi-spec/*: cannot embed directory openapi-spec/.svn: invalid name .svn
  make: *** [Makefile:74: lint] Error 1
  ```

- 解决

  需更改为更明确的文件路径，避免使用 "*"

  ```golang
  // Code generated by "grpc-kit-cli/0.3.1-beta.1". DO NOT EDIT.
  
  package doc
  
  import (
	  "embed"
  )
  
  //go:embed openapi-spec/*.js
  //go:embed openapi-spec/*.json
  //go:embed openapi-spec/*.html
  var Assets embed.FS
  ```

## [0.3.1] - 2023-04-09

### Added

- 使用文档更新

  1. 去除 [gogo](https://github.com/gogo/protobuf) 模块文档； 
  2. 更新 grpc 地址由 https://github.com/golang/protobuf 转变为 https://github.com/protocolbuffers/protobuf-go；

- 多平台镜像构建

  1. 由于阿里云镜像中心不支持存放多架构容器更改为使用腾讯容器镜像服务；
  2. 当前多架构容器仅支持使用 "docker buildx"，暂不支持 "podman"；

- 添加依赖工具的下载

  ```shell
  make protoc
  make protoc-gen-go
  make protoc-gen-go-grpc
  make protoc-gen-grpc-gateway
  make protoc-gen-openapiv2
  ```

- 仅版本号发生变更时才执行 sed

  1. 更改了 scripts/version.sh 中的 update 方法；
  2. 仅当先前与当前版本号不一致才更改 microservice.openapiv2.yaml 文件；
  3. 更改了 /tmp/microservice.openapiv2.yaml 生成临时文件地址；

### Fixed

- 在 "oidc authenticator" 的 logger 存在空指针错误

  异常代码位置

  ```shell
  github.com/grpc-kit/pkg@v0.3.0/cfg/security.go:76
  ```

  当设置的 "oidc issuer" 可访问，但未正常返回 "/.well-known/openid-configuration" 日志输出触发了空指针。

## [0.3.0] - 2023-03-10

### Added

- 新增 "组织代码" 作为所有 proto 包名前缀

  1. 默认 "组织代码" 取值为 "grpc-kit"

  2. 根据规则自动生成内置变量：应用名称、服务包名、服务标题、服务代码

  ```shell
  APPNAME、PROTO_PACKAGE、SERVICE_TITLE、SERVICE_CODE
  ```

- 对 microservice.proto 文件中的功能注解分离并声明式

  1. 分离 "google.api.http" 功能到文件 "microservice.gateway.yaml"

     文档地址：https://github.com/googleapis/googleapis/blob/master/google/api/service.proto

  2. 分离 "grpc.gateway.protoc_gen_openapiv2.options" 功能到文件 "microservice.openapiv2.yaml"

     文档地址：https://github.com/grpc-ecosystem/grpc-gateway/internal/descriptor/openapiconfig/openapiconfig.proto

- 去掉 gogo 模块，升级 grpc-gateway v2 版本

  1. 移除了 https://github.com/gogo/protobuf 的依赖；
  2. 升级了 grpc-gateway 为 v2 版本；

- 重新规范公知类 proto 的文件存放目录

  1. 更改了 https://github.com/grpc-kit/api 原先 proto 路径规范；

  ```shell
  proto/v1/example.proto
  proto/v1/tracing.proto
  ```

  更改为以下格式：

  ```shell
  known/status/v1/response.proto
  known/example/v1/example.proto
  known/config/v1/config.proto
  ```

  2. 更改了 proto 的包名称：

  ```shell
  grpc.kit.api.proto.v1
  ```

  更改为以下前缀：

  ```shell
  grpc_kit.api.known.
  ```

- 更改库 "errors" 为 "errs" 防止对标准库重名

  1. 更改 "github.com/grpc-kit/pkg/errors" 为 "github.com/grpc-kit/pkg/errs"；
  2. 升级 proto 使用 "google.golang.org/protobuf/proto" 版本
  3. 状态使用公知版本 "grpc_kit.api.known.status.Status" 结构体

- 移除 pkg/api 中使用 gogo 类库

  1. 去除由 "protoc-gen-gogo" 生成的 "pb.go" 文件
  2. 统一使用新规范后的 "grpc-kit/api proto" 生成的 "pb.go" 文件

- 使用 gitlab-ci runner 为 shell 添加默认变量

  1. 默认模版添加以下变量；

  ```shell
  # 默认全局变量
  variables:
    CGO_ENABLED: "0"
    GIT_SSL_NO_VERIFY: "true"
    #GO111MODULE: "on"
    #GOPROXY: "https://goproxy.cn"
    #GOSUMDB: "sum.golang.google.cn"
    #GOPRIVATE: ""
    #GOPATH: "/home/gitlab-runner/go"
  ```

### Fixed

- make lint 首次无法正常运行

  1. 首次代码初始化后 "api/" 目录下不存在 "*.pb.go" 代码，导致无法引用；
  2. 通过在执行 `make lint` 之前，做 "proto" 文件的序列化；
