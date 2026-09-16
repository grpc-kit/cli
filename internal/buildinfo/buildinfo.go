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

// Package buildinfo exposes metadata injected into the CLI binary at build time.
package buildinfo

import (
	"fmt"
	"runtime"
	"strconv"
)

var (
	AppName        string
	BuildDate      string
	GitCommit      string
	GitBranch      string
	CLIVersion     string
	CommitUnixTime string
	ReleaseVersion string
)

// Info is the stable JSON representation printed by grpc-kit-cli version.
type Info struct {
	AppName        string `json:"appname"`
	BuildDate      string `json:"build_date"`
	GitCommit      string `json:"git_commit"`
	GitBranch      string `json:"git_branch"`
	GoVersion      string `json:"go_version"`
	Compiler       string `json:"compiler"`
	Platform       string `json:"platform"`
	CLIVersion     string `json:"cli_version"`
	CommitUnixTime int64  `json:"commit_unix_time"`
	ReleaseVersion string `json:"release_version"`
}

// Get returns build metadata with deterministic development-build defaults.
func Get() Info {
	commitUnixTime, err := strconv.ParseInt(CommitUnixTime, 10, 64)
	if err != nil {
		commitUnixTime = 0
	}

	info := Info{
		AppName:        AppName,
		BuildDate:      BuildDate,
		GitCommit:      GitCommit,
		GitBranch:      GitBranch,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		CLIVersion:     CLIVersion,
		CommitUnixTime: commitUnixTime,
		ReleaseVersion: ReleaseVersion,
	}
	if info.GitCommit == "" {
		info.GitCommit = "1234567890123456789012345678901234567890"
	}
	if info.BuildDate == "" {
		info.BuildDate = "1970-01-01T00:00:00Z"
	}
	if info.CLIVersion == "" {
		info.CLIVersion = "0.0.0"
	}
	if info.ReleaseVersion == "" {
		info.ReleaseVersion = "0.0.0"
	}
	return info
}
