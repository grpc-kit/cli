// Copyright © 2020 The gRPC Kit Authors.
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
	"encoding/json"
	"fmt"

	"github.com/grpc-kit/cli/internal/buildinfo"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	short := false
	cmd := &cobra.Command{
		Use:         "version",
		Short:       "Print the version number of grpc-kit-cli",
		Long:        `All software has versions. This is grpc-kit-cli's.`,
		Args:        cobra.NoArgs,
		Annotations: map[string]string{skipUserConfigAnnotation: "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFuncVersion(cmd, short)
		},
	}
	cmd.Flags().BoolVar(&short, "short", false, "print only the release version")
	return cmd
}

func runFuncVersion(cmd *cobra.Command, short bool) error {
	info := buildinfo.Get()
	if short {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), info.ReleaseVersion)
		return err
	}

	rawBody, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(rawBody))
	return err
}
