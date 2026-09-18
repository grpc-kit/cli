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
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const maxManagedFileSize = 4 << 20

type Status string

const (
	StatusManagedUpToDate Status = "managed_up_to_date"
	StatusReady           Status = "ready"
	StatusApplied         Status = "applied"
	StatusConflict        Status = "conflict"
	StatusUnsupported     Status = "unsupported"
)

type Diagnostic struct {
	Code    string
	Path    string
	Line    int
	Message string
}

type Change struct {
	Path         string
	Before       []byte
	After        []byte
	Mode         fs.FileMode
	BeforeSHA256 [sha256.Size]byte
}

type Plan struct {
	Project          Project
	FromCLIVersions  []string
	TargetCLIVersion string
	TargetPkgVersion string
	Status           Status
	Changes          []Change
	Warnings         []Diagnostic
	ApplyBlockers    []Diagnostic
	ManualActions    []Diagnostic
	Conflicts        []Diagnostic
}

func (p Plan) Blocked() bool {
	return p.Status == StatusConflict || p.Status == StatusUnsupported || len(p.Conflicts) != 0
}

// BuildPlan computes the complete write set without modifying the project.
func BuildPlan(projectPath, targetCLIVersion string) (Plan, error) {
	targetVersion, err := normalizeTargetCLIVersion(targetCLIVersion)
	if err != nil {
		return Plan{}, fmt.Errorf("target CLI version: %w", err)
	}
	project, err := LoadProject(projectPath)
	if err != nil {
		return Plan{}, fmt.Errorf("invalid project: %w", err)
	}
	plan := Plan{
		Project:          project,
		TargetCLIVersion: targetVersion,
		TargetPkgVersion: currentTargetPkgVersion,
		Status:           StatusManagedUpToDate,
	}
	if project.SourceFamily == "" {
		if project.SourceCLIVersion != targetVersion {
			plan.Status = StatusUnsupported
			plan.Conflicts = append(plan.Conflicts, Diagnostic{
				Code:    "unsupported_source",
				Path:    "scripts/env",
				Message: fmt.Sprintf("source CLI version %s has no frozen compatibility assets", project.SourceCLIVersion),
			})
			plan.ManualActions, _ = ScanManualActions(project, nil, targetVersion)
			return plan, nil
		}
		// A migrated marker records the target CLI version, not the historical
		// source family. Use the union of assets that this target may have
		// written so a second preview remains idempotent.
		project.SourceFamily = migratedSourceFamily
		plan.Project.SourceFamily = project.SourceFamily
	}
	assetPaths, err := CompatibilityAssetPaths(project.SourceFamily)
	if err != nil {
		return Plan{}, err
	}
	assetSet := make(map[string]struct{}, len(assetPaths))
	for _, assetPath := range assetPaths {
		assetSet[assetPath] = struct{}{}
	}

	versions := make(map[string]struct{})
	err = filepath.WalkDir(project.Root, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relativePath, err := filepath.Rel(project.Root, currentPath)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath)
		_, isAsset := assetSet[relativePath]
		isScriptPatch := relativePath == generateScriptPath
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			if isAsset || relativePath == legacyPublicEmbedPath || isScriptPatch {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unsafe_file", Path: relativePath, Message: "managed path is a directory"})
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if isAsset || relativePath == legacyPublicEmbedPath || isScriptPatch {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unsafe_file", Path: relativePath, Message: "managed path is a symlink"})
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			if isAsset || isScriptPatch {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unsafe_file", Path: relativePath, Message: "managed path is not a regular file"})
			}
			return nil
		}
		if isScriptPatch {
			body, err := readRegularFile(currentPath, maxManagedFileSize)
			if err != nil {
				return fmt.Errorf("read managed file %s: %w", relativePath, err)
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			target, err := renderGenerateScriptTarget(targetVersion)
			if err != nil {
				return fmt.Errorf("render managed file %s: %w", relativePath, err)
			}
			if change, _ := planGenerateScript(body, info.Mode(), target); change != nil {
				plan.Changes = append(plan.Changes, *change)
			}
			return nil
		}

		line, err := readFirstLine(currentPath)
		if err != nil {
			if isAsset || relativePath == legacyPublicEmbedPath {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unreadable_marker", Path: relativePath, Message: err.Error()})
			}
			return nil
		}
		if relativePath == legacyPublicEmbedPath && strings.Contains(line, `grpc-kit-cli/{{ .Global.ReleaseVersion }}`) {
			body, err := readRegularFile(currentPath, maxManagedFileSize)
			if err != nil {
				return err
			}
			if IsKnownLegacyMarkerFile(relativePath, body) {
				plan.Warnings = append(plan.Warnings, Diagnostic{Code: "known_legacy_marker", Path: relativePath, Message: "frozen legacy placeholder output is skipped"})
			} else {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "invalid_marker", Path: relativePath, Message: "legacy placeholder marker does not match the frozen fixture"})
			}
			return nil
		}

		marker, markerErr := ParseMarkerLine(line)
		if markerErr != nil {
			if looksLikeManagedMarker(line) {
				plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "invalid_marker", Path: relativePath, Message: markerErr.Error()})
			}
			return nil
		}
		versions[marker.Version] = struct{}{}
		if !isAsset {
			plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unmapped_managed_file", Path: relativePath, Message: "managed file has no compatibility asset"})
			return nil
		}
		if marker.Version != targetVersion && !marker.Supported() {
			plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "unsupported_version", Path: relativePath, Message: fmt.Sprintf("marker version %s is not supported", marker.Version)})
			return nil
		}

		body, err := readRegularFile(currentPath, maxManagedFileSize)
		if err != nil {
			return fmt.Errorf("read managed file %s: %w", relativePath, err)
		}
		renderProject := project
		if marker.SourceFamily != "" {
			renderProject.SourceFamily = marker.SourceFamily
		}
		after, err := RenderCompatibilityAsset(renderProject, relativePath, targetVersion)
		if err != nil {
			plan.Conflicts = append(plan.Conflicts, Diagnostic{Code: "render_failed", Path: relativePath, Message: err.Error()})
			return nil
		}
		if bytes.Equal(body, after) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		plan.Changes = append(plan.Changes, Change{
			Path:         relativePath,
			Before:       body,
			After:        after,
			Mode:         info.Mode(),
			BeforeSHA256: sha256.Sum256(body),
		})
		return nil
	})
	if err != nil {
		return Plan{}, fmt.Errorf("scan project: %w", err)
	}

	for version := range versions {
		plan.FromCLIVersions = append(plan.FromCLIVersions, version)
	}
	sort.Strings(plan.FromCLIVersions)
	if len(plan.FromCLIVersions) > 1 {
		plan.Warnings = append(plan.Warnings, Diagnostic{Code: "mixed_versions", Message: "managed files contain multiple supported CLI versions: " + strings.Join(plan.FromCLIVersions, ", ")})
	}
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Path < plan.Changes[j].Path })
	sortDiagnostics(plan.Warnings)
	sortDiagnostics(plan.Conflicts)

	if len(plan.Conflicts) != 0 {
		plan.Status = StatusConflict
		for _, conflict := range plan.Conflicts {
			if conflict.Code == "unsupported_source" || conflict.Code == "unsupported_version" {
				plan.Status = StatusUnsupported
				break
			}
		}
	} else if len(plan.Changes) != 0 {
		plan.Status = StatusReady
	}
	plan.ApplyBlockers = inspectApplyBlockers(project.Root, targetVersion)
	changedPaths := make(map[string]struct{}, len(plan.Changes))
	for _, change := range plan.Changes {
		changedPaths[change.Path] = struct{}{}
	}
	manualActions, diagnosticsErr := ScanManualActions(project, changedPaths, targetVersion)
	if diagnosticsErr != nil {
		plan.Warnings = append(plan.Warnings, Diagnostic{Code: "diagnostics_failed", Message: diagnosticsErr.Error()})
	} else {
		plan.ManualActions = manualActions
	}
	sortDiagnostics(plan.ManualActions)
	return plan, nil
}

func looksLikeManagedMarker(line string) bool {
	return strings.HasPrefix(line, `// Code generated by "grpc-kit-cli/`) ||
		strings.HasPrefix(line, `# Code generated by "grpc-kit-cli/`)
}

func readFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	reader := bufio.NewReader(io.LimitReader(file, 4097))
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	if len(line) > 4096 {
		return "", fmt.Errorf("first line exceeds 4096 bytes")
	}
	return strings.TrimSuffix(line, "\n"), nil
}

func sortDiagnostics(diagnostics []Diagnostic) {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].Path == diagnostics[j].Path {
			if diagnostics[i].Line != diagnostics[j].Line {
				return diagnostics[i].Line < diagnostics[j].Line
			}
			return diagnostics[i].Code < diagnostics[j].Code
		}
		return diagnostics[i].Path < diagnostics[j].Path
	})
}
