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

## [0.2.2] - 2022-06-09

### Added

- 新增"VERSION"文件 

用于描述当前分支版本，同时提供给CICD使用，如果当前分支未打成"tag"，则均说明是先行版本号；

## [0.2.1] - 2022-06-09

### Added

- "api/doc"目录内容更改至"public/doc"
- "api/proto"目录更改为"api/{product-code}/{short-name}"
- "cli"新增"repository"参数用于说明代码仓库名
- rpc客户端、服务端实例初始化转移至"cfg"实现
- favicon.ico文件移至自定义http handler中实现
- http接口统一以"/api/"为前缀对外暴露
