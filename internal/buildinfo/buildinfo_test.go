// Copyright © 2026 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package buildinfo

import "testing"

func TestGetDefaults(t *testing.T) {
	previousCommit := GitCommit
	previousDate := BuildDate
	previousCLI := CLIVersion
	previousRelease := ReleaseVersion
	previousUnixTime := CommitUnixTime
	GitCommit, BuildDate, CLIVersion, ReleaseVersion, CommitUnixTime = "", "", "", "", "invalid"
	t.Cleanup(func() {
		GitCommit = previousCommit
		BuildDate = previousDate
		CLIVersion = previousCLI
		ReleaseVersion = previousRelease
		CommitUnixTime = previousUnixTime
	})

	info := Get()
	if info.ReleaseVersion != "0.0.0" || info.CLIVersion != "0.0.0" {
		t.Fatalf("Get() versions = %q/%q, want development defaults", info.ReleaseVersion, info.CLIVersion)
	}
	if info.CommitUnixTime != 0 {
		t.Fatalf("Get().CommitUnixTime = %d, want 0", info.CommitUnixTime)
	}
}
