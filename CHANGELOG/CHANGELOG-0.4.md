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

## [0.4.2] - 2026-08-05

### Added

#### grpc-kit/pkg 模块

- 新增标准化 Access Token Claims 与统一签发器

  1. `auth/token.go` 新增 `CommonClaims` 与 `AccessTokenClaims`，分离 OIDC ID Token 和 OAuth 2.0 Access Token 的语义
  2. Access Token 新增 `client_id` / `scope` / `roles` 声明，并统一携带 `sub` / `iat` / `exp` / `jti` / `tenant` / `groups`
  3. 新增 `BuildAccessTokenClaims` 以及 HS256 / RS256 签名辅助函数；Access Token 的 JWT `typ` 为 `at+jwt`，ID Token 为 `JWT`
  4. 统一本地静态用户、LDAP、OIDC、OAuth2 与 MFA 验证后的令牌签发流程

- 新增基于 `roles` 的授权配置与 Context API

  1. `security.authorization` 新增 `allowed_roles`，`http_users` 新增 `roles`
  2. `cfg.LocalConfig` 新增 `AccessTokenFrom()` / `RolesFrom()`
  3. `rpc` 新增 `ContextWithTokenClaims()` / `GetTokenClaimsFromContext()` 与 roles Context 存取方法
  4. 旧的 `allowed_groups` 继续兼容；`IDTokenFrom()` 和 ID Token Context 方法保留一个兼容周期

- 新增 JWKS 凭证生命周期保护

  1. 创建 KEY_PAIR / X.509 凭证时，支持将 PEM 或 DER 公钥、私钥和证书统一归一化为 DER 存储，PKCS#8 RSA 私钥转换为 PKCS#1 DER
  2. 禁止禁用或删除最后一个活跃 JWKS 凭证，确保系统始终保留可用签名密钥
  3. JWKS 端点同时返回 `ACTIVE` 和 `EXPIRED` 公钥，便于密钥轮换后继续验证尚未过期的历史令牌

- OAuth2 Userinfo 新增 HTTP Basic Auth 支持

  `/builtin/admin/api/v1/oauth2/userinfo` 现可使用 Bearer Access Token 或 `app.yaml` 中配置的 HTTP Basic 用户访问；在未配置数据库时也可返回最小身份信息。

### Changed

#### grpc-kit/pkg 模块

- 分离角色与群组语义

  `roles` 专用于授权，`groups` 仅表示用户的群组成员关系；权限前置门和 `allowed_roles` 均仅检查角色。`allowed_groups` 仅作为旧配置名兼容，其值仍按 role code 解释；同时配置两者时必须保持集合一致，否则拒绝授权。

- 调整 JWT / OIDC 声明模型

  1. 新签发令牌使用 OIDC 标准 `preferred_username`，仍兼容读取旧 `username` 声明
  2. 新 Access Token 不再签发历史 `appid` / `username` 声明，也不再生成 `<username>@localhost` 占位邮箱
  3. `IDTokenClaims` 补充 `nonce` / `azp` / `auth_time` / `acr` / `amr` / `at_hash` / `c_hash` 等 OIDC 标准声明

- 改进角色与群组的有效性计算

  签发令牌时分别计算 role code 与 group code，并过滤未激活、已过期或已删除的成员关系、群组、部门和角色。

- 调整 MFA Challenge 令牌签发上下文

  MFA 挑战保留 `client_id` / `scope` / `tenant` / TTL，确保完成 MFA 后签发的 Access Token 与原始登录请求一致；Challenge 读写改为深拷贝和显式更新暂存密钥。

- 升级依赖

  1. `github.com/golang-jwt/jwt` 从 v4.5.2 升级至 v5.3.1
  2. `github.com/modelcontextprotocol/go-sdk` 从 v1.6.1 升级至 v1.7.0

- 升级注意

  1. 仅包含 `groups` 而不包含 `roles` 的旧 Access Token 不再获得角色权限，升级后应让用户重新登录以换取新令牌
  2. 配置项建议从 `allowed_groups` 迁移至 `allowed_roles`；迁移期间同时配置时，两者的 role code 必须一致
  3. 直接构造 `auth.IDTokenClaims` 的调用方需适配内嵌的 `CommonClaims`；业务 Access Token 请改用 `AccessTokenClaims`、`AccessTokenFrom()` 与新 Token Claims Context API

### Fixed

#### grpc-kit/pkg 模块

- 修复 `SetEmail()` 条件反转导致非空邮箱未被写入的问题，并按上游 `email_verified` 声明设置 GitHub / OAuth2 用户的邮箱验证状态。

- 修复角色创建和更新时遗漏 `status` 字段的问题；创建时未指定状态默认为 `ACTIVE`。

- 修复受保护凭证无法更新 `status` 的问题，保留 `display_name` / `description` / `status` 三个可更新字段。

- 修复 HS256 验证在启用 `skip_expiry_check` 时的过期校验处理，改用 jwt/v5 `WithoutClaimsValidation` 跳过声明校验，但仍保留签名验证。

### Security

#### grpc-kit/pkg 模块

- OIDC 登录改为通过 Provider Discovery 获取 Verifier，并校验上游 `id_token` 的签名、Issuer、Audience 与有效期后再读取 Claims。

- 授权不再将 `groups` 声明回退解释为角色，避免群组成员关系与权限角色混用造成误授权。

## [0.4.1] - 2026-07-25

### Added

#### grpc-kit/api 模块

- 新增企业级凭证（Credential）管理完整 API 定义

  在 `known/admin/v1/` 下新增凭证管理的完整接口，支持 API Key、对称密钥、非对称密钥对、X.509 证书、软件许可证、通用密文六类凭证：

  1. 新增 `Credential` 核心消息（`admin.common.proto`），采用分区设计（标识区 -> 展示区 -> 分类区 -> 管理区 -> 密钥材料区 -> 生命周期区 -> 扩展区 -> 审计区），含 `Type`/`Algorithm`/`Usage`/`Status`/`Source` 嵌套枚举，密钥材料以 `oneof key_material` 表达
  2. 新增 6 个 RPC 方法（`admin.proto`）：`CreateCredential` / `ListCredentials` / `GetCredential` / `UpdateCredential` / `DeleteCredential` / `RevealCredentialSecret`
  3. 新增 `security.proto` 请求/响应消息，支持 cursor/offset 双分页、AIP-160 `filter`、`order_by`、`FieldMask` 局部更新
  4. `RevealCredentialSecret` 支持用户密码 SHA256 二次验证后按类型揭示敏感字段
  5. 新增 `CredentialCode` 种子枚举（`code.proto`），含 `CREDENTIAL_CODE_JWT_SIGNING_V1`（DB code: `jwt-signing-v1`）
  6. 新增 gateway 路由（`admin.gateway.yaml`，前缀 `/builtin/admin/api/v1/credentials`）及 OpenAPI/Swagger 文档

#### grpc-kit/pkg 模块

- 新增凭证（Credential）管理服务实现

  在 `admin/rpc_security_keys.go` 实现 `KnownAdmin` 服务的全部凭证 RPC 方法：

  1. 敏感字段（`api_secret`/`private_key`/`passphrase`/`license_key`/`symmetric_key`）使用 AES-GCM 加密存储，`credentialToProto` 永不映射 `*_encrypted` 字段
  2. `computeFingerprint` 按凭证类型计算 SHA-256 摘要（64 字符 hex），用于幂等去重；`SECRET` 类型 Token 持久化（`type=SECRET + usage=AUTH + source=USER`）跳过指纹检查
  3. `RevealCredentialSecret` 校验当前用户密码 SHA256 后解密返回敏感字段
  4. `CreateCredential` 空 `code` 时自动生成 12 位随机码（复用 `schema.EnsureCode`）
  5. `UpdateCredential` 通过 `FieldMask` 支持局部更新，`code`/`type`/`algorithm`/`key_material`/`fingerprint` 不可修改
  6. `DeleteCredential` 软删除，受保护凭证（`protected=true`）禁止删除

- 新增 Lion ORM `credentials` 实体模型

  在 `lion/credentials/` 新增基于 ent 的凭证实体，覆盖 API Key、密钥对、X.509 证书、许可证、对称密钥等全部字段类型及审计字段。

- 新增凭证种子代码管理

  `admin/well_known_seed_code.go` 新增 `seedCredentialCode()`，将 `CREDENTIAL_CODE_JWT_SIGNING_V1` 映射为 DB code `jwt-signing-v1`，与既有 `DepartmentCode`/`RoleCode`/`AuthProviderCode`/`BootstrapUsername` 风格统一。

- 新增 MCP Server 优雅关闭

  `mcp/mcp.go` 新增 `Server.Close()` 方法，采用 best-effort 策略遍历关闭所有活跃 sessions，收集首个错误但继续关闭剩余 session。

- 新增 MCP 内置 Resources 与 Prompt

  1. `mcp/tools/resources.go` 新增 `RegisterBuiltinResources()`，注册 `grpc-kit://version`、`grpc-kit://openapi-spec/microservice`、`grpc-kit://openapi-spec/admin` 三个内置资源，镜像已公开的 `/version`、`/openapi-spec` HTTP 端点
  2. `mcp/tools/prompts.go` 新增 `RegisterGettingStartedPrompt()` 内置 Prompt `getting_started`，服务名取自 microservice swagger 的 `info.title`

- 新增 AutoBridge HTTP 方法语义标注

  `mcp/tools/bridge.go` 新增 `buildToolAnnotations()`，按 HTTP method 推断 MCP `ToolAnnotations`：GET/HEAD/OPTIONS 只读+幂等、PUT 幂等写入、DELETE 幂等破坏性、POST/PATCH 非幂等；`DestructiveHint` 显式置 `false` 以覆盖 SDK 默认值。

- 新增 MFA 自服务校验

  `admin/rpc_auth_mfa.go` 为 `SetupUserMFA` / `DisableUserMFA` 增加自服务门：用户只能管理自己的 MFA，防止已认证但无角色用户越权操作他人账户。

### Changed

#### grpc-kit/pkg 模块

- 调整 AutoBridge 工具命名格式

  `mcp/tools/bridge.go` 工具名称仅取方法名（snake_case），移除服务名前缀；多 service 场景下 method 重名由 `uniqueToolName` 追加 `_2`/`_3` 兜底。

- 更新凭证种子代码为 JWT 签名版本

  `admin/well_known_seed_code.go` 凭证种子代码调整为 JWT 签名版本，同步调整 `rpc_security_keys_test.go` 相关测试用例。

### Removed

#### grpc-kit/api 模块

- 移除 `CredentialSeedCode` 枚举

  移除旧的 `CredentialSeedCode` 枚举，统一为 `CredentialCode` 种子枚举，与既有种子代码体系风格一致。

#### grpc-kit/pkg 模块

- 移除 MCP 内部运行配置暴露工具

  移除早期 `get_config` 等内部运行配置暴露工具（见 ADR-009），改为通过内置 Resources 镜像公开 HTTP 端点，不引入新安全面。

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
