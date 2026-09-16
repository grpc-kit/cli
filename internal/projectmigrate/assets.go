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
	"embed"
	"fmt"
	"go/parser"
	"go/token"
	"path"
	"strings"
	"text/template"
)

const compatibilityAssetRoot = "assets/v0_5_0/v0_3_8"

var compatibilityAssetPaths = []string{
	"cmd/server/main.go",
	"handler/microservice.go",
	"handler/register.go",
	"handler/rpc_internal.go",
	"handler/shutdown.go",
	"scripts/env",
}

//go:embed assets/v0_5_0/v0_3_8
var compatibilityAssets embed.FS

type assetData struct {
	TargetCLIVersion string
	Repository       string
	ProductCode      string
	ShortName        string
	APIVersion       string
	ServiceTitle     string
}

// CompatibilityAssetPaths returns the existing paths that the selected source
// family can map to pkg v0.5.0-compatible managed content.
func CompatibilityAssetPaths(sourceFamily string) ([]string, error) {
	if !hasCompatibilityAssets(sourceFamily) {
		return nil, fmt.Errorf("source family %q has no frozen compatibility assets", sourceFamily)
	}
	return append([]string(nil), compatibilityAssetPaths...), nil
}

// RenderCompatibilityAsset renders one frozen asset entirely in memory.
func RenderCompatibilityAsset(project Project, relativePath, targetCLIVersion string) ([]byte, error) {
	if !hasCompatibilityAssets(project.SourceFamily) {
		return nil, fmt.Errorf("source family %q has no frozen compatibility assets", project.SourceFamily)
	}
	if !containsAssetPath(relativePath) {
		return nil, fmt.Errorf("path %q has no compatibility asset", relativePath)
	}
	version, err := normalizeTargetCLIVersion(targetCLIVersion)
	if err != nil {
		return nil, fmt.Errorf("target CLI version: %w", err)
	}
	templatePath := path.Join(compatibilityAssetRoot, relativePath+".tmpl")
	body, err := compatibilityAssets.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("read compatibility asset %s: %w", relativePath, err)
	}
	tmpl, err := template.New(relativePath).Option("missingkey=error").Parse(string(body))
	if err != nil {
		return nil, fmt.Errorf("parse compatibility asset %s: %w", relativePath, err)
	}
	data := assetData{
		TargetCLIVersion: version,
		Repository:       project.ModulePath,
		ProductCode:      project.ProductCode,
		ShortName:        project.ShortName,
		APIVersion:       project.APIVersion,
		ServiceTitle:     identifierTitle(project.ProductCode) + identifierTitle(project.ShortName),
	}
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return nil, fmt.Errorf("render compatibility asset %s: %w", relativePath, err)
	}
	result := rendered.Bytes()
	if bytes.Contains(result, []byte("{{")) || bytes.Contains(result, []byte("}}")) {
		return nil, fmt.Errorf("render compatibility asset %s: unresolved template expression", relativePath)
	}
	marker, err := ParseMarkerLine(firstLine(result))
	if err != nil || marker.Version != version {
		return nil, fmt.Errorf("render compatibility asset %s: invalid target marker", relativePath)
	}
	if strings.HasSuffix(relativePath, ".go") {
		if _, err := parser.ParseFile(token.NewFileSet(), relativePath, result, parser.AllErrors); err != nil {
			return nil, fmt.Errorf("render compatibility asset %s: invalid Go syntax: %w", relativePath, err)
		}
	}
	return result, nil
}

func hasCompatibilityAssets(sourceFamily string) bool {
	return sourceFamily == "v0.3.8" || sourceFamily == "v0.3.9-beta.1"
}

func containsAssetPath(relativePath string) bool {
	for _, candidate := range compatibilityAssetPaths {
		if candidate == relativePath {
			return true
		}
	}
	return false
}

func firstLine(body []byte) string {
	line, _, _ := bytes.Cut(body, []byte("\n"))
	return string(line)
}

func identifierTitle(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(value[:1]) + value[1:]
}
