// Copyright © 2020 Li MingQing <mingqing@henji.org>
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
		name:  "config/app-dev-local.toml",
		parse: true,
		body: `
# https://github.com/grpc-kit/cfg/blob/master/app-sample.toml

# 基础服务配置
[services]
  # 服务注册的前缀，全局统一
  root-path = "service"
  # 服务注册的空间，全局统一
  namespace = "example"
  # 服务的代码，设置后不可变更
  service-code = "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}"
  # 接口网关的地址
  api-endpoint = "{{ .Global.APIEndpoint }}"
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc-address = "127.0.0.1:10081"
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http-address = "127.0.0.1:10080"
  # 服务注册，外部网络可连接的grpc地址（一般等同于grpc-address）
  public-address = "127.0.0.1:10081"

# 服务注册配置
#[discover]
  #driver = "etcdv3"
  #heartbeat = 15
  #endpoints = [ "http://127.0.0.1:2379" ]

# 认证鉴权配置
[security]
  enable = false

# 关系数据配置
[database]
  driver = "postgres"
  dbname = "demo"
  user = "demo"
  password = "grpc-kit"
  host = "127.0.0.1"
  port = 5432
  sslmode = "disable"
  connect-timeout = 10

# 缓存服务配置
[cachebuf]
  enable = true
  driver = "redis"
  address = "127.0.0.1:6379"
  password = ""

# 日志调试配置
[debugger]
  # 日志输出级别，可取：panic、fatal、error、warn、info、debug、trace
  log-level = "debug"
  # 日志输出的格式，可取：json、text
  log-format = "text"

# 链路追踪配置
[opentracing]
  enable = true
  host = "127.0.0.1"
  port = 16831

# 应用私有配置
[independent]
  name = "grpc-kit"
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "config/app-sample.toml",
		parse: true,
		body: `
# https://github.com/grpc-kit/cfg/blob/master/app-sample.toml

# 基础服务配置
[services]
  # 服务注册的前缀，全局统一
  root-path = "service"
  # 服务注册的空间，全局统一
  namespace = "example"
  # 服务的代码，设置后不可变更
  service-code = "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}"
  # 接口网关的地址
  api-endpoint = "{{ .Global.APIEndpoint }}"
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc-address = "127.0.0.1:10081"
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http-address = "127.0.0.1:10080"
  # 服务注册，外部网络可连接的grpc地址（一般等同于grpc-address）
  public-address = "127.0.0.1:10081"

# 服务注册配置
[discover]
  driver = "etcdv3"
  heartbeat = 15
  endpoints = [ "http://127.0.0.1:2379" ]
  #endpoints = [ "https://node1:2379", "https://node2:2379", "https://node3:2379" ]
  #[discover.tls]
  #  ca-file = "/opt/certs/etcd-ca.pem"
  #  cert-file = "/opt/certs/etcd.pem"
  #  key-file = "/opt/certs/etcd-key.pem"

# 认证鉴权配置
[security]
  enable = false

# 关系数据配置
[database]
  driver = "postgres"
  dbname = "demo"
  user = "demo"
  password = "grpc-kit"
  host = "127.0.0.1"
  port = 5432
  sslmode = "disable"
  connect-timeout = 10

# 缓存服务配置
[cachebuf]
  enable = true
  driver = "redis"
  address = "127.0.0.1:6379"
  password = ""

# 日志调试配置
[debugger]
  # 日志输出级别，可取：panic、fatal、error、warn、info、debug、trace
  log-level = "debug"
  # 日志输出的格式，可取：json、text
  log-format = "text"

# 链路追踪配置
[opentracing]
  enable = true
  host = "127.0.0.1"
  port = 16831

# 应用私有配置
[independent]
  name = "grpc-kit"
`,
	})
}
