# CHANGELOG-0.4

| 名称         | 说明                     |
|------------|------------------------|
| Added      | 添加新功能                  |
| Changed    | 功能的变更                  |
| Deprecated | 未来会删除                  |
| Removed    | 之前为Deprecated状态，此版本被移除 |
| Fixed      | 功能的修复                  |
| Security   | 有关安全问题的修复              |

## [Unreleased]

## [0.4.0] - 2026-07-08

### Added

#### grpc-kit/api 模块

- 新增 `KnownAdmin` 内置管理服务完整 Proto 定义

  在 `known/admin/v1/` 目录下新增以下 Proto 文件，定义了完整的后台管理 API 接口：

  1. `admin.proto` — `KnownAdmin` 服务入口，聚合所有管理 RPC 方法
  2. `admin.common.proto` — 公共消息定义（AuthProvider、User、Role、Department、Group、Menu、Policy、OAuth2Client 等）
  3. `admin.gateway.yaml` — grpc-gateway 路由配置
  4. `admin.openapiv2.yaml` / `admin.swagger.json` — OpenAPI / Swagger 文档
  5. `auth.proto` — 登录认证、MFA 多因素认证、认证提供方管理、登录选项列表
  6. `code.proto` — Code 校验规则定义
  7. `config.proto` — 本地配置管理、全局设置管理
  8. `database.proto` — 数据库初始化接口
  9. `department.proto` — 部门管理及成员关系
  10. `group.proto` — 群组管理及成员关系
  11. `menu.proto` — 菜单管理
  12. `oauth2.proto` — OAuth2 客户端管理
  13. `policy.proto` — 策略管理
  14. `roles.proto` — 角色管理（含成员、菜单、策略关联）
  15. `security.proto` — 安全密钥与 OAuth2 Discovery/JWKS/Userinfo
  16. `service.proto` — 服务字典与 Action 管理
  17. `setting.proto` — 全局设置管理
  18. `user.proto` — 用户管理

- 新增多因素认证（MFA）接口

  支持基于 TOTP 的多因素认证，包含：
  1. 登录时 MFA 二步验证（`VerifyAuthMFA`）
  2. 登录态首次配置 MFA（`StartAuthMFASetup` / `ConfirmAuthMFASetup`）
  3. 用户级 MFA 管理（`SetupUserMFA` / `ConfirmUserMFA` / `DisableUserMFA`）

- 新增 OAuth2 客户端管理接口

  支持 OAuth2 客户端的完整 CRUD 操作，包含客户端密钥的自动生成与 bcrypt 哈希存储。

- 新增登录页认证提供方列表接口

  `ListLoginOptions` 接口无需认证，返回最小字段集供登录页展示可用认证方式。

- 新增全局设置管理接口

  支持按分类获取和更新全局配置（`GetGlobalSettings` / `ListGlobalSettings` / `UpdateGlobalSettings`）。

- 新增本地配置查询接口

  支持获取本地配置列表及按模块获取配置（`ListLocalConfigs` / `GetLocalConfigs`）。

- 新增数据库初始化接口

  `CreateDatabaseInitialize` 支持通过密码哈希初始化管理员账户。

- 新增服务字典与 Action 管理接口

  `ListServices` / `ListServiceActions` 用于查询已注册的微服务及其操作定义。

- 新增角色策略关联接口

  支持角色与策略的绑定管理（`ListRolePolicies` / `CreateRolePolicies` / `UpdateRolePolicies` / `DeleteRolePolicy`）。

- 新增角色菜单关联接口

  支持角色与菜单的绑定管理（`ListRoleMenus` / `CreateRoleMenus` / `UpdateRoleMenus` / `DeleteRoleMenu`）。

- 引入通用可见性枚举与 `protected` 字段

  为部门、角色、认证提供方等实体添加 `protected` 字段，防止误删内置数据。

#### grpc-kit/pkg 模块

- 新增 Admin 管理服务完整实现

  在 `admin/` 目录下实现了 `KnownAdmin` 服务的所有 RPC 方法，涵盖：
  1. 认证登录（`rpc_auth_login.go`）— 支持本地静态用户、LDAP、OIDC 回调登录
  2. MFA 多因素认证（`rpc_auth_mfa.go`、`mfa_challenge.go`、`mfa_gate.go`）— TOTP 验证、恢复码机制
  3. 认证回调（`rpc_auth_callback.go`）— OIDC 回调用户信息入库
  4. 用户管理（`rpc_users.go`、`social_users.go`、`static_users.go`）
  5. 角色管理（`rpc_roles.go`、`principal_role_helpers.go`）
  6. 部门管理（`rpc_departments.go`）
  7. 群组管理（`rpc_groups.go`）
  8. 菜单管理（`rpc_menus.go`）
  9. 策略管理（`rpc_policies.go`、`role_policy_helpers.go`）
  10. 角色策略关联（`rpc_role_policies.go`）
  11. OAuth2 客户端管理（`rpc_oauth2_clients.go`）
  12. 安全密钥（`rpc_security_keys.go`）— JWT 签名证书与 JWKS 管理
  13. 全局设置（`rpc_global_settings.go`）
  14. 本地配置（`rpc_local_config.go`）
  15. 服务字典（`rpc_services.go`）— 解析 `admin.gateway.yaml` 并支持 Swagger 集成
  16. 数据库初始化（`rpc_database.go`）— 内置菜单、部门、管理员账户种子数据
  17. 微信开放平台登录（`wechat_open.go`）
  18. 内置种子代码管理（`well_known_seed_code.go`）

- 新增 Lion ORM 数据模型层

  基于 [ent](https://entgo.io/) 框架实现完整的数据持久化层，包含以下实体模型：
  1. `users` — 用户
  2. `useridentities` — 用户身份（本地密码、社交登录）
  3. `usermemberships` — 用户成员关系
  4. `userprofiles` — 用户属性扩展
  5. `authproviders` — 认证提供方
  6. `credentials` — 凭证（安全密钥）
  7. `departments` — 部门
  8. `groups` — 群组
  9. `menus` — 菜单
  10. `roles` — 角色
  11. `rolemenus` — 角色-菜单关联
  12. `rolepolicies` — 角色-策略关联
  13. `principalroles` — 主体-角色关联
  14. `policies` — 策略
  15. `globalsettings` — 全局设置
  16. `oauth2clients` — OAuth2 客户端
  17. `oauth2codes` — OAuth2 授权码

- 新增 OPA 权限策略引擎集成

  集成 [Open Policy Agent (OPA)](https://www.openpolicyagent.org/) 作为权限策略验证引擎：
  1. 支持内置 Rego 规则与 RBAC 数据
  2. 支持 `opa_envoy_plugin` 外部鉴权服务
  3. 支持动态 data provider 配置
  4. 支持获取本地 RBAC 数据转化为 Envoy RBAC 结构
  5. 为自定义 HTTP handler 复用鉴权配置

- 新增基于 CloudEvents 的审计日志模块

  实现完整的 gRPC 审计拦截器，支持：
  1. Unary 与 Stream 方法的审计事件记录
  2. 基于 [CloudEvents](https://cloudevents.io/) 规范的事件封装
  3. 通过 Kafka（IBM/sarama）传输审计事件
  4. 审计事件推送失败指标（`grpc_kit.audit_event.send_errors`）
  5. 支持 JSON 序列化配置
  6. 忽略 RPC 健康检查的审计事件

- 新增缓存模块

  1. 支持 Redis 缓存（含 Cluster / Sentinel 模式）
  2. 实现内存 LRU 缓存
  3. 支持 TTL 过期策略
  4. Redis 操作集成 OpenTelemetry 链路与指标
  5. 使用 gob 序列化缓存至 Redis
  6. 自动添加内部缓存键前缀
  7. 统一缓存操作返回 bool 类型

- 新增 GRN（通用资源名称）解析与匹配

  借鉴 AWS ARN 设计 6 段固定结构的资源标识符，支持通配符匹配，与 Rego 端 `grn_match` 函数行为一致。

- 新增 Automations 流程编排配置

  支持工作流配置参数获取与脚本内容获取。

- 新增配置快照（Config Snapshot）功能

- 新增前端静态文件托管增强

  1. 支持后台前端 history fallback 效果
  2. 内嵌 Swagger UI 提供 OpenAPI 文档展示
  3. 解析 `admin.gateway.yaml` 微服务网关配置

- 新增加密工具

  1. `crypto/bcrypt.go` — bcrypt 哈希工具
  2. `crypto/aes.go` — ent 字段 AES 加解密
  3. `crypto/sha.go` — 公共 SHA256 计算函数

- 新增 LDAP 登录支持

  支持 LDAP 用户登录流程，包含属性同步功能。

- 新增静态用户登录认证

  支持基于本地配置文件的静态用户登录，静态用户加入 `user_id` 属性。

- 新增 `/debug/` 端点内部网络访问限制

- 新增返回 `x-` 开头 HTTP 请求头的支持

- 新增更丰富的错误信息自定义能力

### Changed

#### grpc-kit/api 模块

- 统一 ID 字段类型为 `int64`

  原先部分实体 ID 使用 `int32`，统一更改为 `int64` 以支持更大范围的标识符。

- 统一名称字段为系统标识符（`code`）与展示名称（`display_name`）

  将 `name` 字段更改为 `code`，仅支持字母编码，用于系统标识。

- 重命名 `expired_at` 为 `expires_at`

  统一过期时间字段命名。

- 重命名 `order_weight` 为 `sort_order`

  统一排序字段命名。

- 重命名用户关系模型

  1. `UserGroup` → `GroupMember`
  2. `UserDepartment` → `DepartmentMember`

- 统一列表 RPC 方法使用 `List` 前缀

- 重构成员关系模型以统一管理群组和部门成员

- 重构国际化名称结构，移除 `I18NName` 中的 `default_language` 字段

- 更新分页请求消息描述以增强清晰度和一致性

- 更改手机号格式以兼容国际化

- 更改 `idcard` 为 `national_id`

- 统一接口前缀

#### grpc-kit/pkg 模块

- 升级 OpenTelemetry 至 v1.43.0

  1. 移除 grpc stream otelgrpc，改用 grpc stats 实现
  2. 降低指标 bucket 边界数
  3. 改进 otel grpc 使用方式
  4. 未获取到 trace id 时使用 UUID 填充审计 ID

- 升级 Kafka SDK：Shopify/sarama → IBM/sarama

- 升级 OPA 至 v1.x 版本

- 升级 grpc-gateway 至 v2.28.0

- 统一 ID 字段类型从 `int32` 更改为 `int64`

- 将 `name` 更改为 `code`，仅支持字母编码

- 重命名 `expired_at` 为 `expires_at`

- 注册至配置中心使用 JSON 代替 YAML 格式

- 为 admin 等自定义 HTTP handler 添加鉴权中间件

- 优化 Redis 与 Memory 缓存实现

- 统一参数设置使用 `WithXXX` 风格

- 更改 Envoy 请求 URL 函数

- 移除 URL-scheme 定义，添加 scheme metadata 传递

- 对初始化函数 `InitX` 更改以避免不必要暴露

- 独立 gRPC 审计拦截器至单独模块

- 对 ctx 中获取的数据独立至 `rpc` 模块

- 减少 user 认证相关对外可视方法

- 统一缓存操作返回 bool 类型

- 对 JSON marshal 提取为公共配置

- 标准化本地配置快照中的服务名称

- 统一成员处理逻辑（跨部门和群组）

### Removed

#### grpc-kit/api 模块

- 移除已废弃的 Proto 文件

  移除 `policy_attachments`、`policy_statements`、`resources`、`resource_types` 等 Proto 定义，精简 API 结构。

- 移除 `Scope` 消息，更新 `Action` 消息结构

- 移除 `surface_mask` 字段

- 移除与资源和角色权限相关的未使用消息

- 移除 ent demo schema 及代码生成文件（`ent/`、`ent/schema/demo.go`、`ent/entc.go`、`ent/generate.go`）

#### grpc-kit/pkg 模块

- 移除已废弃的 `policy_statement.proto` 文件及其生成代码

- 移除 `surface_mask` 字段及菜单相关功能

- 移除 Services 实体管理代码

- 移除未使用的 `Scopes` 和 `PermissionBindings` 事务客户端

- 移除 `lion_user_roles` 边及关联功能

- 移除不再使用的 SQLite3 依赖

- 移除未使用的 schema 文件及关联事务客户端

### Fixed

#### grpc-kit/pkg 模块

- 修复审计拦截器日志记录方法（使用 `Warn` 替代 `Warnf`）
- 修复推送审计日志 response 时的空指针错误
- 修复认证关闭下 Security Authorization 的空指针问题
- 修复 HTTP 与 gRPC 之间链路请求 ID 传递丢失
- 修复 JWT 为 HS256 签名时缺失 ctx 设置 groups
- 修复对象存储默认上传分块大小为 64M
- 修复配置项 `independent` 中 `time.Duration` 类型无法解析
- 修复用户 ID 类型转换问题（int → int64）
- 修复 JWT 令牌验证逻辑
- 修复部门成员类型默认值处理
- 修复认证提供商类型设置问题
- 修复 HTTP basic auth 认证无法识别 `password_hash`
- 修复用户名为空时使用 subject 作为用户名
- 修复 URL 由 `certs` 至 `jwks`
- 修复共享实例不能存储 RPC method 变化的内容
- 修复遗漏 `hide_children_in_menu` 后缀
- 修复默认 RBAC 示例属性 principals 编写错误
- 修复自定义 HTTP handler 仅在开启鉴权时输出 input 日志
- 修复对 RBAC 数据添加包头
- 修复 Open Policy Agent 相关包导入路径至 v1 版本
