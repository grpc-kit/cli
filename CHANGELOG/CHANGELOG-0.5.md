# CHANGELOG-0.5

| 名称         | 说明                     |
|------------|------------------------|
| Added      | 添加新功能                  |
| Changed    | 功能的变更                  |
| Deprecated | 未来会删除                  |
| Removed    | 之前为Deprecated状态，此版本被移除 |
| Fixed      | 功能的修复                  |
| Security   | 有关安全问题的修复              |

## [Unreleased]

### Added

- 新增 `grpc-kit-cli project migrate [path]`，默认只读预览旧项目的托管文件变更和人工迁移清单；`--apply` 仅在正式稳定版 CLI、clean Git 工作树和输入摘要复核通过时更新已有托管文件。
- 首期迁移目标固定为正式 `github.com/grpc-kit/pkg v0.5.0`。命令严格区分 CLI marker 版本和 pkg 版本，不修改 `go.mod`、业务源码及其他用户管理文件，也不提供可能被误解为 CLI 自更新的根级 `upgrade`/`update` alias。

### Changed

#### grpc-kit/cli 模块

- **Breaking**：服务模板日志从 logrus 切换为 Go 标准库 `log/slog`。

  1. `Microservice`、`IndependentCfg`、flow client 与 MCP registrar 的 logger 类型改为 `*slog.Logger`；`GetLogger`、`WithLogger`、`WithWorkflow`、`flow.NewClient` 等 API 名称保持不变
  2. RPC、MCP 和 shutdown 日志改用 slog Context API；新模板不再直接依赖 `github.com/sirupsen/logrus`
  3. 新模板固定使用 `github.com/grpc-kit/pkg v0.5.0`

- **Breaking**：初始化、注册、注销和服务启动链路统一直接传递 `context.Context`。

  1. `NewMicroservice(ctx, ...)`、`IndependentCfg.Init(ctx, ...)` 和 `LocalConfig.Init(ctx)` 直接接收启动上下文
  2. `LocalConfig.HTTPHandlerFrontend(ctx, ...)`、`sd.Register(ctx, ...)`、`LocalConfig.Deregister(ctx)` 与 `sd.Registry.Deregister(ctx)` 直接接收调用方上下文
  3. `rpc.Server.StartBackground(ctx)` 直接接收启动上下文，使启动等待可以响应取消；旧生成项目中的 `StartBackground()` 调用需补充 ctx
  4. 移除临时的 `InitContext`、`DeregisterContext`、`HTTPHandlerFrontendContext` 和 `sd.RegisterContext` 双入口，不保留无 ctx 兼容包装

- **Breaking**：`errs.Status.WithLogger` 增加 ctx 首参并删除 `WithLoggerContext`；自定义业务调用需迁移为 `WithLogger(ctx, logger, format, err)`。

- **Changed**：服务模板与 `pkg v0.5.0` 的 `go` 指令提升至 `1.25.13`（随安全补丁升级 grpc v1.83.1、x/net v0.56.0、x/text v0.39.0 并修复 Go 标准库已知漏洞）；生成与升级项目需使用 Go 1.25.13 及以上版本构建。

#### 旧生成项目迁移

升级到 `github.com/grpc-kit/pkg v0.5.0` 时，旧项目需要同步完成以下机械迁移：

1. 将 logger 类型和构造逻辑从 logrus 改为 `*slog.Logger`
2. 为 `NewMicroservice`、`Init`、`HTTPHandlerFrontend`、`sd.Register`、`Deregister`、`StartBackground` 和 `WithLogger` 调用补充已有 ctx
3. 删除对临时 `*Context` 方法及 logrus 兼容入口的调用
4. 运行代码生成、单元测试与构建，确认自定义 handler、注册流程和关闭流程均已适配新签名
