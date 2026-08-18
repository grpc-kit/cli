package template

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grpc-kit/cli/config"
)

// TestServiceTemplateRendersMCPExtension 防止 MCP 扩展点在模板变更时被遗漏，
// 同时验证生成的 Go 文件均可被 parser 解析。
func TestServiceTemplateRendersMCPExtension(t *testing.T) {
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir(%s): %v", tempDir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	generator, err := New(config.Config{
		Global: config.GlobalConfig{
			Type:           TypeService,
			APIEndpoint:    "api",
			ProductCode:    "demo",
			ShortName:      "echo",
			ReleaseVersion: "0.0.1-test",
			Repository:     "example.com/acme/echo",
			Organization:   "acme",
			Appname:        "echo",
			ServiceCode:    "demo.echo",
			ProtoPackage:   "demo.echo.v1",
			ServiceTitle:   "Echo",
		},
		Template: config.TemplateConfig{Service: config.TemplateService{APIVersion: "v1"}},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	root := filepath.Join(tempDir, "echo")
	assertContains := func(relativePath string, expected ...string) {
		t.Helper()
		content, err := os.ReadFile(filepath.Join(root, relativePath))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", relativePath, err)
		}
		for _, value := range expected {
			if !strings.Contains(string(content), value) {
				t.Errorf("%s does not contain %q", relativePath, value)
			}
		}
	}
	assertContains("handler/register.go", "DO NOT EDIT", "privateMCPHandle()")
	assertContains("handler/private.go", "func (m *Microservice) privateMCPHandle() error", "return fmt.Errorf")
	assertContains("modeler/mcp/registrar.go", "func (r *Registrar) Register", "server is nil")
	assertContains("modeler/mcp/registrar_test.go", "session.CallTool", "session.ReadResource", "session.GetPrompt")
	assertContains("config/app-dev-local.yaml", "aiconnector:", "mcp_server:")
	assertContains("go.mod", "go 1.25.0", "github.com/grpc-kit/pkg v0.4.2", "github.com/modelcontextprotocol/go-sdk v1.7.0")
	assertContains("Makefile", ">> synchronize Go module dependencies", "@${GO} mod tidy")
	assertContains("AGENTS.md", "## Shared Skills", "scripts/skills/skills/generate-release-changelog/SKILL.md")

	// 黑盒接口 E2E 模版框架，详见 test/e2e/README.md
	assertContains("test/e2e/README.md", "//go:build e2e", "make test-e2e")
	assertContains("test/e2e/client/client.go", "//go:build e2e", "func (r *Response) HasRPCCode(")
	assertContains("test/e2e/fixture/main.go", "//go:build e2e", "func Main(m *testing.M)")
	assertContains("test/e2e/fixture/env.go", "//go:build e2e", `defaultServiceCode = "echo.v1.demo"`)
	assertContains("test/e2e/fixture/auth.go", "//go:build e2e", "example.com/acme/echo/test/e2e/client")
	assertContains("test/e2e/demo/testmain_test.go", "//go:build e2e", "fixture.Main(m)")
	assertContains("test/e2e/demo/suite_demo_test.go", "//go:build e2e", "/api/healthz", "/api/demo", "Unauthenticated")
	assertContains("Makefile", "test-e2e:", "vet -tags=e2e")
	assertContains("AGENTS.md", "test/e2e/")

	for _, relativePath := range []string{
		"handler/private.go",
		"modeler/mcp/registrar.go",
		"modeler/mcp/handler.go",
		"test/e2e/client/client.go",
		"test/e2e/fixture/env.go",
		"test/e2e/fixture/main.go",
		"test/e2e/fixture/auth.go",
		"test/e2e/demo/testmain_test.go",
		"test/e2e/demo/suite_demo_test.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relativePath))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", relativePath, err)
		}
		if strings.Contains(string(content), "DO NOT EDIT") {
			t.Errorf("business-editable file %s is incorrectly marked generated", relativePath)
		}
	}

	fset := token.NewFileSet()
	e2eRoot := filepath.Join(root, "test", "e2e")
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		file, parseErr := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if parseErr != nil {
			return parseErr
		}
		if !strings.HasPrefix(path, e2eRoot+string(filepath.Separator)) {
			return nil
		}
		// 黑盒纯度约束：test/e2e/ 只允许依赖标准库与 test/e2e/ 内部包，
		// 禁止 import 本服务的业务包（api/、handler/、internal/、modeler/）。
		for _, spec := range file.Imports {
			importPath := strings.Trim(spec.Path.Value, `"`)
			if !strings.HasPrefix(importPath, "example.com/acme/echo") {
				continue
			}
			if !strings.HasPrefix(importPath, "example.com/acme/echo/test/e2e/") {
				relativePath, _ := filepath.Rel(root, path)
				t.Errorf("e2e file %s violates black-box purity by importing %q", relativePath, importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("parse generated Go files: %v", err)
	}
}
