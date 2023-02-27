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

	"github.com/grpc-kit/pkg/vars"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of grpc-kit-cli",
	Long:  `All software has versions. This is grpc-kit-cli's.`,
	RunE:  runFuncVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runFuncVersion(cmd *cobra.Command, args []string) error {
	rawBody, err := json.MarshalIndent(vars.GetVersion(), "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(rawBody))
	return nil
}
