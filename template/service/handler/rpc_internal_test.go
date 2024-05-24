package handler

import (
	"context"
	"testing"

	statusv1 "github.com/grpc-kit/pkg/api/known/status/v1"
)

func TestInternal(t *testing.T) {
	req := &statusv1.HealthCheckRequest{
		Service: m.baseCfg.Services.ServiceCode,
	}

	_, err := m.HealthCheck(context.TODO(), req)
	if err != nil {
		t.Errorf("HealthCheck test fail: %v", err)
	}
}
