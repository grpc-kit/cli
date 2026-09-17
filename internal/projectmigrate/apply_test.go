package projectmigrate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyAndReplanAreIdempotent(t *testing.T) {
	root := writePlanProject(t, true)
	initGitRepository(t, root)
	beforeSnapshot := snapshotProject(t, root)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusReady || len(plan.ApplyBlockers) != 0 {
		t.Fatalf("plan status = %s, blockers = %#v", plan.Status, plan.ApplyBlockers)
	}
	if err := Apply(context.Background(), &plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if plan.Status != StatusApplied {
		t.Fatalf("Apply() status = %s", plan.Status)
	}
	afterSnapshot := snapshotProject(t, root)
	assertOnlyPlannedChanges(t, beforeSnapshot, afterSnapshot, plan.Changes)
	registerBody, err := os.ReadFile(filepath.Join(root, "handler", "register.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(registerBody), "privateMCPHandle") || !strings.Contains(string(registerBody), "StartBackground(ctx)") {
		t.Fatalf("unexpected register.go:\n%s", registerBody)
	}
	info, err := os.Stat(filepath.Join(root, "handler", "register.go"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("register.go mode = %o, want 640", info.Mode().Perm())
	}

	second, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if second.Status != StatusManagedUpToDate || len(second.Changes) != 0 {
		t.Fatalf("second plan status = %s, changes = %v", second.Status, changePaths(second.Changes))
	}
	if !hasDiagnostic(second.ApplyBlockers, "dirty_worktree", "") {
		t.Fatalf("second plan blockers = %#v", second.ApplyBlockers)
	}
	if err := Apply(context.Background(), &second); err != nil {
		t.Fatalf("idempotent Apply() error = %v", err)
	}
}

type snapshotEntry struct {
	body []byte
	mode os.FileMode
}

func snapshotProject(t *testing.T, root string) map[string]snapshotEntry {
	t.Helper()
	snapshot := make(map[string]snapshotEntry)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(relativePath)] = snapshotEntry{body: body, mode: info.Mode()}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func assertOnlyPlannedChanges(t *testing.T, before, after map[string]snapshotEntry, changes []Change) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("file count changed: before=%d after=%d", len(before), len(after))
	}
	changeSet := make(map[string]Change, len(changes))
	for _, change := range changes {
		changeSet[change.Path] = change
	}
	for path, beforeEntry := range before {
		afterEntry, exists := after[path]
		if !exists {
			t.Fatalf("file %s was deleted", path)
		}
		change, planned := changeSet[path]
		if planned {
			if !bytes.Equal(afterEntry.body, change.After) {
				t.Errorf("file %s does not equal planned output", path)
			}
		} else if !bytes.Equal(beforeEntry.body, afterEntry.body) {
			t.Errorf("unplanned file %s changed", path)
		}
		if beforeEntry.mode != afterEntry.mode {
			t.Errorf("file %s mode changed from %v to %v", path, beforeEntry.mode, afterEntry.mode)
		}
	}
}

func TestApplyChangesRollsBackOnFailure(t *testing.T) {
	root := t.TempDir()
	firstPath := filepath.Join(root, "first.txt")
	secondPath := filepath.Join(root, "second.txt")
	writeFile(t, firstPath, []byte("first-before"), 0o600)
	writeFile(t, secondPath, []byte("second-before"), 0o600)
	changes := []Change{
		{Path: "first.txt", Before: []byte("first-before"), After: []byte("first-after"), Mode: 0o600},
		{Path: "second.txt", Before: []byte("second-before"), After: []byte("second-after"), Mode: 0o600},
	}
	calls := 0
	writer := func(path string, body []byte, mode os.FileMode) error {
		calls++
		if calls == 2 {
			return errors.New("injected failure")
		}
		return os.WriteFile(path, body, mode)
	}
	if err := applyChanges(root, changes, writer); err == nil || !strings.Contains(err.Error(), "injected failure") {
		t.Fatalf("applyChanges() error = %v", err)
	}
	first, _ := os.ReadFile(firstPath)
	second, _ := os.ReadFile(secondPath)
	if string(first) != "first-before" || string(second) != "second-before" {
		t.Fatalf("rollback result first=%q second=%q", first, second)
	}
}

func TestVerifyChangeInputDetectsMutation(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "managed.go")
	writeFile(t, path, []byte("before"), 0o600)
	change := Change{Path: "managed.go", Before: []byte("before"), Mode: 0o600}
	change.BeforeSHA256 = [32]byte{1}
	if err := verifyChangeInput(root, change); err == nil || !strings.Contains(err.Error(), "content changed") {
		t.Fatalf("verifyChangeInput() error = %v", err)
	}
}

func TestApplyRewritesLegacyGenerateScript(t *testing.T) {
	root := writePlanProject(t, true)
	writeFile(t, filepath.Join(root, "scripts", "generate.sh"), legacyGenerateScriptFixture, 0o755)
	initGitRepository(t, root)
	beforeSnapshot := snapshotProject(t, root)

	plan, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusReady || len(plan.ApplyBlockers) != 0 {
		t.Fatalf("plan status = %s, blockers = %#v", plan.Status, plan.ApplyBlockers)
	}
	if got := changePaths(plan.Changes); len(got) != 7 {
		t.Fatalf("change paths = %v, want 7 managed files", got)
	}
	if err := Apply(context.Background(), &plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if plan.Status != StatusApplied {
		t.Fatalf("Apply() status = %s", plan.Status)
	}
	afterSnapshot := snapshotProject(t, root)
	assertOnlyPlannedChanges(t, beforeSnapshot, afterSnapshot, plan.Changes)
	scriptBody, err := os.ReadFile(filepath.Join(root, "scripts", "generate.sh"))
	if err != nil {
		t.Fatal(err)
	}
	expected, err := renderGenerateScriptTarget("0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(scriptBody, expected) {
		t.Fatal("generate.sh was not rewritten to the rendered target content")
	}
	marker, err := ParseMarkerLine(secondLine(scriptBody))
	if err != nil || marker.Version != "0.4.0" || marker.CommentPrefix != "#" {
		t.Fatalf("migrated generate.sh marker = %#v, want script marker 0.4.0", marker)
	}
	info, err := os.Stat(filepath.Join(root, "scripts", "generate.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("generate.sh mode = %o, want 755", info.Mode().Perm())
	}

	second, err := BuildPlan(root, "0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if second.Status != StatusManagedUpToDate || len(second.Changes) != 0 {
		t.Fatalf("second plan status = %s, changes = %v", second.Status, changePaths(second.Changes))
	}
	if err := Apply(context.Background(), &second); err != nil {
		t.Fatalf("idempotent Apply() error = %v", err)
	}
}

func initGitRepository(t *testing.T, root string) {
	t.Helper()
	commands := [][]string{
		{"init"},
		{"add", "."},
		{"-c", "user.name=grpc-kit-test", "-c", "user.email=test@grpc-kit.invalid", "commit", "-m", "fixture"},
	}
	for _, args := range commands {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, output)
		}
	}
}
