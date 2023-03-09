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

import (
	"fmt"
)

func (t *templateService) fileDirectoryApi() {
	t.files = append(t.files, &templateFile{
		name: fmt.Sprintf("api/%v/%v/%v/microservice.proto",
			t.config.Global.ProductCode,
			t.config.Global.ShortName,
			t.config.Template.Service.APIVersion),
		parse: true,
		body: `
syntax = "proto3";

package {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }};

option go_package = "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }};{{ .Global.ShortName }}{{ .Template.Service.APIVersion }}";

// 引入依赖的外部proto文件
import "github.com/grpc-kit/api/known/status/v1/response.proto";

// 同组RPC方法对应一个proto文件，以该组RPC名称的小写字母为文件名
import "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }}/demo.proto";

// 该微服务支持的 RPC 方法定义
service {{ title .Global.ProductCode }}{{ title .Global.ShortName }} {
  rpc HealthCheck(grpc_kit.api.known.status.v1.HealthCheckRequest) returns (grpc_kit.api.known.status.v1.HealthCheckResponse) {}
  rpc Demo(DemoRequest) returns (DemoResponse) {}
}
`,
	})

	t.files = append(t.files, &templateFile{
		name: fmt.Sprintf("api/%v/%v/%v/demo.proto",
			t.config.Global.ProductCode,
			t.config.Global.ShortName,
			t.config.Template.Service.APIVersion),
		parse: true,
		body: `
syntax = "proto3";

// 根据具体的微服务名称做更改
package {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }};

option go_package = "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }};{{ .Global.ShortName }}{{ .Template.Service.APIVersion }}";

// 引入google公共类型
import "google/protobuf/empty.proto";

// 引入第三方依赖的proto文件
import "github.com/grpc-kit/api/known/example/v1/example.proto";

// DemoRequest Demo方法请求可使用的接口参数
message DemoRequest {
  // UUID 资源编号
  string uuid = 1;

  // Ping 资源内容
  grpc_kit.api.known.example.v1.ExampleRequest ping = 2;
}

// DemoResponse Demo方法响应的具体内容
message DemoResponse {

  message Pong {
    // UUID 资源编号
    string uuid = 1;

    // Pong 单个资源响应内容
    grpc_kit.api.known.example.v1.ExampleResponse pong = 2;
  }

  // Pong 返回创建的资源
  Pong pong = 1;

  // Content 多个资源响应内容（无分页属性）
  repeated grpc_kit.api.known.example.v1.ExampleResponse content = 2;

  // Ping 返回更新的资源
  grpc_kit.api.known.example.v1.ExampleResponse ping = 3;

  // Empty 返回空的内容
  google.protobuf.Empty empty = 4;
}
`,
	})

	t.files = append(t.files, &templateFile{
		name: fmt.Sprintf("api/%v/%v/%v/microservice.gateway.yaml",
			t.config.Global.ProductCode,
			t.config.Global.ShortName,
			t.config.Template.Service.APIVersion),
		parse: true,
		body: `
type: google.api.Service
config_version: 3

http:
  rules:
  - selector: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ title .Global.ProductCode }}{{ title .Global.ShortName }}.HealthCheck
    get: "/healthz"

  - selector: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ title .Global.ProductCode }}{{ title .Global.ShortName }}.Demo
    post: "/api/demo"
    body: "*"
    response_body: "pong"
    additional_bindings:
      - get: "/api/demo"
      - get: "/api/demo/{uuid}"
        response_body: "pong.pong"
      - put: "/api/demo/{uuid}"
        body: "ping"
        response_body: "ping"
      - delete: "/api/demo/{uuid}"
        response_body: "empty"
`,
	})

	t.files = append(t.files, &templateFile{
		name: fmt.Sprintf("api/%v/%v/%v/microservice.openapiv2.yaml",
			t.config.Global.ProductCode,
			t.config.Global.ShortName,
			t.config.Template.Service.APIVersion),
		parse: true,
		body: `
openapiOptions:
  # grpc.gateway.protoc_gen_openapiv2.options.Swagger
  # 对应 swagger 属性，一般不做更改
  file:
    - file: "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }}/microservice.proto"
      option:
        swagger: "2.0"
        info:
          title: "{{ .Global.ProductCode }}-{{ .Global.ShortName }}-{{ .Template.Service.APIVersion }}"
          contact:
            name: "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}"
            url: "http://{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}.{{ .Global.APIEndpoint }}"
          license:
            name: "Apache License 2.0"
            url: "https://github.com/grpc-kit/cli/blob/main/LICENSE"
          version: "0.0.0"
        host: "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}.{{ .Global.APIEndpoint }}"
        base_path: "/"
        schemes:
        - "HTTP"
        consumes:
        - "application/json"
        produces:
        - "application/json"
        securityDefinitions:
          security:
            BasicAuth:
              type: "TYPE_BASIC"
            ApiKeyAuth:
              type: "TYPE_API_KEY"
              in: "IN_HEADER"
              name: "Authorization: Bearer <token>"
        security:
        - securityRequirement:
            BasicAuth: {}
        - securityRequirement:
            ApiKeyAuth: {}
        responses:
          "4xx":
            description: "客户端参数异常"
            schema:
              jsonSchema:
                ref: ".grpc_kit.api.known.status.v1.ErrorResponse"
          "5xx":
            description: "服务端处理异常"
            schema:
              jsonSchema:
                ref: ".grpc_kit.api.known.status.v1.ErrorResponse"
        external_docs:
          description: 'Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}"'
          url: "https://grpc-kit.com"

  # grpc.gateway.protoc_gen_openapiv2.options.Operation
  # 对应 proto 中 service 的 rpc
  method:
    - method: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ title .Global.ProductCode }}{{ title .Global.ShortName }}.HealthCheck
      option:
        tags:
        - "internal"
        description: '请务删除！\n 接口格式：/healthz?service=test1.v1.opsaid\n 请求成功访问状态码200，且内容为：{"status": "SERVING"}'
        summary: "健康检测"
        responses:
          "200":
            examples:
              "application/json": '{"value": "the input value"}'

    - method: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ title .Global.ProductCode }}{{ title .Global.ShortName }}.Demo
      option:
        tags:
          - "demo"
        description: "这里做一些较长的使用描述\n 1. POST 用于创建资源，非幂等\n 2. GET 用于获取资源，幂等\n 3. PUT 用于更新资源，幂等\n 4. DELETE 用于删除资源，幂等"
        summary: "示例 RESTful/RPC 接口"
        responses:
          "204":
            examples:
              "application/json": '{}'

  # grpc.gateway.protoc_gen_openapiv2.options.Schema
  # 对应 proto 的 message
  message:
    - message: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.DemoRequest
      option:
        # 请求示例
        example: '{ "ping": { "name": "grpc-kit" } }'
        jsonSchema:
          # 必填字段
          required:
            - "ping"
          # 对结构体更详细的描述
          description: "结构体其他更详细的描述"
    - message: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.DemoResponse
      option:
        example: '{"uuid":"99feafb5-bed6-4daf-927a-69a2ab80c485", "pong": { "name": "grpc-kit" } }'

  # grpc.gateway.protoc_gen_openapiv2.options.JSONSchema
  # 对应 proto 的 message 下各属性
  field:
    - field: {{ .Global.ProductCode }}.{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.DemoRequest.uuid
      option:
        # 字段均不做 description 注解，在定义 proto 属性时添加
        # description: "请求的 ping 属性"
        default: "99feafb5-bed6-4daf-927a-69a2ab80c485"
`,
	})
}
