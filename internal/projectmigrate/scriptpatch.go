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
	"crypto/sha256"
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"

	_ "embed"
)

const generateScriptPath = "scripts/generate.sh"

const (
	// generateScriptAssetSHA256 pins the embedded frozen asset, which carries
	// the {{ .TargetCLIVersion }} marker expression. Whenever
	// template/service/scripts/generate.sh.tmpl changes, re-freeze
	// assets/v0_5_0/scripts/generate.sh.tmpl and update this digest.
	generateScriptAssetSHA256 = "e911a752d48765d560902873785ae04c72aa6fcfcb0bd0e31372aeae449e221e"
	// legacyGenerateScriptV038SHA256 pins the frozen v0.3.8 template blob.
	// The v0.3.8 scaffold rendered this .tmpl through text/template with no
	// template expressions and no leading-trim effect, so the blob is
	// byte-identical to the scripts/generate.sh found in v0.3.8 projects.
	legacyGenerateScriptV038SHA256 = "8ca0626c1f4a5e35ac1004557f662410c63a1fa5402a5916044c59dbff101c48"
)

// scripts/generate.sh keeps its shebang on line 1, so since CLI 0.5.1 the
// service template puts the ownership marker on line 2. Pre-0.5.1 scaffolds
// have no marker at all; those generations stay under the one narrow
// content-evidence exception to the first-line marker rule (see the
// maintain-project-migration SOP, "修复生成器缺陷"): a whole-file rewrite is
// authorized only when the on-disk content matches a frozen digest exactly.
// Content similarity never grants ownership. Marker-carrying copies produced
// by 0.5.1 scaffolds and migrations should migrate through the marker rule in
// future releases (recognizing the shebang+line-2 shape for this path) instead
// of growing legacyGenerateScriptSHA256.

//go:embed assets/v0_5_0/scripts/generate.sh.tmpl
var generateScriptAsset []byte

// legacyGenerateScriptSHA256 maps exact frozen legacy digests to the CLI
// source family that produced them. It only absorbs pre-0.5.1 generations
// without a marker. Source families without a locally available release
// artifact (v0.3.9-beta.1) stay unmapped and surface as manual actions
// instead of rewrites.
var legacyGenerateScriptSHA256 = map[string]string{
	legacyGenerateScriptV038SHA256: "v0.3.8",
}

// renderGenerateScriptTarget renders the frozen asset with the target CLI
// version and validates the shebang plus the line-2 marker, mirroring the
// guarantees RenderCompatibilityAsset enforces for marker assets.
func renderGenerateScriptTarget(targetCLIVersion string) ([]byte, error) {
	version, err := normalizeTargetCLIVersion(targetCLIVersion)
	if err != nil {
		return nil, fmt.Errorf("target CLI version: %w", err)
	}
	parsed, err := template.New(generateScriptPath).Option("missingkey=error").Parse(string(generateScriptAsset))
	if err != nil {
		return nil, fmt.Errorf("parse generate script asset: %w", err)
	}
	var rendered bytes.Buffer
	if err := parsed.Execute(&rendered, assetData{TargetCLIVersion: version}); err != nil {
		return nil, fmt.Errorf("render generate script target: %w", err)
	}
	result := rendered.Bytes()
	if bytes.Contains(result, []byte("{{")) || bytes.Contains(result, []byte("}}")) {
		return nil, fmt.Errorf("render generate script target: unresolved template expression")
	}
	if !bytes.HasPrefix(result, []byte("#!/bin/bash\n")) {
		return nil, fmt.Errorf("render generate script target: missing shebang")
	}
	marker, err := ParseMarkerLine(secondLine(result))
	if err != nil || marker.Version != version {
		return nil, fmt.Errorf("render generate script target: invalid target marker")
	}
	return result, nil
}

// generateScriptEvidence classifies on-disk scripts/generate.sh content by
// exact SHA-256 (frozen legacy generations) or byte equality with the
// rendered target (marker-carrying current content).
type generateScriptEvidence int

const (
	generateScriptUnknown generateScriptEvidence = iota
	generateScriptCurrent
	generateScriptLegacy
)

func classifyGenerateScript(body, target []byte) generateScriptEvidence {
	digest := fmt.Sprintf("%x", sha256.Sum256(body))
	if _, ok := legacyGenerateScriptSHA256[digest]; ok {
		return generateScriptLegacy
	}
	if bytes.Equal(body, target) {
		return generateScriptCurrent
	}
	return generateScriptUnknown
}

// planGenerateScript authorizes a whole-file rewrite of scripts/generate.sh
// only for exact frozen legacy content; anything else is left to manual
// action so locally modified scripts are never overwritten.
func planGenerateScript(body []byte, mode fs.FileMode, target []byte) (change *Change, upToDate bool) {
	switch classifyGenerateScript(body, target) {
	case generateScriptCurrent:
		return nil, true
	case generateScriptLegacy:
		return &Change{
			Path:         generateScriptPath,
			Before:       body,
			After:        bytes.Clone(target),
			Mode:         mode,
			BeforeSHA256: sha256.Sum256(body),
		}, false
	default:
		return nil, false
	}
}

// scanGenerateScriptAction reports scripts/generate.sh when its content is
// neither the rendered target nor a frozen legacy version and the managed
// plan is not already rewriting the whole file. Absent files need nothing,
// and unsafe shapes are reported as unsafe_file conflicts by BuildPlan, so
// every read failure stays silent here.
func scanGenerateScriptAction(project Project, changedPaths map[string]struct{}, actions *[]Diagnostic, dedup map[string]struct{}, targetCLIVersion string) error {
	if _, excluded := changedPaths[generateScriptPath]; excluded {
		return nil
	}
	body, err := readRegularFile(filepath.Join(project.Root, filepath.FromSlash(generateScriptPath)), maxManagedFileSize)
	if err != nil {
		return nil
	}
	target, err := renderGenerateScriptTarget(targetCLIVersion)
	if err != nil {
		return fmt.Errorf("scan %s: %w", generateScriptPath, err)
	}
	if classifyGenerateScript(body, target) != generateScriptUnknown {
		return nil
	}
	appendDiagnostic(actions, dedup, Diagnostic{
		Code:    "generate_script_modified",
		Path:    generateScriptPath,
		Message: "content is neither the current template nor a frozen legacy version; reconcile the protoc include paths and swagger output handling against the latest service template manually",
	})
	return nil
}

func secondLine(body []byte) string {
	_, rest, _ := bytes.Cut(body, []byte("\n"))
	line, _, _ := bytes.Cut(rest, []byte("\n"))
	return string(line)
}
