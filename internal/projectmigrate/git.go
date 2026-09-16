// Copyright © 2026 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package projectmigrate

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func inspectApplyBlockers(root, targetVersion string) []Diagnostic {
	var blockers []Diagnostic
	if _, err := ValidateTargetCLIVersion(targetVersion); err != nil {
		blockers = append(blockers, Diagnostic{Code: "release_version", Message: err.Error()})
	}
	if err := ensureGitWorktree(context.Background(), root); err != nil {
		blockers = append(blockers, Diagnostic{Code: "git_worktree", Message: err.Error()})
		return blockers
	}
	dirty, err := gitOutput(context.Background(), root, "status", "--porcelain", "--untracked-files=normal", "--", ".")
	if err != nil {
		blockers = append(blockers, Diagnostic{Code: "git_status", Message: err.Error()})
	} else if len(bytes.TrimSpace(dirty)) != 0 {
		blockers = append(blockers, Diagnostic{Code: "dirty_worktree", Message: "project worktree is not clean"})
	}
	return blockers
}

func ensureGitWorktree(ctx context.Context, root string) error {
	if _, err := gitOutput(ctx, root, "rev-parse", "--show-toplevel"); err != nil {
		return fmt.Errorf("project is not in a Git worktree: %w", err)
	}
	return nil
}

func ensureGitTracked(ctx context.Context, root, relativePath string) error {
	if _, err := gitOutput(ctx, root, "ls-files", "--error-unmatch", "--", relativePath); err != nil {
		return fmt.Errorf("%s is not tracked by Git: %w", relativePath, err)
	}
	return nil
}

func gitOutput(ctx context.Context, root string, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-C", root}, args...)
	cmd := exec.CommandContext(ctx, "git", commandArgs...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git %v: %s", args, bytes.TrimSpace(exitErr.Stderr))
		}
		return nil, err
	}
	return output, nil
}
