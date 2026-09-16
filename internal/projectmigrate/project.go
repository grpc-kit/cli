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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const (
	maxEnvironmentSize = 64 << 10
	maxGoModSize       = 1 << 20
)

var (
	projectNamePattern = regexp.MustCompile(`^[a-z0-9]{4,}$`)
	apiVersionPattern  = regexp.MustCompile(`^v[1-9][0-9]*(?:(?:alpha|beta)[1-9][0-9]*)?$`)
)

var derivedEnvironmentLines = map[string]string{
	"APPNAME":      `${PRODUCT_CODE}-${SHORT_NAME}-${API_VERSION}`,
	"SERVICE_CODE": `${SHORT_NAME}.${API_VERSION}.${PRODUCT_CODE}`,
}

// Project contains the identity required to render managed migration assets.
type Project struct {
	Root             string
	ModulePath       string
	SourceCLIVersion string
	SourceFamily     string
	ProductCode      string
	ShortName        string
	APIVersion       string
}

type environment struct {
	Marker      Marker
	CLIVersion  string
	ProductCode string
	ShortName   string
	APIVersion  string
}

// LoadProject identifies a project without executing its scripts or changing
// any files.
func LoadProject(path string) (Project, error) {
	root, err := filepath.Abs(path)
	if err != nil {
		return Project{}, fmt.Errorf("resolve project path: %w", err)
	}
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return Project{}, fmt.Errorf("inspect project root: %w", err)
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
		return Project{}, fmt.Errorf("project root must be a real directory, not a symlink")
	}

	environmentBody, err := readRegularFile(filepath.Join(root, "scripts", "env"), maxEnvironmentSize)
	if err != nil {
		return Project{}, fmt.Errorf("read scripts/env: %w", err)
	}
	env, err := parseEnvironment(environmentBody)
	if err != nil {
		return Project{}, fmt.Errorf("parse scripts/env: %w", err)
	}

	goModBody, err := readRegularFile(filepath.Join(root, "go.mod"), maxGoModSize)
	if err != nil {
		return Project{}, fmt.Errorf("read go.mod: %w", err)
	}
	parsedGoMod, err := modfile.Parse("go.mod", goModBody, nil)
	if err != nil {
		return Project{}, fmt.Errorf("parse go.mod: %w", err)
	}
	if parsedGoMod.Module == nil || parsedGoMod.Module.Mod.Path == "" {
		return Project{}, fmt.Errorf("parse go.mod: module path is missing")
	}
	modulePath := parsedGoMod.Module.Mod.Path
	if err := module.CheckPath(modulePath); err != nil {
		return Project{}, fmt.Errorf("parse go.mod: invalid module path %q: %w", modulePath, err)
	}

	return Project{
		Root:             root,
		ModulePath:       modulePath,
		SourceCLIVersion: env.CLIVersion,
		SourceFamily:     env.Marker.SourceFamily,
		ProductCode:      env.ProductCode,
		ShortName:        env.ShortName,
		APIVersion:       env.APIVersion,
	}, nil
}

func parseEnvironment(body []byte) (environment, error) {
	scanner := bufio.NewScanner(bytes.NewReader(body))
	if !scanner.Scan() {
		return environment{}, fmt.Errorf("file is empty")
	}
	marker, err := ParseMarkerLine(scanner.Text())
	if err != nil {
		return environment{}, err
	}

	values := make(map[string]string, 4)
	lineNumber := 1
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found || key == "" {
			return environment{}, fmt.Errorf("line %d is not a literal KEY=value assignment", lineNumber)
		}
		if expected, ok := derivedEnvironmentLines[key]; ok {
			if value != expected {
				return environment{}, fmt.Errorf("line %d has unexpected derived value for %s", lineNumber, key)
			}
			continue
		}
		switch key {
		case "CLI_VERSION", "PRODUCT_CODE", "SHORT_NAME", "API_VERSION":
		default:
			return environment{}, fmt.Errorf("line %d assigns unsupported key %q", lineNumber, key)
		}
		if _, exists := values[key]; exists {
			return environment{}, fmt.Errorf("line %d duplicates %s", lineNumber, key)
		}
		if value == "" || strings.ContainsAny(value, " \t$`;&|(){}") {
			return environment{}, fmt.Errorf("line %d has non-literal value for %s", lineNumber, key)
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return environment{}, fmt.Errorf("scan file: %w", err)
	}

	for _, key := range []string{"CLI_VERSION", "PRODUCT_CODE", "SHORT_NAME", "API_VERSION"} {
		if values[key] == "" {
			return environment{}, fmt.Errorf("required key %s is missing", key)
		}
	}
	cliVersion, err := normalizeVersion(values["CLI_VERSION"])
	if err != nil {
		return environment{}, fmt.Errorf("invalid CLI_VERSION: %w", err)
	}
	if cliVersion != marker.Version {
		return environment{}, fmt.Errorf("marker version %q does not match CLI_VERSION %q", marker.Version, cliVersion)
	}
	if !projectNamePattern.MatchString(values["PRODUCT_CODE"]) {
		return environment{}, fmt.Errorf("PRODUCT_CODE %q does not match %s", values["PRODUCT_CODE"], projectNamePattern)
	}
	if !projectNamePattern.MatchString(values["SHORT_NAME"]) {
		return environment{}, fmt.Errorf("SHORT_NAME %q does not match %s", values["SHORT_NAME"], projectNamePattern)
	}
	if !apiVersionPattern.MatchString(values["API_VERSION"]) {
		return environment{}, fmt.Errorf("API_VERSION %q does not match %s", values["API_VERSION"], apiVersionPattern)
	}

	return environment{
		Marker:      marker,
		CLIVersion:  cliVersion,
		ProductCode: values["PRODUCT_CODE"],
		ShortName:   values["SHORT_NAME"],
		APIVersion:  values["API_VERSION"],
	}, nil
}

func readRegularFile(path string, limit int64) ([]byte, error) {
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 || !pathInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("%s must be a regular file, not a symlink", path)
	}
	if pathInfo.Size() > limit {
		return nil, fmt.Errorf("%s exceeds the %d-byte limit", path, limit)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return nil, fmt.Errorf("%s changed while it was being opened", path)
	}
	body, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > limit {
		return nil, fmt.Errorf("%s exceeds the %d-byte limit", path, limit)
	}
	return body, nil
}
