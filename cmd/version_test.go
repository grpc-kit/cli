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
	"encoding/json"
	"testing"

	"github.com/grpc-kit/cli/internal/buildinfo"
)

func TestVersionShortUsesCommandOutput(t *testing.T) {
	previous := buildinfo.ReleaseVersion
	buildinfo.ReleaseVersion = "v0.5.0"
	t.Cleanup(func() { buildinfo.ReleaseVersion = previous })

	cmd := newVersionCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--short"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got, want := output.String(), "v0.5.0\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestVersionJSONUsesStableFields(t *testing.T) {
	cmd := newVersionCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var info buildinfo.Info
	if err := json.Unmarshal(output.Bytes(), &info); err != nil {
		t.Fatalf("Unmarshal output: %v\n%s", err, output.String())
	}
	if info.ReleaseVersion == "" || info.GoVersion == "" || info.Platform == "" {
		t.Fatalf("incomplete version output: %+v", info)
	}
}
