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
	"io"

	"github.com/pmezard/go-difflib/difflib"
)

func WritePlan(output io.Writer, plan Plan) error {
	if _, err := fmt.Fprintf(output, "Status: %s\nTarget CLI: %s\nTarget pkg: %s\n", plan.Status, plan.TargetCLIVersion, plan.TargetPkgVersion); err != nil {
		return err
	}
	if err := writeDiagnostics(output, "Warning", plan.Warnings); err != nil {
		return err
	}
	if err := writeDiagnostics(output, "Apply blocked", plan.ApplyBlockers); err != nil {
		return err
	}
	if err := writeDiagnostics(output, "Manual action", plan.ManualActions); err != nil {
		return err
	}
	if len(plan.ManualActions) != 0 {
		if _, err := fmt.Fprintln(output, "Project compatibility: not verified; complete the manual actions before building."); err != nil {
			return err
		}
	}
	if err := writeDiagnostics(output, "Conflict", plan.Conflicts); err != nil {
		return err
	}
	if len(plan.Changes) == 0 {
		_, err := fmt.Fprintln(output, "Managed changes: none")
		return err
	}
	if _, err := fmt.Fprintf(output, "Managed changes: %d\n", len(plan.Changes)); err != nil {
		return err
	}
	for _, change := range plan.Changes {
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(change.Before)),
			B:        difflib.SplitLines(string(change.After)),
			FromFile: "a/" + change.Path,
			ToFile:   "b/" + change.Path,
			Context:  3,
		})
		if err != nil {
			return err
		}
		if _, err := io.WriteString(output, diff); err != nil {
			return err
		}
	}
	return nil
}

func writeDiagnostics(output io.Writer, label string, diagnostics []Diagnostic) error {
	for _, diagnostic := range diagnostics {
		location := diagnostic.Path
		if location == "" {
			location = "project"
		}
		if diagnostic.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, diagnostic.Line)
		}
		if _, err := fmt.Fprintf(output, "%s [%s] %s: %s\n", label, diagnostic.Code, location, diagnostic.Message); err != nil {
			return err
		}
	}
	return nil
}
