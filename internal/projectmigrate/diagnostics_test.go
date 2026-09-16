package projectmigrate

import (
	"strings"
	"testing"
)

func TestScanManualActions(t *testing.T) {
	goMod := `module git.lmq.io/kaopuvm/spider

go 1.25.0

require (
	github.com/grpc-kit/pkg v0.4.3
	github.com/sirupsen/logrus v1.9.3
)
`
	root := writeProject(t, validEnvironment, goMod)
	writeFile(t, root+"/user.go", []byte(`package spider

import "github.com/sirupsen/logrus"

type IndependentCfg struct { logger *logrus.Entry }
type Option func(*IndependentCfg)

func (c *IndependentCfg) Init(opts ...Option) error { return nil }

func migrate(ctx, lc, base, server, handler, status, mux, assets any, logger *logrus.Entry) {
	logger.Warnf("value: %v", 1)
	logger.WithFields(logrus.Fields{"key": "value"})
	logrus.WithField("package", "call")
	_ = logger.Data
	lc.Init()
	base.Deregister()
	base.HTTPHandlerFrontend(mux, assets)
	server.StartBackground()
	handler.NewMicroservice(lc)
	status.WithLogger(logger, "%v", nil)
}
`), 0o600)
	writeFile(t, root+"/registry.go", []byte(`package spider

import "github.com/grpc-kit/pkg/sd"

func register(conn, ttl any) {
	sd.Register(conn, "name", "address", "value", ttl)
}
`), 0o600)
	writeFile(t, root+"/managed.go", []byte(`package spider
import "github.com/sirupsen/logrus"
var managed *logrus.Entry
`), 0o600)

	project, err := LoadProject(root)
	if err != nil {
		t.Fatal(err)
	}
	actions, err := ScanManualActions(project, map[string]struct{}{"managed.go": {}})
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"go_version", "pkg_version", "logrus_dependency", "logrus_import", "logrus_type", "logrus_custom", "printf_logger", "context_signature", "context_call", "errs_with_logger"} {
		if !hasDiagnosticCode(actions, code) {
			t.Errorf("actions do not contain %q: %#v", code, actions)
		}
	}
	for _, messagePart := range []string{"WithField", "Entry.Data", "sd.Register"} {
		if !hasDiagnosticMessage(actions, messagePart) {
			t.Errorf("actions do not contain message %q: %#v", messagePart, actions)
		}
	}
	for _, action := range actions {
		if action.Path == "managed.go" {
			t.Fatalf("managed change produced a manual action: %#v", action)
		}
		if action.Path != "go.mod" && action.Line == 0 {
			t.Errorf("source action has no line: %#v", action)
		}
	}
}

func hasDiagnosticMessage(diagnostics []Diagnostic, part string) bool {
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic.Message, part) {
			return true
		}
	}
	return false
}

func TestWritePlanIncludesManualActionLocation(t *testing.T) {
	plan := Plan{
		Status:           StatusManagedUpToDate,
		TargetCLIVersion: "0.4.0",
		TargetPkgVersion: "v0.5.0",
		ManualActions: []Diagnostic{
			{Code: "context_call", Path: "handler/register.go", Line: 38, Message: "pass context"},
		},
	}
	var output strings.Builder
	if err := WritePlan(&output, plan); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Manual action [context_call] handler/register.go:38: pass context") {
		t.Fatalf("output = %q", output.String())
	}
	if !strings.Contains(output.String(), "Project compatibility: not verified") {
		t.Fatalf("output does not explain compatibility status: %q", output.String())
	}
}

func hasDiagnosticCode(diagnostics []Diagnostic, code string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}
