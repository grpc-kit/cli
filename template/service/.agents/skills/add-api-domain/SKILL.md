---
name: add-api-domain
description: 为 grpc-kit 服务新增 API 领域，包括 proto 契约、HTTP 映射、处理器和 API 测试。适用于新增领域或相关 RPC。
---

# 新增 API 领域

在生成服务的根目录执行任务。修改文件前，先阅读 `AGENTS.md`、`scripts/env`、`Makefile` 和 `test/e2e/README.md`。检查一个已有 API 领域，确认当前 package 和文件约定，不要假设模板占位路径仍然适用。

1. 在现有 `api/` package 下定义带版本的消息和服务方法，并遵循仓库现有的 AIP 约定。在 `microservice.proto` 中增加 import 和 RPC，然后更新对应的 gateway 与 OpenAPI YAML。只有需求确实涉及持久化时才修改 Ent Schema。
2. 需要生成绑定代码来实现 handler 时运行 `make generate`，并检查生成差异。不要修改带有 `DO NOT EDIT` 标记的文件，也不要把代码生成描述成能够完成业务 handler 实现。
3. 在允许业务编辑的 handler 文件中实现 RPC，并使用已有的注册扩展点。传输层映射和领域行为应与相邻 handler 保持一致。
4. 在实现代码旁增加有实际断言的单元测试。当 API 契约可从外部观察时，在 `test/e2e/<domain>/` 下增加黑盒测试。复用 `test/e2e/client` 和 `test/e2e/fixture`，添加 `e2e` build tag 和每个 package 的 `TestMain`，只断言外部 HTTP 契约。不能用空测试、无条件跳过或没有断言的 stub 充当已完成测试。
5. 运行 `make test`，完成代码生成和单元测试；该命令不会运行 E2E 测试。如果本地已有配置正确的服务，再运行 `make test-e2e`。否则可在条件允许时使用 `go vet -tags=e2e ./test/e2e/...` 做编译和静态检查，并明确说明没有执行运行时 E2E。

完成后列出修改过的契约文件、生成产物、handler 和测试，并报告准确的验证结果以及因环境条件未执行的检查。
