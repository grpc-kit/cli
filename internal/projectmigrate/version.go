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
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// ValidateTargetCLIVersion accepts only a stable, release-build semantic
// version and returns the marker representation without a leading v.
func ValidateTargetCLIVersion(raw string) (string, error) {
	version, err := normalizeTargetCLIVersion(raw)
	if err != nil {
		return "", err
	}
	canonical := "v" + version
	if semver.Prerelease(canonical) != "" {
		return "", fmt.Errorf("CLI release version %q is a prerelease", raw)
	}

	return version, nil
}

// normalizeTargetCLIVersion accepts a non-zero canonical version for preview.
// Apply adds the stronger stable-release requirement above.
func normalizeTargetCLIVersion(raw string) (string, error) {
	version, err := normalizeVersion(raw)
	if err != nil {
		return "", err
	}
	if semver.Compare("v"+version, "v0.0.0") <= 0 {
		return "", fmt.Errorf("CLI release version %q is a development default", raw)
	}
	return version, nil
}

func normalizeVersion(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("CLI release version was not injected")
	}
	if trimmed != raw {
		return "", fmt.Errorf("CLI release version %q contains surrounding whitespace", raw)
	}
	canonical := trimmed
	if !strings.HasPrefix(canonical, "v") {
		canonical = "v" + canonical
	}
	if !semver.IsValid(canonical) || semver.Canonical(canonical) != canonical {
		return "", fmt.Errorf("CLI release version %q is not canonical semantic version", raw)
	}
	return strings.TrimPrefix(canonical, "v"), nil
}
