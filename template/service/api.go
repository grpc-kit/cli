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

// 引入依赖的外部proto文件
import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "github.com/grpc/grpc-proto/grpc/health/v1/health.proto";
import "github.com/googleapis/googleapis/google/api/annotations.proto";
import "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options/annotations.proto";

// 同组RPC方法对应一个proto文件，以该组RPC名称的小写字母为文件名
import "api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }}/demo.proto";

// 结合本项目，推荐做以下设置
option (gogoproto.marshaler_all) = true;
option (gogoproto.sizer_all) = true;
option (gogoproto.unmarshaler_all) = true;
option (gogoproto.goproto_registration) = true;
option (gogoproto.goproto_getters_all) = true;
option (gogoproto.goproto_unrecognized_all) = false;
option (gogoproto.goproto_unkeyed_all) = false;
option (gogoproto.goproto_sizecache_all) = false;
option (gogoproto.equal_all) = true;
option (gogoproto.compare_all) = true;
option (gogoproto.messagename_all) = false;

// 转换为swagger接口文档的相关设置
option (grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger) = {
    host: "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}.{{ .Global.APIEndpoint }}",
    info: {
        title: "{{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}",
        contact: {
            name: "gRPC Kit",
            url: "https://grpc-kit.com"
        },
        license:{
            name: "Apache License 2.0"
        }
        version: "v0.0.0",
    },
    security_definitions: {
        security: {
            key: "BasicAuth",
            value: {
                type: TYPE_BASIC
            }
        }
    },
    security: {
        security_requirement: {
            key: "BasicAuth",
            value: {}
        }
    },
    responses: {
        key: "500",
        value: {
            description: '{"code": 500, "error": "internal error", "message": "internal error", "details": []}'
        }
    }
};

// 该微服务支持的RPC方法定义
service {{ title .Global.ProductCode }}{{ title .Global.ShortName }} {
    rpc HealthCheck(grpc.health.v1.HealthCheckRequest) returns (grpc.health.v1.HealthCheckResponse) {
        option (google.api.http) = {
            get: "/healthz"
        };
        option (grpc.gateway.protoc_gen_swagger.options.openapiv2_operation) = {
            tags: "internal"
            summary: "健康检测"
            description: '请务删除！\n 接口格式：/healthz?service={{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}\n 请求成功访问状态码200，且内容为：{"status": "SERVING"}'
        };
    };

    rpc Demo(DemoRequest) returns (DemoResponse) {
        option (google.api.http) = {
            post: "/demo"
            body: "*"
            response_body: "pong"
            additional_bindings {
                get: "/demo"
                response_body: "content"
            }
            additional_bindings {
                get: "/demo/{uuid}"
                response_body: "pong.pong"
            }
            additional_bindings {
                put: "/demo/{uuid}"
                body: "ping"
                response_body: "ping"
            }
            additional_bindings {
                delete: "/demo/{uuid}"
                response_body: "empty"
            }
        };
        option (grpc.gateway.protoc_gen_swagger.options.openapiv2_operation) = {
            tags: "demo",
            summary: "示例RESTful/RPC接口",
            description: "这里做一些较长的使用描述\n 1. POST 用于创建资源，非幂等\n 2. GET 用于获取资源，幂等\n 3. PUT 用于更新资源，幂等\n 4. DELETE 用于删除资源，幂等",
            deprecated: false,
            responses: {
                key: "204",
                value: {
                    description: "no content"
                }
            }
        };
    }

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

// 引入google公共类型
import "google/protobuf/empty.proto";

// 引入第三方依赖的proto文件
import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options/annotations.proto";

// 引入项目通用的proto文件
import "github.com/grpc-kit/api/proto/v1/example.proto";

// 结合本项目，推荐做以下设置
option (gogoproto.marshaler_all) = true;
option (gogoproto.sizer_all) = true;
option (gogoproto.unmarshaler_all) = true;
option (gogoproto.goproto_registration) = true;
option (gogoproto.goproto_getters_all) = true;
option (gogoproto.goproto_unrecognized_all) = false;
option (gogoproto.goproto_unkeyed_all) = false;
option (gogoproto.goproto_sizecache_all) = false;
option (gogoproto.equal_all) = true;
option (gogoproto.compare_all) = true;
option (gogoproto.messagename_all) = false;

message DemoRequest {
    option (grpc.gateway.protoc_gen_swagger.options.openapiv2_schema) = {
        json_schema: {
            description: "Demo方法请求可使用的接口参数",
        },
        example: { value: '{ "ping": { "name": "grpc-kit" } }' }
    };

    // UUID 资源编号
    string uuid = 1;

    // Ping 资源内容
    grpc.kit.api.proto.v1.ExampleRequest ping = 2;
}

message DemoResponse {
    option (grpc.gateway.protoc_gen_swagger.options.openapiv2_schema) = {
        json_schema: {
            description: "Demo方法响应的具体内容",
        },
        example: { value: '{"uuid":"99feafb5-bed6-4daf-927a-69a2ab80c485", "pong": { "name": "grpc-kit" } }' }
    };

    message Pong {
        // UUID 资源编号
        string uuid = 1;

        // Pong 单个资源响应内容
        grpc.kit.api.proto.v1.ExampleResponse pong = 2;
    }

    // Pong 返回创建的资源
    Pong pong = 1;

    // Content 多个资源响应内容（无分页属性）
    repeated grpc.kit.api.proto.v1.ExampleResponse content = 2;

    // Ping 返回更新的资源
    grpc.kit.api.proto.v1.ExampleResponse ping = 3;

    // Empty 返回空的内容
    google.protobuf.Empty empty = 4;
}
`,
	})
}
