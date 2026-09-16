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
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type atomicWriteFunc func(path string, body []byte, mode fs.FileMode) error

// Apply writes exactly the changes in plan after revalidating Git state and
// every input digest. A failure rolls back files written by this invocation.
func Apply(ctx context.Context, plan *Plan) error {
	if plan == nil {
		return fmt.Errorf("plan is nil")
	}
	if plan.Blocked() {
		return fmt.Errorf("migration plan is blocked by conflicts")
	}
	if len(plan.Changes) == 0 {
		plan.Status = StatusManagedUpToDate
		return nil
	}
	if _, err := ValidateTargetCLIVersion(plan.TargetCLIVersion); err != nil {
		return fmt.Errorf("refuse apply: %w", err)
	}
	if err := ensureGitWorktree(ctx, plan.Project.Root); err != nil {
		return err
	}
	status, err := gitOutput(ctx, plan.Project.Root, "status", "--porcelain", "--untracked-files=normal", "--", ".")
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(status)) != 0 {
		return fmt.Errorf("project worktree is not clean")
	}
	for _, change := range plan.Changes {
		if err := ensureGitTracked(ctx, plan.Project.Root, change.Path); err != nil {
			return err
		}
		if err := verifyChangeInput(plan.Project.Root, change); err != nil {
			return err
		}
	}
	if err := applyChanges(plan.Project.Root, plan.Changes, writeFileAtomic); err != nil {
		return err
	}
	plan.Status = StatusApplied
	return nil
}

func applyChanges(root string, changes []Change, writer atomicWriteFunc) error {
	applied := make([]Change, 0, len(changes))
	for _, change := range changes {
		absolutePath, err := projectFilePath(root, change.Path)
		if err != nil {
			return rollbackChanges(root, applied, writer, err)
		}
		if err := writer(absolutePath, change.After, change.Mode); err != nil {
			return rollbackChanges(root, applied, writer, fmt.Errorf("write %s: %w", change.Path, err))
		}
		applied = append(applied, change)
		written, err := readRegularFile(absolutePath, maxManagedFileSize)
		if err != nil || !bytes.Equal(written, change.After) {
			if err == nil {
				err = fmt.Errorf("written content did not match the plan")
			}
			return rollbackChanges(root, applied, writer, fmt.Errorf("verify written file %s: %w", change.Path, err))
		}
	}
	return nil
}

func rollbackChanges(root string, applied []Change, writer atomicWriteFunc, cause error) error {
	var rollbackErrs []error
	for index := len(applied) - 1; index >= 0; index-- {
		change := applied[index]
		absolutePath, err := projectFilePath(root, change.Path)
		if err == nil {
			err = writer(absolutePath, change.Before, change.Mode)
		}
		if err != nil {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("rollback %s: %w", change.Path, err))
		}
	}
	if len(rollbackErrs) == 0 {
		return cause
	}
	return errors.Join(append([]error{cause}, rollbackErrs...)...)
}

func verifyChangeInput(root string, change Change) error {
	absolutePath, err := projectFilePath(root, change.Path)
	if err != nil {
		return err
	}
	body, err := readRegularFile(absolutePath, maxManagedFileSize)
	if err != nil {
		return fmt.Errorf("revalidate %s: %w", change.Path, err)
	}
	if sha256.Sum256(body) != change.BeforeSHA256 {
		return fmt.Errorf("revalidate %s: content changed after planning", change.Path)
	}
	info, err := os.Lstat(absolutePath)
	if err != nil {
		return err
	}
	if info.Mode() != change.Mode {
		return fmt.Errorf("revalidate %s: file mode changed after planning", change.Path)
	}
	return nil
}

func projectFilePath(root, relativePath string) (string, error) {
	if relativePath == "" || filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("unsafe project path %q", relativePath)
	}
	cleanRelative := filepath.Clean(filepath.FromSlash(relativePath))
	if cleanRelative == "." || cleanRelative == ".." || len(cleanRelative) >= 3 && cleanRelative[:3] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("unsafe project path %q", relativePath)
	}
	absolutePath := filepath.Join(root, cleanRelative)
	withinRoot, err := filepath.Rel(root, absolutePath)
	if err != nil || withinRoot == ".." || len(withinRoot) >= 3 && withinRoot[:3] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("unsafe project path %q", relativePath)
	}
	return absolutePath, nil
}

func writeFileAtomic(path string, body []byte, mode fs.FileMode) (returnErr error) {
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".grpc-kit-migrate-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	closed := false
	defer func() {
		if !closed {
			if closeErr := temporary.Close(); returnErr == nil && closeErr != nil {
				returnErr = closeErr
			}
		}
		if removeErr := os.Remove(temporaryPath); returnErr == nil && removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			returnErr = removeErr
		}
	}()
	if err := temporary.Chmod(mode.Perm()); err != nil {
		return err
	}
	if _, err := temporary.Write(body); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	closed = true
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	return nil
}
