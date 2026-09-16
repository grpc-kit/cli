# test/ 目录规划规范

参考基准：Kubernetes 测试组织方式 + Go 官方测试惯例。

## 1. 测试类型归属

| 测试类型 | 位置 | 隔离方式 | 运行入口 |
| --- | --- | --- | --- |
| 白盒/灰盒单测 | 与源码同目录（colocated，`*_test.go`） | 无 build tag | `make test`（`go test ./...`） |
| 依赖真实外部后端的集成测试 | 与源码同目录 | 环境变量门控 + `t.Skip` | `make test` + 设置环境变量 |
| 黑盒 E2E | `test/e2e/` | `//go:build e2e` | `make test-e2e`（前置 `make run`） |

**黑盒纯度约束**：`test/e2e/` 下代码只允许 import 标准库、`test/e2e/client`、
`test/e2e/fixture`，禁止导入本服务的业务包（`api/`、`handler/`、`internal/`、`modeler/`）
—— 保证测的是对外 HTTP 契约而非内部实现。

这条约束还带来一个实用副作用：`go vet -tags=e2e ./test/e2e/...` 只依赖标准库，
不需要先跑 `make generate`，因此在没有 protoc 工具链的环境里也能做编译检查。

## 2. 目录结构（按业务域分包）

```text
test/
└── e2e/
    ├── client/     # 薄 HTTP 客户端：鉴权注入、统一错误体解码（google.rpc.Status）
    ├── fixture/    # 环境引导：env（BaseURL/ServiceCode/runID/前缀）、
    │               # main（Main(m) 探活 guard）、auth（调用主体与资源工厂挂载点）
    └── demo/       # 内置示例接口域：HealthCheck 与 Demo 的全部 HTTP 绑定（D1–D8）
```

- **目录创建时机**：域目录随该域首个 suite 落地时创建，不预建空目录。
- **`demo/` 的处置**：它测的是模版内置的示例接口。删除 `handler/rpc_demo.go` 与
  `api/**/demo.proto` 时，请一并删除 `test/e2e/demo/`。
- **Go 语言约束**（域分包的既定代价与对策）：
  1. `TestMain` 是 per-package 的：各域包 `testmain_test.go` 一行
     `func TestMain(m *testing.M) { fixture.Main(m) }` 接入统一引导；
  2. 每个域包是独立测试进程：runID 按进程生成（隔离性更好）；
  3. Makefile `test-e2e` 固定 `-p 1` 串行化包执行，防共享环境互踩。

## 3. 命名与隔离约定

- suite 文件：`suite_<domain>_<topic>_test.go`；
- 包名 = 域名（如 `package demo`）；
- 文件头注释引用被测文档章节与用例 ID 表，便于 grep 对照；
- Test 函数 `Test<Domain><Behavior>`，子测试名以用例 ID 开头（如 `"D5 ..."`）；
- helper 归属：域内私有 helper 留在 suite 文件；跨域共享的一律 export 进 `fixture/`；
- 有状态资源一律「创建即注册 `t.Cleanup`」，清理失败只 `t.Logf`，不中断其他清理；
- 测试数据带 `env.Prefix()`（含 runID）前缀，避免并发与历史残留互踩；
- `test/e2e/` 下所有 `.go` 文件（含 `client`/`fixture` 的非 `_test` 文件）统一
  `//go:build e2e`；
- 单测、集成测试留在源码旁，禁止迁入 `test/`。

## 4. 运行方式

```shell
make run                 # 终端 1：本地起服务（config/app-dev-local.yaml）
make test-e2e            # 终端 2：全量 E2E（-p 1 串行）

# 按域/用例过滤
go test -tags=e2e -p 1 ./test/e2e/demo/... -run TestDemo -v
```

`.gitignore` 忽略了 `config/app-dev-*.yaml`，因此 CI 或全新环境里没有那份配置。
E2E 不自带配置，被测服务的运行配置由环境提供（挂载、密钥管理或前置阶段下发），
CI 里通过 `E2E_CONFIG_FILE` 指定要读哪一份：

```shell
make build
./build/service -c config/app-dev-local.yaml     # 终端 1，或 CI 里的 ${E2E_CONFIG_FILE}
make test-e2e                                    # 终端 2
```

用例只依赖该配置的三处内容——`services.http_address`、`services.service_code`、
`security.authentication.http_users`。与默认值不一致时用下面的环境变量覆盖，不要改用例：

| 变量 | 默认值 | 对应配置项 |
| --- | --- | --- |
| `E2E_BASE_URL` | `http://127.0.0.1:8080` | `services.http_address` |
| `E2E_SERVICE_CODE` | 由模版生成 | `services.service_code`（健康检查的必填入参） |
| `E2E_BASIC_USERNAME` | `user1` | `security.authentication.http_users[].username` |
| `E2E_BASIC_PASSWORD` | `grpc-kit-cli` | `security.authentication.http_users[].password` |

若项目改用 OIDC / Bearer 认证，请在 `fixture/auth.go` 的 `RequireAuthed` 里换取
access_token 并用 `client.WithToken` 构造，把失败收敛成一次 `t.Fatalf`。

常见故障：

| 现象 | 原因与处理 |
| --- | --- |
| `e2e: 探活失败 ...: connection refused` | 服务没起来，先 `make run` |
| 探活返回 `501 Unimplemented` | 端口被另一个 grpc-kit 服务占用（网关收到请求但后端没有本服务），换端口或用 `E2E_BASE_URL` 指向正确实例 |
| 探活返回 `404 NotFound` | `E2E_SERVICE_CODE` 与服务配置的 `service_code` 不一致 |
| 用例返回 `401 Unauthenticated` | 配置里的 `http_users` 与 `E2E_BASIC_USERNAME` / `E2E_BASIC_PASSWORD` 不匹配 |

## 5. testdata 约定

E2E 静态数据（golden file、OpenAPI 样本等）放 `test/e2e/<域>/testdata/`，
Go 工具链自动忽略该目录名。
