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
		name:  "config/app-dev-local.yaml",
		parse: true,
		body: `
# https://github.com/grpc-kit/pkg/blob/master/cfg/app-sample.yaml

# 基础服务配置
services:
  # 服务注册的前缀，全局统一
  root_path: service
  # 服务注册的空间，全局统一
  namespace: example
  # 服务的代码，设置后不可变更
  service_code: {{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}
  # 接口网关的地址
  api_endpoint: {{ .Global.APIEndpoint }}
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc_address: 127.0.0.1:10081
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http_address: 127.0.0.1:10080
  # 服务注册，外部网络可连接的grpc地址（一般等同于grpc-address）
  public_address: ""

# 服务注册配置
#discover:
  #driver: etcdv3
  #heartbeat: 15
  #endpoints:
  #  - http://127.0.0.1:2379
  #discover:
  #  tls:
  #    ca_file: /opt/certs/etcd-ca.pem
  #    cert_file: /opt/certs/etcd.pem
  #    key_file: /opt/certs/etcd-key.pem

# 认证鉴权配置
security:
  enable: true
  # 认证：谁在登录
  authentication:
    # 跳过认证的rpc方法
    insecure_rpcs:
      - HealthCheck
    #oidc_provider:
      #issuer: https://accounts.example.com
      #config:
        # 必须验证token.aud是否与client_id相等
        #client_id: example
        # 允许的签名算法类别
        #supported_signing_algs:
        #  - RS256
        # 忽略token.aud与client_id的验证
        #skip_client_id_check: true
        # 忽略token是否过期的验证
        #skip_expiry_check: false
        # 忽略token issuer的验证
        #skip_issuer_check: true
        # 是否跳过issuer的ca验证
        #insecure_skip_verify: true
    http_users:
      - username: user1
        password: pass1
  # TODO; 鉴权：能做什么
  #authorization:

# 关系数据配置
database:
  driver: postgres
  dbname: demo
  user: demo
  password: grpc-kit
  host: 127.0.0.1
  port: 5432
  sslmode: disable
  connect_timeout: 10

# 缓存服务配置
cachebuf:
  enable: true
  driver: redis
  address: 127.0.0.1:6379
  password: ""

# 日志调试配置
debugger:
  enable_pprof: true
  log_level: debug
  log_format: text

# 链路追踪配置
opentracing:
  enable: true
  host: 127.0.0.1
  port: 6831
  log_fields:
    http_body: false
    http_response: false

# 应用私有配置
independent:
  name: grpc-kit
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "config/app-sample.yaml",
		parse: true,
		body: `
# https://github.com/grpc-kit/pkg/blob/master/cfg/app-sample.yaml

# 基础服务配置
services:
  # 服务注册的前缀，全局统一
  root_path: service
  # 服务注册的空间，全局统一
  namespace: example
  # 服务的代码，设置后不可变更
  service_code: {{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}
  # 接口网关的地址
  api_endpoint: {{ .Global.APIEndpoint }}
  # 服务所监听的grpc地址（如未设置，自动监听在127.0.0.1的随机端口）
  grpc_address: 127.0.0.1:10081
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http_address: 127.0.0.1:10080
  # 服务注册，外部网络可连接的grpc地址（一般等同于grpc-address）
  public_address: ""

# 服务注册配置
#discover:
  #driver: etcdv3
  #heartbeat: 15
  #endpoints:
  #  - http://127.0.0.1:2379
  #discover:
  #  tls:
  #    ca_file: /opt/certs/etcd-ca.pem
  #    cert_file: /opt/certs/etcd.pem
  #    key_file: /opt/certs/etcd-key.pem

# 认证鉴权配置
security:
  enable: true
  # 认证：谁在登录
  authentication:
    # 跳过认证的rpc方法
    #insecure_rpcs:
    #  - SearchHosts
    #oidc_provider:
      #issuer: https://accounts.example.com
      #config:
        # 必须验证token.aud是否与client_id相等
        #client_id: example
        # 允许的签名算法类别
        #supported_signing_algs:
        #  - RS256
        # 忽略token.aud与client_id的验证
        #skip_client_id_check: true
        # 忽略token是否过期的验证
        #skip_expiry_check: false
        # 忽略token issuer的验证
        #skip_issuer_check: true
        # 是否跳过issuer的ca验证
        #insecure_skip_verify: true
    http_users:
      - username: user1
        password: pass1
  # TODO; 鉴权：能做什么
  #authorization:

# 关系数据配置
database:
  driver: postgres
  dbname: demo
  user: demo
  password: grpc-kit
  host: 127.0.0.1
  port: 5432
  sslmode: disable
  connect_timeout: 10

# 缓存服务配置
cachebuf:
  enable: true
  driver: redis
  address: 127.0.0.1:6379
  password: ""

# 日志调试配置
debugger:
  enable_pprof: true
  log_level: debug
  log_format: text

# 链路追踪配置
opentracing:
  enable: true
  host: 127.0.0.1
  port: 6831
  log_fields:
    http_body: false
    http_response: false

# 应用私有配置
independent:
  name: grpc-kit
`,
	})
}
