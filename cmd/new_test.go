// Copyright © 2026 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grpc-kit/cli/config"
	"github.com/grpc-kit/cli/template"
)

func TestNewCommandGeneratesToRequestedOutput(t *testing.T) {
	previous := cfgType
	cfgType = config.Config{
		Global: config.GlobalConfig{
			Type:        template.TypeService,
			ProductCode: "demo",
			ShortName:   "echo",
			GitDomain:   "example.com",
		},
	}
	t.Cleanup(func() { cfgType = previous })

	outputPath := filepath.Join(t.TempDir(), "custom-output")
	cmd := newNewCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--output", outputPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputPath, "go.mod")); err != nil {
		t.Fatalf("generated go.mod: %v", err)
	}
	if !strings.Contains(output.String(), "Generate code templates type: service") {
		t.Fatalf("output = %q", output.String())
	}
}

func TestNewCommandRejectsExtraArguments(t *testing.T) {
	cmd := newNewCommand()
	cmd.SetArgs([]string{"unexpected"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want argument validation error")
	}
}
