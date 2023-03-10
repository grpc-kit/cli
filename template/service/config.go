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

package service

func (t *templateService) fileDirectoryConfig() {
	t.files = append(t.files, &templateFile{
		name:  "config/app-dev-local.yaml",
		parse: true,
		body: `
# https://github.com/grpc-kit/pkg/blob/main/cfg/app-sample.yaml

# 基础服务配置
services:
  # 服务注册的前缀，全局统一
  root_path: service
  # 服务注册的空间，全局统一
  namespace: example
  # 服务的代码，名称唯一且必填，格式：应用短名.接口版本.产品代码
  service_code: {{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}
  # 接口网关的地址
  api_endpoint: {{ .Global.APIEndpoint }}
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc_address: 0.0.0.0:10081
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http_address: 0.0.0.0:8080

# 服务注册配置
#discover:
#  driver: etcdv3
#  heartbeat: 15
#  endpoints:
#    - http://127.0.0.1:2379

# 认证鉴权配置
security:
  enable: true
  authentication:
    insecure_rpcs:
      - Check
      - HealthCheck
    http_users:
      - username: user1
        password: grpc-kit-cli

# 日志调试配置
debugger:
  enable_pprof: true
  log_level: debug
  log_format: text

# 应用私有配置
independent:
  name: grpc-kit
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "config/app-mini.yaml",
		parse: true,
		body: `
# https://github.com/grpc-kit/pkg/blob/main/cfg/app-mini.yaml

# 基础服务配置
services:
  # 服务的代码，名称唯一且必填，格式：应用短名.接口版本.产品代码
  service_code: {{ .Global.ServiceCode }}
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http_address: 127.0.0.1:8080
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc_address: 127.0.0.1:10081

# 应用私有配置
independent:
  name: grpc-kit
`,
	})
}
