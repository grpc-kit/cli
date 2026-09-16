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

package cmd

import (
	"fmt"

	"github.com/grpc-kit/cli/internal/buildinfo"
	"github.com/grpc-kit/cli/internal/projectmigrate"
	"github.com/spf13/cobra"
)

type projectMigrateOptions struct {
	apply bool
}

func newProjectMigrateCommand() *cobra.Command {
	options := projectMigrateOptions{}
	cmd := &cobra.Command{
		Use:          "migrate [path]",
		Short:        "Preview or apply managed-file migrations for an existing project",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := "."
			if len(args) == 1 {
				projectPath = args[0]
			}

			if options.apply {
				if _, err := projectmigrate.ValidateTargetCLIVersion(buildinfo.ReleaseVersion); err != nil {
					return fmt.Errorf("refuse --apply: %w", err)
				}
			}

			return runFuncProjectMigrate(cmd, projectPath, options)
		},
	}
	cmd.Flags().BoolVar(&options.apply, "apply", false, "apply the managed-file migration plan")
	return cmd
}

func runFuncProjectMigrate(cmd *cobra.Command, projectPath string, options projectMigrateOptions) error {
	plan, err := projectmigrate.BuildPlan(projectPath, buildinfo.ReleaseVersion)
	if err != nil {
		return err
	}
	if err := projectmigrate.WritePlan(cmd.OutOrStdout(), plan); err != nil {
		return fmt.Errorf("write migration preview: %w", err)
	}
	if plan.Blocked() {
		return fmt.Errorf("migration plan is %s", plan.Status)
	}
	if !options.apply {
		return nil
	}
	if err := projectmigrate.Apply(cmd.Context(), &plan); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(cmd.OutOrStdout(), "Managed files migrated: %d\n", len(plan.Changes)); err != nil {
		return err
	}
	if len(plan.ManualActions) != 0 {
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Managed files migrated; project compatibility is not yet verified.")
	}
	return err
}
