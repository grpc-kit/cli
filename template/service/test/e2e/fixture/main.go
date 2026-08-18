//go:build e2e

// Main 供各业务域包的 TestMain 复用：探活 guard 与进程退出码收敛于此。
// 各域包（test/e2e/<domain>）一行接入：func TestMain(m *testing.M) { fixture.Main(m) }。
package fixture

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

// Main 是各业务域包 TestMain 的统一入口：探活 guard + m.Run。
func Main(m *testing.M) {
	env := Init()

	// 前置 guard：确认被测服务可达，失败时给出可操作的提示而不是一排连接错误。
	if err := probe(env); err != nil {
		fmt.Fprintf(os.Stderr, "e2e: 探活失败 %s（%v）；请先在另一终端执行 make run\n",
			env.BaseURL, err)
		os.Exit(1)
	}
	fmt.Printf("e2e: 目标服务 %s，serviceCode=%s，runID=%s\n",
		env.BaseURL, env.ServiceCode, env.RunID)

	os.Exit(m.Run())
}

// probe 用免认证的健康检查接口确认服务在线。
// /api/healthz 对应 HealthCheck 方法，它在 config/app-dev-local.yaml 的
// security.authentication.insecure_rpcs 列表内，匿名可访问；
// 该方法要求 service 入参等于本服务的 service_code，否则返回 404 NotFound。
func probe(env *Env) error {
	target := env.BaseURL + "/api/healthz?service=" + url.QueryEscape(env.ServiceCode)

	probeClient := &http.Client{Timeout: 3 * time.Second}
	resp, err := probeClient.Get(target)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// 常见原因：E2E_SERVICE_CODE 与服务配置的 service_code 不一致，
		// 或该端口上跑的是另一个服务（此时网关会返回 501 Unimplemented）。
		return fmt.Errorf("GET /api/healthz?service=%s 返回 HTTP %d（期望 200）；"+
			"请确认端口未被其他服务占用，且 service_code 一致",
			env.ServiceCode, resp.StatusCode)
	}
	return nil
}
