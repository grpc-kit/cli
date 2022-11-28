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
  grpc_address: 127.0.0.1:10081
  # 服务所监听的http地址（如未设置，则不开启gateway服务）
  http_address: 127.0.0.1:8080

# 服务注册配置
#discover:
#  driver: etcdv3
#  heartbeat: 15
#  endpoints:
#    - http://127.0.0.1:2379

# 认证鉴权配置
security:
  enable: false

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
		name:  "config/app-sample.yaml",
		parse: true,
		body: `
# https://github.com/grpc-kit/pkg/blob/main/cfg/app-sample.yaml

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
discover:
  driver: etcdv3
  heartbeat: 15
  endpoints:
    - http://127.0.0.1:2379
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
      - SearchHosts
    oidc_provider:
      issuer: https://accounts.example.com
      config:
        # 必须验证token.aud是否与client_id相等
        client_id: example
        # 允许的签名算法类别
        supported_signing_algs:
          - RS256
        # 忽略token.aud与client_id的验证
        skip_client_id_check: true
        # 忽略token是否过期的验证
        skip_expiry_check: false
        # 忽略token issuer的验证
        skip_issuer_check: true
        # 是否跳过issuer的ca验证
        insecure_skip_verify: true
    http_users:
      - username: user1
        password: pass1
  # TODO; 鉴权：能做什么
  authorization:

# 关系数据配置
database:
  enable: true
  driver: mysql
  #driver: postgres
  dbname: demo
  username: demo
  password: password
  address: 192.168.31.200:3306
  parameters: ""
  #address: 192.168.31.200:5432
  #parameters: "sslmode=disable"
  connection_pool:
    max_idle_time: 1800s
    max_life_time: 21600s
    max_idle_conns: 300
    max_open_conns: 300

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
    http_body: true
    http_response: true

# 事件通道配置
cloudevents:
  protocol: "kafka_sarama"
  kafka_sarama:
    topic: "uptime-test"
    brokers:
      - 127.0.0.1:19092
      - 127.0.0.1:29092
      - 127.0.0.1:39092
    config:
      net:
        max_open_requests: 5
        dial_timeout: 30s
        read_timeout: 30s
        write_timeout: 30s
        tls:
          enable: false
        sasl:
          enable: true
          mechanism: "SCRAM-SHA-256"
          user: "uptime"
          password: "testkey"
        keep_alive: 40s
      metadata:
        retry:
          max: 3
          backoff: 250ms
        refresh_frequency: 10m0s
        full: true
        allow_auto_topic_creation: false
      producer:
        max_message_bytes: 1000000
        required_acks: 1
        timeout: 10s
        return:
          successes: false
          errors: true
        flush:
          bytes: 104857600
          frequency: 30s
          max_messages: 999
        retry:
          max: 3
          backoff: 100ms
      consumer:
        group:
          session:
            timeout: 10s
          heartbeat:
            interval: 3s
          rebalance:
            strategy: range
            timeout: 55s
            retry:
              max: 4
              backoff: 2s
        retry:
          backoff: 2s
        fetch:
          min: 1
        max_wait_time: 250ms
        max_processing_time: 100ms
        return:
          errors: true
        offsets:
          auto_commit:
            enable: true
            interval: 1s
          retry:
            max: 3
      version: "2.4.0"

# 应用私有配置
independent:
  name: grpc-kit
`,
	})
}
