package projectmigrate

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestBuildPlanForSpiderShape(t *testing.T) {
	root := writePlanProject(t, false)
	writeManagedStub(t, root, "cmd/server/main.go", "0.3.9-beta.1", "main")
	writeManagedStub(t, root, "handler/microservice.go", "0.3.9-beta.1", "handler")
	writeManagedStub(t, root, "handler/rpc_internal.go", "0.3.8-beta.1", "handler")
	writeManagedStub(t, root, "handler/shutdown.go", "0.3.9-beta.1", "handler")
	writeFile(t, filepath.Join(root, "handler", "register.go"), []byte("package handler\n// user managed\n"), 0o640)
	writeFile(t, filepath.Join(root, legacyPublicEmbedPath), []byte(legacyPublicEmbedFixture), 0o640)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	if plan.Status != StatusReady || plan.Blocked() {
		t.Fatalf("BuildPlan() status = %s, conflicts = %#v", plan.Status, plan.Conflicts)
	}
	wantChanges := []string{
		"cmd/server/main.go",
		"handler/microservice.go",
		"handler/rpc_internal.go",
		"handler/shutdown.go",
		"scripts/env",
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, wantChanges) {
		t.Fatalf("change paths = %v, want %v", got, wantChanges)
	}
	if !hasDiagnostic(plan.Warnings, "known_legacy_marker", legacyPublicEmbedPath) || !hasDiagnostic(plan.Warnings, "mixed_versions", "") {
		t.Fatalf("warnings = %#v", plan.Warnings)
	}
	if slices.Contains(changePaths(plan.Changes), "handler/register.go") {
		t.Fatal("user-managed register.go entered Changes")
	}
	registerBody, err := os.ReadFile(filepath.Join(root, "handler", "register.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(registerBody) != "package handler\n// user managed\n" {
		t.Fatal("preview modified register.go")
	}
}

func TestBuildPlanRejectsUnmappedManagedFile(t *testing.T) {
	root := writePlanProject(t, false)
	writeManagedStub(t, root, "old/generated.go", "0.3.8", "old")

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusConflict || !hasDiagnostic(plan.Conflicts, "unmapped_managed_file", "old/generated.go") {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestBuildPlanRejectsModifiedLegacyMarkerFixture(t *testing.T) {
	root := writePlanProject(t, false)
	writeFile(t, filepath.Join(root, legacyPublicEmbedPath), []byte(legacyPublicEmbedFixture+"// changed\n"), 0o640)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusConflict || !hasDiagnostic(plan.Conflicts, "invalid_marker", legacyPublicEmbedPath) {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestBuildPlanDoesNotCreateMissingAssets(t *testing.T) {
	root := writePlanProject(t, false)
	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, []string{"scripts/env"}) {
		t.Fatalf("change paths = %v", got)
	}
	if _, err := os.Stat(filepath.Join(root, "handler", "microservice.go")); !os.IsNotExist(err) {
		t.Fatalf("preview created a missing asset: %v", err)
	}
}

func TestBuildPlanRewritesLegacyGenerateScript(t *testing.T) {
	root := writePlanProject(t, false)
	writeFile(t, filepath.Join(root, "scripts", "generate.sh"), legacyGenerateScriptFixture, 0o755)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusReady || plan.Blocked() {
		t.Fatalf("plan status = %s, conflicts = %#v", plan.Status, plan.Conflicts)
	}
	wantChanges := []string{"scripts/env", "scripts/generate.sh"}
	if got := changePaths(plan.Changes); !slices.Equal(got, wantChanges) {
		t.Fatalf("change paths = %v, want %v", got, wantChanges)
	}
	target, err := renderGenerateScriptTarget("0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range plan.Changes {
		if change.Path != generateScriptPath {
			continue
		}
		if !bytes.Equal(change.After, target) {
			t.Fatal("generate.sh rewrite does not equal the rendered target")
		}
		marker, err := ParseMarkerLine(secondLine(change.After))
		if err != nil || marker.Version != "0.4.0" || marker.CommentPrefix != "#" {
			t.Fatalf("generate.sh rewrite marker = %#v, want script marker 0.4.0", marker)
		}
	}
	if hasDiagnosticCode(plan.ManualActions, "generate_script_modified") {
		t.Fatalf("planned rewrite also produced a manual action: %#v", plan.ManualActions)
	}
}

func TestBuildPlanSkipsCurrentGenerateScript(t *testing.T) {
	root := writePlanProject(t, false)
	scriptPath := filepath.Join(root, "scripts", "generate.sh")
	target, err := renderGenerateScriptTarget("0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, scriptPath, target, 0o755)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, []string{"scripts/env"}) {
		t.Fatalf("change paths = %v, want only scripts/env", got)
	}
	if hasDiagnosticCode(plan.ManualActions, "generate_script_modified") {
		t.Fatalf("current content produced a manual action: %#v", plan.ManualActions)
	}
	body, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(body, target) {
		t.Fatal("preview modified an up-to-date generate.sh")
	}
}

func TestBuildPlanLeavesModifiedMarkedGenerateScriptToManualAction(t *testing.T) {
	root := writePlanProject(t, false)
	target, err := renderGenerateScriptTarget("0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "scripts", "generate.sh"), append(bytes.Clone(target), []byte("# local edit\n")...), 0o755)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusReady || plan.Blocked() {
		t.Fatalf("modified marked generate.sh must not block the plan: status = %s, conflicts = %#v", plan.Status, plan.Conflicts)
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, []string{"scripts/env"}) {
		t.Fatalf("change paths = %v, want only scripts/env", got)
	}
	if !hasDiagnostic(plan.ManualActions, "generate_script_modified", "scripts/generate.sh") {
		t.Fatalf("manual actions = %#v, want generate_script_modified", plan.ManualActions)
	}
}

func TestBuildPlanLeavesModifiedGenerateScriptToManualAction(t *testing.T) {
	root := writePlanProject(t, false)
	writeFile(t, filepath.Join(root, "scripts", "generate.sh"), append(bytes.Clone(legacyGenerateScriptFixture), []byte("# local edit\n")...), 0o755)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusReady || plan.Blocked() {
		t.Fatalf("modified generate.sh must not block the plan: status = %s, conflicts = %#v", plan.Status, plan.Conflicts)
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, []string{"scripts/env"}) {
		t.Fatalf("change paths = %v, want only scripts/env", got)
	}
	if !hasDiagnostic(plan.ManualActions, "generate_script_modified", "scripts/generate.sh") {
		t.Fatalf("manual actions = %#v, want generate_script_modified", plan.ManualActions)
	}
}

func TestBuildPlanIgnoresMissingGenerateScript(t *testing.T) {
	root := writePlanProject(t, false)
	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if got := changePaths(plan.Changes); !slices.Equal(got, []string{"scripts/env"}) {
		t.Fatalf("change paths = %v, want only scripts/env", got)
	}
	if hasDiagnosticCode(plan.ManualActions, "generate_script_modified") {
		t.Fatalf("missing file produced a manual action: %#v", plan.ManualActions)
	}
}

func TestBuildPlanConflictsWhenGenerateScriptPathIsUnsafe(t *testing.T) {
	t.Run("directory", func(t *testing.T) {
		root := writePlanProject(t, false)
		if err := os.Mkdir(filepath.Join(root, "scripts", "generate.sh"), 0o755); err != nil {
			t.Fatal(err)
		}
		plan, err := BuildPlan(root, "0.4.0")
		if err != nil {
			t.Fatal(err)
		}
		if plan.Status != StatusConflict || !hasDiagnostic(plan.Conflicts, "unsafe_file", "scripts/generate.sh") {
			t.Fatalf("plan = %#v, want unsafe_file conflict", plan)
		}
	})
	t.Run("symlink", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlink behavior requires elevated privileges on some Windows hosts")
		}
		root := writePlanProject(t, false)
		target := filepath.Join(root, "actual-generate.sh")
		if err := os.WriteFile(target, legacyGenerateScriptFixture, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(root, "scripts", "generate.sh")); err != nil {
			t.Fatal(err)
		}
		plan, err := BuildPlan(root, "0.4.0")
		if err != nil {
			t.Fatal(err)
		}
		if plan.Status != StatusConflict || !hasDiagnostic(plan.Conflicts, "unsafe_file", "scripts/generate.sh") {
			t.Fatalf("plan = %#v, want unsafe_file conflict", plan)
		}
	})
}

func TestWritePlanIncludesUnifiedDiff(t *testing.T) {
	root := writePlanProject(t, false)
	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := WritePlan(&output, plan); err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"Status: ready", "--- a/scripts/env", "+++ b/scripts/env", "CLI_VERSION=0.4.0"} {
		if !strings.Contains(output.String(), value) {
			t.Errorf("output does not contain %q:\n%s", value, output.String())
		}
	}
}

func writePlanProject(t *testing.T, allManaged bool) string {
	t.Helper()
	root := writeProject(t, validEnvironment, "module git.lmq.io/kaopuvm/spider\n\ngo 1.25.0\n")
	if allManaged {
		writeManagedStub(t, root, "cmd/server/main.go", "0.3.9-beta.1", "main")
		writeManagedStub(t, root, "handler/microservice.go", "0.3.9-beta.1", "handler")
		writeManagedStub(t, root, "handler/register.go", "0.3.9-beta.1", "handler")
		writeManagedStub(t, root, "handler/rpc_internal.go", "0.3.8-beta.1", "handler")
		writeManagedStub(t, root, "handler/shutdown.go", "0.3.9-beta.1", "handler")
	}
	return root
}

func writeManagedStub(t *testing.T, root, relativePath, version, packageName string) {
	t.Helper()
	body := []byte(`// Code generated by "grpc-kit-cli/` + version + `". DO NOT EDIT.

package ` + packageName + "\n")
	writeFile(t, filepath.Join(root, filepath.FromSlash(relativePath)), body, 0o640)
}

func writeFile(t *testing.T, path string, body []byte, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, mode); err != nil {
		t.Fatal(err)
	}
}

func changePaths(changes []Change) []string {
	paths := make([]string, 0, len(changes))
	for _, change := range changes {
		paths = append(paths, change.Path)
	}
	return paths
}

func hasDiagnostic(diagnostics []Diagnostic, code, path string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Path == path {
			return true
		}
	}
	return false
}
