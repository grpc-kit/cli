# project migrate 维护 SOP

## 1. 先判断变化属于哪一类

| 变化 | 默认处理 |
|---|---|
| 仅影响以后 `new` 创建的项目 | 更新 `template/service` 与模板测试；不要修改历史兼容资产 |
| pkg 同一目标族的补丁版本，且对既有托管代码向后兼容 | 保留最低版本基线，增加候选版本兼容测试；不因“版本更新”重渲染资产 |
| pkg 新签名要求托管代码变化，或提高最低 Go 版本 | 建立新的明确迁移目标和资产根；同步 Plan、诊断、测试和文档 |
| 增加可迁移的历史 CLI 来源 | 从 tag 或确认 commit 冻结来源摘要，证明其可复用某一资产族后再登记 |
| 修复生成器缺陷 | 修复最新模板；历史文件只有在路径和全文证据充分时增加窄例外，不放宽 marker |
| 业务文件或扩展点变化 | 产生 ManualAction；除非文件仍满足托管所有权规则，否则不得自动写入 |

`v0.x` 版本不能只凭版本号假设兼容。先检查正式 tag、目标模块 `go.mod`、CHANGELOG 和实际 API 签名，再确定属于“扩展测试矩阵”还是“新迁移目标”。

## 2. 建立可审计的版本事实

1. 确认 pkg 目标是已发布 tag，而不是本地分支、伪版本或 CHANGELOG 标题。
2. 记录 pkg 最低 Go 版本及受影响 API；核对 context、logger、注册、启动和关闭签名。
3. 区分三个版本：来源 CLI、执行迁移的目标 CLI、目标 pkg。禁止用其中一个推导另外两个。
4. `internal/projectmigrate/target.go` 是当前迁移目标 pkg/Go 版本的代码内来源；修改目标时仍需用 `rg` 找出模板、资产目录、Makefile 默认值、README 和 CHANGELOG 中的所有引用，避免只改常量或测试字符串。
5. 如果正式 CLI 版本尚未决定，可以完成 preview 和 fixture 验收，但必须继续拒绝 apply；不得先写一个猜测版本。

## 3. 维护模板与兼容资产

1. 更新 `template/service` 时，先判断变化是否只服务新项目。MCP、Flow、E2E 等新增能力不得自动进入旧来源资产。
2. 需要更新旧项目托管代码时，在 `internal/projectmigrate/assets/<pkg-target>/<source-family>/` 维护独立资产。资产必须：
   - 保留完整首行 marker；
   - 只使用 `scripts/env` 和 `go.mod` 可恢复的参数；
   - 不引用来源 fixture 不存在的本地符号；
   - 渲染后无模板表达式残留，并可通过 Go 语法检查。
3. 新增来源版本时，从发布 tag 或确认 commit 冻结路径和 SHA-256。真实脏工作区只能用于只读调查，不能作为 fixture 来源。
4. 新增 pkg 迁移目标时保留旧目标资产和回归测试。不要原地把 `v0.5.0` 资产目录解释成 `v0.6.0`。
5. 用户管理文件中的破坏点进入 `ManualAction`。诊断应提供稳定 code、文件、行号和建议，并排除本次 Change Plan 已整体替换的路径。
6. logrus/slog 用户源码可以推荐使用 grpc-kit 共享技能 `migrate-logrus-to-slog`。技能由用户显式授权 AI 编辑并负责语义迁移、依赖清理和完整项目验证；它不是 CLI 自动写入能力，也不能让未托管文件进入 Change Plan。

## 4. 正式 pkg 兼容矩阵

当前最低基线使用：

```shell
make test-compatibility PKG_COMPAT_VERSION=v0.5.0
```

当候选补丁版本发布后，同时执行最低版本和候选版本，例如：

```shell
make test-compatibility PKG_COMPAT_VERSION=v0.5.0
make test-compatibility PKG_COMPAT_VERSION=v0.5.8
```

最低版本证明迁移资产没有误用较新 API；候选版本证明向前兼容。命令必须使用正式模块且无本地 `replace`。不要使用 `@latest`，也不要删除最低版本门禁来换取候选版本通过。

如果候选版本失败：

- 属于 pkg 非预期回归：先报告并修复/发布 pkg，不改 CLI 资产掩盖问题；
- 属于明确的新 pkg 契约：建立新迁移目标，更新资产和 ManualAction；
- 只影响新建模板：更新 `template/service`，保持历史迁移资产不变。

## 5. 验收顺序

1. 运行 `gofmt` 和 `make test`。
2. 对最低 pkg 和本次候选 pkg 分别运行 `make test-compatibility PKG_COMPAT_VERSION=<version>`。
3. 对每个支持来源验证 preview：状态、warning、ManualAction 和 Changes 路径必须稳定，preview 前后项目逐字节不变。
4. 在从真实服务 clean commit 导出的临时 Git fixture 上执行 apply；断言 Git diff 路径严格等于 Plan，且没有 create/delete/mode change。
5. 再次 preview 必须为 `managed_up_to_date`。
6. 如果本次声明“完整项目兼容”，还必须在临时 fixture 中完成人工待办并执行该项目的生成、测试和构建。托管资产编译通过不能替代这一步。
7. 运行 CLI 与文档仓库的 `git diff --check`，更新 README、目标 CHANGELOG 和路线文档中的真实验收结果。

## 6. 停止条件

出现以下任一情况时，不得宣称可以发布或执行真实项目 apply：

- pkg 目标没有正式 tag，或正式模块无法获取；
- 来源版本没有发布产物/确认 commit 和冻结摘要；
- 兼容资产需要来源项目不存在的符号或必须创建文件；
- 正式 CLI 版本未决定、为空、`v0.0.0` 或 prerelease；
- 真实项目 Git 不 clean；
- Plan 有 unsupported/conflict，或 compatibility test 失败；
- 只完成托管文件验证，却声称业务项目已经完成编译兼容迁移。
