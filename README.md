# gRPC Kit

主要基于以下几个核心类库实现：

- [grpc](https://github.com/protocolbuffers/protobuf-go)
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

## 简要概述

快速生成微服务模版，为同一产品提供统一的治理方式，提高多人协作效率

## 生成代码并运行（本机方式）

### 创建应用模版

```shell
mkdir -p $GOPATH/src/github.com/opsaid
cd $GOPATH/src/github.com/opsaid

grpc-kit-cli new -t service -p opsaid -s test1
```

### 下载依赖的环境

```shell
make protoc
make protoc-gen-go
make protoc-gen-go-grpc
make protoc-gen-grpc-gateway
make protoc-gen-openapiv2
```

### 运行应用代码

```shell
make run
```

## 生成代码并运行（容器方式）

### 创建应用模版

```shell
docker run \
    --rm \
    -v $(pwd):/usr/local/src \
    -w /usr/local/src \
    ccr.ccs.tencentyun.com/grpc-kit/cli:0.3.0 \
    grpc-kit-cli new -t service -p opsaid -s test1
```

### 运行应用代码

```shell
docker run -i -t --rm \
    -v $GOPATH/pkg:/go/pkg \
    -v $(pwd):/usr/local/src \
    -w /usr/local/src \
    --network host \
    ccr.ccs.tencentyun.com/grpc-kit/cli:0.3.0 \
    make run
```

## 服务访问测试

- 微服务接口文档

```shell
http://127.0.0.1:8080/openapi-spec/
```

- 微服务编译版本

```shell
# curl http://127.0.0.1:8080/version | python -m json.tool

{
    "appname": "test1.v1.opsaid",
    "build_date": "2023-01-13T09:10:45Z",
    "git_commit": "1234567890123456789012345678901234567890",
    "git_branch": "",
    "go_version": "go1.18.5",
    "compiler": "gc",
    "platform": "darwin/amd64",
    "cli_version": "0.2.3",
    "commit_unix_time": 0,
    "release_version": "0.1.0"
}
```

- 微服务性能数据

```shell
# curl http://127.0.0.1:8080/metrics
```

```shell
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000114581
go_gc_duration_seconds{quantile="0.25"} 0.000873528
go_gc_duration_seconds{quantile="0.5"} 0.002296699
go_gc_duration_seconds{quantile="0.75"} 0.003722618
go_gc_duration_seconds{quantile="1"} 0.010592338
go_gc_duration_seconds_sum 0.033207328
go_gc_duration_seconds_count 12
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 19

...
```

- 微服务健康探测

探测流量仅到 gateway 不会调度到 grpc 服务。

```shell
# curl http://127.0.0.1:8080/ping
OK
```

探测流量同时到 gateway 与 grpc 服务。

```shell
# curl 'http://127.0.0.1:8080/healthz?service=test1.v1.opsaid'
{"status":"SERVING"}
```

- 示例 demo 接口

```shell
# curl -u user1:grpc-kit-cli http://127.0.0.1:8080/api/demo
```
