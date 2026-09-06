package template

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/grpc-kit/cli/config"
)

// TestServiceTemplateRendersExtensions 防止服务扩展点在模板变更时被遗漏，
// 同时验证生成的 Go 文件均可被 parser 解析。
func TestServiceTemplateRendersExtensions(t *testing.T) {
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
	assertContains("AGENTS.md", "## Service Skills", ".agents/skills/add-api-domain/SKILL.md")
	assertContains(".agents/skills/README.md", "# 服务技能", "模板打包限制")
	assertContains(".agents/skills/add-api-domain/SKILL.md", "name: add-api-domain", "make test-e2e")
	for _, relativePath := range []string{"AGENTS.md", "Makefile", filepath.Join("scripts", "env"), filepath.Join("test", "e2e", "README.md")} {
		if _, err := os.Stat(filepath.Join(root, relativePath)); err != nil {
			t.Errorf("skill repository reference %s is unavailable: %v", relativePath, err)
		}
	}
	skillContent, err := os.ReadFile(filepath.Join(root, ".agents", "skills", "add-api-domain", "SKILL.md"))
	if err != nil {
		t.Fatalf("ReadFile(SKILL.md): %v", err)
	}
	for _, forbidden := range []string{"template/service/v1", "example.com/acme/echo"} {
		if strings.Contains(string(skillContent), forbidden) {
			t.Errorf("SKILL.md contains template-specific path %q", forbidden)
		}
	}

	claudeContent, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("ReadFile(CLAUDE.md): %v", err)
	}
	if got, want := string(claudeContent), "@AGENTS.md\n"; got != want {
		t.Errorf("CLAUDE.md = %q, want %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); err != nil {
		t.Errorf("CLAUDE.md import target AGENTS.md is unavailable: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".claude", "skills")); !os.IsNotExist(err) {
		t.Errorf("unexpected .claude/skills bridge: %v", err)
	}
	for _, relativePath := range []string{
		"CLAUDE.md.tmpl",
		filepath.Join(".agents", "skills", "add-api-domain", "SKILL.md.tmpl"),
	} {
		if _, err := os.Stat(filepath.Join(root, relativePath)); !os.IsNotExist(err) {
			t.Errorf("unexpected template suffix in generated service %s: %v", relativePath, err)
		}
	}
	makefile, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		t.Fatalf("ReadFile(Makefile): %v", err)
	}
	if strings.Contains(string(makefile), "agents-sync:") {
		t.Error("unexpected agents-sync target")
	}

	validateGeneratedSkill(t, filepath.Join(root, ".agents", "skills", "add-api-domain"))
	compareTemplateSkillAssets(t, previousDir)

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
		filepath.Join(".agents", "skills", "README.md"),
		filepath.Join(".agents", "skills", "add-api-domain", "SKILL.md"),
		"CLAUDE.md",
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
		if strings.HasPrefix(string(content), "// Code generated by") {
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

func validateGeneratedSkill(t *testing.T, skillDir string) {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		t.Fatalf("ReadFile(SKILL.md): %v", err)
	}
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 || parts[0] != "" {
		t.Fatal("SKILL.md must start with YAML frontmatter")
	}
	fields := make(map[string]string)
	for _, line := range strings.Split(parts[1], "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok {
			fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	name := fields["name"]
	description := fields["description"]
	if name != filepath.Base(skillDir) {
		t.Errorf("skill name %q does not match directory %q", name, filepath.Base(skillDir))
	}
	if valid := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`).MatchString(name); !valid || len(name) > 64 {
		t.Errorf("invalid skill name %q", name)
	}
	if description == "" || len(description) > 1024 {
		t.Errorf("invalid skill description length %d", len(description))
	}
	body := strings.TrimSpace(parts[2])
	if body == "" || strings.Contains(body, "TODO") {
		t.Error("skill body is empty or contains an unfinished placeholder")
	}
	if strings.HasPrefix(string(content), "// Code generated by") {
		t.Error("service skill is incorrectly marked generated")
	}
}

func compareTemplateSkillAssets(t *testing.T, templateDir string) {
	t.Helper()

	const root = "service/.agents/skills"
	diskFiles := make([]string, 0)
	err := filepath.WalkDir(filepath.Join(templateDir, "service", ".agents", "skills"), func(filePath string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filePath != filepath.Join(templateDir, "service", ".agents", "skills") &&
			(strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "_")) {
			return &fs.PathError{Op: "embed", Path: filePath, Err: fs.ErrInvalid}
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return &fs.PathError{Op: "embed", Path: filePath, Err: fs.ErrInvalid}
		}
		if entry.Name() == "go.mod" {
			return &fs.PathError{Op: "embed", Path: filePath, Err: fs.ErrInvalid}
		}
		if entry.IsDir() {
			return nil
		}
		relativePath, err := filepath.Rel(templateDir, filePath)
		if err != nil {
			return err
		}
		diskFiles = append(diskFiles, filepath.ToSlash(relativePath))
		return nil
	})
	if err != nil {
		t.Fatalf("walk disk skill template: %v", err)
	}

	embeddedFiles := make([]string, 0)
	err = fs.WalkDir(Assets, root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			embeddedFiles = append(embeddedFiles, filePath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk embedded skill template: %v", err)
	}
	sort.Strings(diskFiles)
	sort.Strings(embeddedFiles)
	if strings.Join(diskFiles, "\x00") != strings.Join(embeddedFiles, "\x00") {
		t.Errorf("embedded skill files = %v, disk skill files = %v", embeddedFiles, diskFiles)
	}
}
