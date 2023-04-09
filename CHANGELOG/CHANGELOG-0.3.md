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