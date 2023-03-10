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
	"os"

	"github.com/grpc-kit/cli/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgType config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grpc-kit-cli",
	Short: "grpc-kit-cli is the main command, used to build your product easy.",
	Long: `gRPC Kit Cli is used to quickly generate code templates for the same product
to comply with the same specifications.

Complete documentation is available at https://grpc-kit.com/.
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		// os.Exit(1)
	}
}

func init() {
	// 注册方法，可在command的Execute方法执行前(preRun)调用
	cobra.OnInitialize(initConfig, unmarshalConfig)

	// 全局参数设置，也就是可在其他command中生效
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file, priority higher than flags (default \"$HOME/.grpc-kit-cli.yaml\")")
	rootCmd.PersistentFlags().StringVarP(&cfgType.Global.Type, "type", "t", "service", "code templates type, current support: service")
	rootCmd.PersistentFlags().StringVar(&cfgType.Global.APIEndpoint, "api-endpoint", "api.grpc-kit.com", "api gateway of the company")
	rootCmd.PersistentFlags().StringVar(&cfgType.Global.GitDomain, "git-domain", "github.com", "the git domain name used to save this code")
	rootCmd.PersistentFlags().StringVarP(&cfgType.Global.ProductCode, "product-code", "p", "", "the product code (must match regex \"^([a-z0-9]){4,}$\")")
	rootCmd.PersistentFlags().StringVarP(&cfgType.Global.ShortName, "short-name", "s", "", "application name under the product (must match regex \"^([a-z0-9]){4,}$\")")
	rootCmd.PersistentFlags().StringVarP(&cfgType.Global.Repository, "repository", "r", "git-domain/product-code/short-name", "name to use for go module (e.g., github.com/user/repo)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grpc-kit-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func unmarshalConfig() {
	// 优先级以用户在 cli 上指定的值最高
	productCode := cfgType.Global.ProductCode

	if err := viper.Unmarshal(&cfgType); err != nil {
		// return err
	}

	if productCode != "" {
		cfgType.Global.ProductCode = productCode
	}
}
