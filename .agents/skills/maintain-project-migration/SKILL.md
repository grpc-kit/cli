---
name: maintain-project-migration
description: 维护 grpc-kit CLI 服务模板、project migrate 兼容资产和正式 pkg 版本矩阵。适用于模板变化、pkg 版本升级、增加来源 fixture、调整 ManualAction 或发布迁移能力；不用于代替业务项目完成人工源码迁移。
---

# 维护项目迁移能力

在 `grpc-kit/cli` 仓库根目录工作。开始修改前读取：

- `../adm/docs/roadmap/cli/01-cli-project-upgrade-feasibility.md`
- `CHANGELOG/CHANGELOG-0.5.md` 及本次目标版本对应的 CHANGELOG
- `internal/projectmigrate/target.go`，以及同目录的资产、来源 fixture、Plan、诊断和兼容测试
- [维护 SOP](references/sop.md)

先分类变更，再决定是否触碰旧项目迁移资产。最新 `template/service`、某个历史来源模板和某个 pkg 迁移目标是三个不同对象；不能因为新建项目模板变化就把新功能注入旧项目。

始终保持以下约束：

- `TargetCLIVersion` 只表示执行迁移的正式 CLI 版本；pkg 目标版本独立记录。
- 只有首行完整 grpc-kit-cli `DO NOT EDIT` marker 的已有普通文件可进入 Change Plan；`scripts/generate.sh` 的 marker 位于第二行（首行是 shebang），0.5.1 之前无 marker 的存量世代凭"路径精确 + 全文 SHA-256 匹配冻结摘要"的内容证据窄例外整文件托管，维护规则见 SOP §3。
- 不自动修改用户代码、`go.mod`、`go.sum`，不创建、删除、重命名或 chmod 项目文件。
- 用户可以另行使用共享技能 `migrate-logrus-to-slog` 完成业务源码迁移；该授权属于 AI 编辑工作流，不得反向扩大 `project migrate` 的 Change Plan。
- 兼容资产只引用来源 fixture 已有的项目内符号；新功能走新建模板或独立 feature 流程。
- 正式 pkg 测试不使用本地 `replace`，也不使用不可复现的 `@latest`。
- 不擅自修改 `VERSION`、创建 tag、发布产物或对脏的真实服务执行 apply。

完成维护后报告：变更分类、支持的来源与目标版本、资产写集合、最低/候选 pkg 验证结果、未执行的完整项目检查及其原因。
