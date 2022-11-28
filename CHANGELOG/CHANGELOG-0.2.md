# CHANGELOG-0.2

| 名称         | 说明                     |
|------------|------------------------|
| Added      | 添加新功能                  |
| Changed    | 功能的变更                  |
| Deprecated | 未来会删除                  |
| Removed    | 之前为Deprecated状态，此版本被移除 |
| Fixed      | 功能的修复                  |
| Security   | 有关安全问题的修复              |

## [Unreleased]

## [0.2.3] - 2022-11-28

### Added

grpc-kit/pkg 模块

- 支持针对用户组进行鉴权

1. http_users 新增 "groups" 属性；
2. 如果配置了 "security.authorization.allowed_groups" 则所有需要认证鉴权的接口必须属于该组里面，否则会403；
3. 用户组区分大小写

- 添加最小化配置示例

1. 未配置的模块，将不开启该功能
2. 添加 "app-mini.yaml" 示例

- 添加健康检测服务可对外部网络

1. 添加 "HTTP GET /ping" 接口，不过 grpc 服务
2. 区别于 "HTTP GET /healthz"，该接口过 grpc 服务

grpc-kit/cli 模块

- 支持自定义应用短名称

1. 通过自定义"SHORTNAME"变量；

- 所有脚本应用shell更改

1. ubuntu默认为dash，明确使用"/bin/bash"，而非"/bin/sh"，会导致部分shell不支持

- 更改默认 http 服务端口 10080 至 8080

1. 由于chrome等浏览器默认对"10080"端口存在"ERR_UNSAFE_PORT"告警，所以更改http默认为"8080"

- cli、pkg 组件版本号统一

1. 为了解决统一编写变更记录(CHANGELOG.md)

- 支持自动生成 kubernetes 编排模版

1. 新增 "DEPLOY_ENV" 变量，表示部署环境，如：dev、test、prod
2. 新增 "BUILD_ENV" 变量，表示构建环境，一般用户自定义，默认为：local
3. 新增指令 `make manifests` 生成基于 Kubernetes 的编排清单
4. 模版路径：scripts/templates/kubernetes

## [0.2.2] - 2022-06-30

### Changed

- 版本号格式的变更

旧规则以"v"为前缀，更新后不在带"v"为前缀，如："v0.1.0-beta.3"，变更为："0.1.0-beta.3"

### Added

- 新增"VERSION"文件 

用于描述当前分支版本，同时提供给CICD使用，如果当前分支未打成"tag"，则均说明是先行版本号，同时版本去掉以"v"开头；

- 新增"Makefile"的帮助说明

```
make help
```

## [0.2.1] - 2022-06-09

### Added

- "api/doc"目录内容更改至"public/doc"
- "api/proto"目录更改为"api/{product-code}/{short-name}"
- "cli"新增"repository"参数用于说明代码仓库名
- rpc客户端、服务端实例初始化转移至"cfg"实现
- favicon.ico文件移至自定义http handler中实现
- http接口统一以"/api/"为前缀对外暴露
