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

	for _, relativePath := range []string{"handler/private.go", "modeler/mcp/registrar.go", "modeler/mcp/handler.go"} {
		content, err := os.ReadFile(filepath.Join(root, relativePath))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", relativePath, err)
		}
		if strings.Contains(string(content), "DO NOT EDIT") {
			t.Errorf("business-editable file %s is incorrectly marked generated", relativePath)
		}
	}

	fset := token.NewFileSet()
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		_, parseErr := parser.ParseFile(fset, path, nil, parser.AllErrors)
		return parseErr
	})
	if err != nil {
		t.Fatalf("parse generated Go files: %v", err)
	}
}
