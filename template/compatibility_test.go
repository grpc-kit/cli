// Copyright © 2026 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package template

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/grpc-kit/cli/config"
	"golang.org/x/mod/modfile"
)

const runTemplateCompatibilityTestEnv = "GRPC_KIT_RUN_TEMPLATE_COMPATIBILITY_TEST"

// TestGeneratedServiceCompilesAgainstReleasedPkg is an opt-in release gate.
// It intentionally uses the public module proxy and rejects local replacements.
func TestGeneratedServiceCompilesAgainstReleasedPkg(t *testing.T) {
	if os.Getenv(runTemplateCompatibilityTestEnv) != "1" {
		t.Skip("set " + runTemplateCompatibilityTestEnv + "=1 to run the generated-service release gate")
	}

	root := filepath.Join(t.TempDir(), "echo")
	generator, err := New(config.Config{
		Global: config.GlobalConfig{
			Type:           TypeService,
			APIEndpoint:    "api",
			ProductCode:    "demo",
			ShortName:      "echo",
			ReleaseVersion: "v0.5.0",
			Repository:     "example.com/grpc-kit/echo",
			Organization:   "grpc-kit",
			Appname:        "demo-echo-v1",
			ServiceCode:    "echo.v1.demo",
			ProtoPackage:   "grpc_kit.api.demo.echo.v1",
			ServiceTitle:   "DemoEcho",
		},
		Template: config.TemplateConfig{Service: config.TemplateService{APIVersion: "v1"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = generator.GenerateTo(root); err != nil {
		t.Fatalf("generate service: %v", err)
	}

	moduleBody, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	module, err := modfile.Parse("go.mod", moduleBody, nil)
	if err != nil {
		t.Fatalf("parse generated go.mod: %v", err)
	}
	if len(module.Replace) != 0 {
		t.Fatalf("generated go.mod contains local replacements: %+v", module.Replace)
	}
	var pkgVersion string
	for _, requirement := range module.Require {
		if requirement.Mod.Path == "github.com/grpc-kit/pkg" {
			pkgVersion = requirement.Mod.Version
		}
	}
	if pkgVersion != "v0.5.0" {
		t.Fatalf("generated pkg version = %q, want v0.5.0", pkgVersion)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	gopathOutput, err := exec.CommandContext(ctx, "go", "env", "GOPATH").Output()
	if err != nil {
		t.Fatalf("go env GOPATH: %v", err)
	}
	environment := append(os.Environ(), "GOPATH="+strings.TrimSpace(string(gopathOutput)), "GOWORK=off")
	for _, command := range [][]string{
		{"make", "generate"},
		{"go", "test", "./...", "-count=1"},
		{"go", "build", "./..."},
	} {
		process := exec.CommandContext(ctx, command[0], command[1:]...)
		process.Dir = root
		process.Env = environment
		if output, runErr := process.CombinedOutput(); runErr != nil {
			t.Fatalf("%s: %v\n%s", strings.Join(command, " "), runErr, output)
		}
	}
}
