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
	"fmt"
	"regexp"

	"github.com/grpc-kit/cli/template"
	"github.com/grpc-kit/pkg/vars"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new code templates for your product",
	Long: `Create a new code templates for your product. It will only be used when 
it is newly created. For example:

./grpc-kit-cli new -t service -p opsaid -s test1
`,
	RunE:          runFuncNew,
	SilenceUsage:  true,
	SilenceErrors: false,
}

func init() {
	rootCmd.AddCommand(newCmd)

	// 只在该command下生效的参数
	newCmd.Flags().StringVar(&cfgType.Template.Service.APIVersion,
		"api-version", "v1", "api version, like: v1alpha1, v1beta1, v1")
}

func runFuncNew(cmd *cobra.Command, args []string) error {
	re := regexp.MustCompile(`^([a-z0-9]){4,}$`)

	// 必须存在的参数校验
	if cfgType.Global.ProductCode == "" {
		return fmt.Errorf("must set -p or --product-code")
	}
	if cfgType.Global.ShortName == "" {
		return fmt.Errorf("must set -s or --short-name")
	}
	if cfgType.Global.ReleaseVersion == "" {
		cfgType.Global.ReleaseVersion = vars.ReleaseVersion
		if cfgType.Global.ReleaseVersion == "" {
			cfgType.Global.ReleaseVersion = "v0.0.0"
		}
	}
	if !re.MatchString(cfgType.Global.ProductCode) {
		return fmt.Errorf("product-code: %v, not match regex", cfgType.Global.ProductCode)
	}
	if !re.MatchString(cfgType.Global.ShortName) {
		return fmt.Errorf("short-name: %v, not match regex", cfgType.Global.ShortName)
	}
	if cfgType.Global.Repository == "" || cfgType.Global.Repository == "git-domain/product-code/short-name" {
		cfgType.Global.Repository = fmt.Sprintf("%v/%v/%v",
			cfgType.Global.GitDomain, cfgType.Global.ProductCode, cfgType.Global.ShortName)
	}

	fmt.Println(
		fmt.Sprintf("Generate code templates type: %v, use git repos: %v",
			cfgType.Global.Type, cfgType.Global.Repository))

	t, err := template.New(cfgType)
	if err != nil {
		return err
	}

	if err := t.Generate(); err != nil {
		return err
	}

	return nil
}
