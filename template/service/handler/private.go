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

func (m *Microservice) privateHTTPHandle(mux *http.ServeMux) error {
	// 这里属于自定义 http 接口，访问 /favicon.ico 不会产生链路数据
	// 如需捕获链路数据，参考文档：https://grpc-kit.com/docs/spec-cfg/observables/
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "")
	})

	return nil
}
