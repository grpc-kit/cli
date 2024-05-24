package template

import (
	"embed"
	"fmt"

	"github.com/grpc-kit/cli/config"
)

//go:embed service/*
var Assets embed.FS

const (
	// TypeService 服务端模版
	TypeService = "service"
)

// Template 代码模版的接口定义
type Template interface {
	Generate() error
}

// New xx
func New(c config.Config) (Template, error) {
	switch c.Global.Type {
	case TypeService:
		s, err := newService(c)
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, fmt.Errorf("not support template type: %v", c.Global.Type)
	}
}
