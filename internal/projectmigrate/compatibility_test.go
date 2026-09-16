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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompatibilityAssetsCompileAgainstReleasedPkg is intentionally opt-in:
// it invokes the Go toolchain and may populate module/build caches. It
// validates the frozen upgrade assets against a released module, without a
// local replace directive. GRPC_KIT_COMPAT_PKG_VERSION defaults to v0.5.0 and
// can select another released version for a compatibility matrix:
//
//	GRPC_KIT_RUN_COMPATIBILITY_TEST=1 GRPC_KIT_COMPAT_PKG_VERSION=v0.5.8 \
//	  go test ./internal/projectmigrate \
//	  -run TestCompatibilityAssetsCompileAgainstReleasedPkg -count=1
func TestCompatibilityAssetsCompileAgainstReleasedPkg(t *testing.T) {
	if os.Getenv("GRPC_KIT_RUN_COMPATIBILITY_TEST") != "1" {
		t.Skip("set GRPC_KIT_RUN_COMPATIBILITY_TEST=1 to compile the fixture against a released pkg version")
	}
	pkgVersion := os.Getenv("GRPC_KIT_COMPAT_PKG_VERSION")
	if pkgVersion == "" {
		pkgVersion = currentTargetPkgVersion
	}
	normalizedPkgVersion, err := normalizeVersion(pkgVersion)
	if err != nil || strings.Contains(normalizedPkgVersion, "-") {
		t.Fatalf("GRPC_KIT_COMPAT_PKG_VERSION must be a stable canonical semantic version: %q", pkgVersion)
	}
	pkgVersion = "v" + normalizedPkgVersion

	root := t.TempDir()
	project := Project{
		Root:         root,
		ModulePath:   "example.com/oneops/spider",
		ProductCode:  "oneops",
		ShortName:    "spider",
		APIVersion:   "v1",
		SourceFamily: "v0.3.8",
	}
	for _, relativePath := range compatibilityAssetPaths {
		body, err := RenderCompatibilityAsset(project, relativePath, "0.4.0")
		if err != nil {
			t.Fatalf("render %s: %v", relativePath, err)
		}
		writeCompatibilityFixtureFile(t, root, relativePath, body)
	}

	fixtureFiles := map[string]string{
		"go.mod": `module example.com/oneops/spider

go ` + currentTargetGoVersion + `

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/grpc-kit/pkg ` + pkgVersion + `
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	google.golang.org/grpc v1.83.1
)
`,
		"modeler/independent_cfg.go": `package modeler

import (
	"context"
	"log/slog"
)

type ClientIndependentOption func(*IndependentCfg)

type IndependentCfg struct {
	logger *slog.Logger
}

func (c *IndependentCfg) Init(ctx context.Context, opts ...ClientIndependentOption) error {
	_ = ctx
	for _, opt := range opts {
		opt(c)
	}
	return nil
}
`,
		"handler/private.go": `package handler

import (
	"net/http"

	"example.com/oneops/spider/modeler"
)

func (m *Microservice) privateExtended() ([]modeler.ClientIndependentOption, error) {
	return nil, nil
}

func (m *Microservice) privateHTTPHandle(*http.ServeMux) error {
	return nil
}
`,
		"api/oneops/spider/v1/service.go": `package v1

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterOneopsSpiderServer(grpc.ServiceRegistrar, any) {}

func RegisterOneopsSpiderHandlerFromEndpoint(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error {
	return nil
}
`,
		"internal/security/embed.go": `package security

import "embed"

//go:embed auth.rego
var Assets embed.FS
`,
		"internal/security/auth.rego": "package oneops.spider.v1\n",
		"public/embed.go": `package public

import "embed"

//go:embed gateway.yaml
var Assets embed.FS
`,
		"public/gateway.yaml": "swagger: '2.0'\n",
	}
	for relativePath, body := range fixtureFiles {
		writeCompatibilityFixtureFile(t, root, relativePath, []byte(body))
	}

	runFixtureGoCommand(t, root, "mod", "tidy")
	runFixtureGoCommand(t, root, "test", "./handler")
	runFixtureGoCommand(t, root, "build", "./...")
}

func writeCompatibilityFixtureFile(t *testing.T, root, relativePath string, body []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func runFixtureGoCommand(t *testing.T, root string, arguments ...string) {
	t.Helper()
	command := exec.Command("go", arguments...)
	command.Dir = root
	command.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=auto")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("go %s: %v\n%s", strings.Join(arguments, " "), err, output)
	}
}
