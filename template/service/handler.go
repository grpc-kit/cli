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

func (t *templateService) fileDirectoryHandler() {
	t.files = append(t.files, &templateFile{
		name:  "handler/microservice.go",
		parse: true,
		body: `
// Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.

package handler

import (
	"github.com/grpc-kit/pkg/cfg"
	"github.com/grpc-kit/pkg/rpc"
	"github.com/sirupsen/logrus"

	"{{ .Global.Repository }}/modeler"
)

// Microservice 该微服务的结构
type Microservice struct {
	code    string                  // 服务代码
	server  *rpc.Server             // 服务定义
	client  *rpc.Client             // 服务调用
	logger  *logrus.Entry           // 全局日志
	baseCfg *cfg.LocalConfig        // 基础配置
	thisCfg *modeler.IndependentCfg // 个性配置
}

// NewMicroservice 全局只实例化一次
func NewMicroservice(lc *cfg.LocalConfig) (*Microservice, error) {
	// 基础配置初始化
	if err := lc.Init(); err != nil {
		return nil, err
	}

	m := &Microservice{
		code:    lc.Services.ServiceCode,
		logger:  lc.GetLogger(),
		baseCfg: lc,
		thisCfg: &modeler.IndependentCfg{},
	}

	// 自定义配置实例化
	if err := lc.GetIndependent(m.thisCfg); err != nil {
		return m, err
	}

	// RPC客户端服务端
	client, err := lc.GetRPCClient()
	if err != nil {
		return m, err
	}
	server, err := lc.GetRPCServer()
	if err != nil {
		return m, err
	}
	m.client = client
	m.server = server

	// 其他个性扩展逻辑
	if err := m.privateExtended(); err != nil {
		return m, err
	}

	// 自定义配置初始化
	if err := m.thisCfg.Init(); err != nil {
		return m, err
	}

	return m, nil
}
`,
	})

	t.files = append(t.files, &templateFile{
		name: "handler/private.go",
		body: `
package handler

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc"
)

func (m *Microservice) privateExtended() error {
	clientOpts := m.baseCfg.GetClientDialOption()
	clientUnaryHandlers := m.baseCfg.GetClientUnaryInterceptor()
	clientStreamHandlers := m.baseCfg.GetClientStreamInterceptor()

	m.client.UseDialOption(clientOpts...).
		UseUnaryInterceptor(clientUnaryHandlers...).
		UseStreamInterceptor(clientStreamHandlers...)

	m.server.UseServerOption(m.baseCfg.GetUnaryInterceptor(m.privateUnaryServerInterceptor()...),
		m.baseCfg.GetStreamInterceptor(m.privateStreamServerInterceptor()...))

	return nil
}

func (m *Microservice) privateUnaryServerInterceptor() []grpc.UnaryServerInterceptor {
	return nil
}

func (m *Microservice) privateStreamServerInterceptor() []grpc.StreamServerInterceptor {
	return nil
}

func (m *Microservice) privateHTTPHandle(mux *http.ServeMux) {
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "")
	})
}
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "handler/register.go",
		parse: true,
		body: `
// Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.

package handler

import (
	"context"
	"net/http"

	"{{ .Global.Repository }}/public/doc"
	pb "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }}"
)

// Register 用于服务启动前环境准备
func (m *Microservice) Register(ctx context.Context) error {
	pb.Register{{ title .Global.ProductCode }}{{ title .Global.ShortName }}Server(m.server.Server(), m)

	// 注册服务信息
	mux, err := m.baseCfg.Register(ctx, pb.Register{{ title .Global.ProductCode }}{{ title .Global.ShortName }}HandlerFromEndpoint)
	if err != nil {
		return err
	}

	// 注册API文档
    mux.Handle("/openapi-spec/", http.FileServer(http.FS(doc.Assets)))

	// 这里添加其他自定义实现
	m.privateHTTPHandle(mux)

	// 注册HTTP网关
	if err := m.server.RegisterGateway(mux); err != nil {
		return err
	}

	// 开启gRPC与HTTP服务
	if err := m.server.StartBackground(); err != nil {
		return err
	}

	return nil
}
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "handler/rpc_demo.go",
		parse: true,
		body: `
package handler

import (
	"context"

    "github.com/gogo/protobuf/types"
    "github.com/grpc-kit/pkg/api"

	pb "{{ .Global.Repository }}/api/{{ .Global.ProductCode }}/{{ .Global.ShortName }}/{{ .Template.Service.APIVersion }}"
)

// Demo test
func (m Microservice) Demo(ctx context.Context, req *pb.DemoRequest) (*pb.DemoResponse, error) {
	m.logger.Warnf("test demo warn: %v", "func Demo")

	result := &pb.DemoResponse{
	    // GET /api/demo
	    Content:  []*api.ExampleResponse{
	        {Name: "grpc-kit-cli"},
            {Name: "grpc-kit-cfg"},
            {Name: "grpc-kit-pkg"},
            {Name: "grpc-kit-api"},
            {Name: "grpc-kit-web"},
            {Name: "grpc-kit-doc"},
        },
        Ping:  &api.ExampleResponse{},
        // POST /api/demo
        // GET /api/demo/{uuid}
        Pong: &pb.DemoResponse_Pong{
	        Uuid: "99feafb5-bed6-4daf-927a-69a2ab80c485",
	        Pong: &api.ExampleResponse{},
        },
        // DELETE /api/demo/{uuid}
        Empty: &types.Empty{},
    }

    if req.Ping != nil {
        result.Ping.Name = req.Ping.Name
        result.Pong.Pong.Name = req.Ping.Name
    }

    if req.Uuid == "99feafb5-bed6-4daf-927a-69a2ab80c485" {
        result.Pong.Pong.Name = "grpc-kit"
    }

    return result, nil
}
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "handler/rpc_internal.go",
		parse: true,
		body: `
// Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.

package handler

import (
    "context"

    "github.com/grpc-kit/pkg/errors"
    hz "google.golang.org/grpc/health/grpc_health_v1"
)

// HealthCheck 用于健康检测
func (m Microservice) HealthCheck(ctx context.Context, req *hz.HealthCheckRequest) (*hz.HealthCheckResponse, error) {
    if req.Service == m.code {
        return &hz.HealthCheckResponse{
            Status: hz.HealthCheckResponse_SERVING,
        }, nil
    }

    return nil, errors.NotFound(ctx).WithMessage("unknown service").Err()
}
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "handler/shutdown.go",
		parse: true,
		body: `
// Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.

package handler

import (
	"context"
	"time"
)

// Shutdown 优雅关闭gRPC与HTTP服务
func (m *Microservice) Shutdown(ctx context.Context) error {
	m.logger.Warnf("Shutdown server begin")

	if err := m.baseCfg.Deregister(); err != nil {
		return err
	}

	// 最长等待关闭的时间，例如超过30秒则强制关闭gateway
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// 阻塞等待，直到所有连接正常或超时退出
	if err := m.server.Shutdown(ctx); err != nil {
		return err
	}

	m.logger.Warnf("Shutdown server end")
	return nil
}
`,
	})
}
