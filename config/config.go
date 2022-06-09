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

package config

// Config 用于表示配置文件的根结构
type Config struct {
	Global   GlobalConfig   `mapstructure:"global"`
	Template TemplateConfig `mapstructure:"template"`
}

// GlobalConfig 用于表示全局通用配置结构
type GlobalConfig struct {
	// 创建模版的类型
	Type string `mapstructure:"-"`
	// 接口网关地址
	APIEndpoint string `mapstructure:"api_endpoint"`
	// GIT代码仓库的域名
	GitDomain string `mapstructure:"git_domain"`
	// 服务所在的产品代号
	ProductCode string `mapstructure:"product_code"`
	// 服务的名称
	ShortName string
	// cli的版本
	ReleaseVersion string `mapstructure:"release_version"`
	// 代码仓库地址
	Repository string `mapstructure:"repository"`
}

// TemplateConfig 用于表示各代码类型的模版配置结构
type TemplateConfig struct {
	Service TemplateService `mapstructure:"service"`
}

// TemplateService 用于表示微服务类型模版的配置结构
type TemplateService struct {
	// 服务的接口版本
	APIVersion string `mapstructure:"api-version"`
}
