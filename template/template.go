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

package template

import (
	"fmt"

	"github.com/grpc-kit/cli/config"
	"github.com/grpc-kit/cli/template/service"
)

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
	default:
		return nil, fmt.Errorf("not support template type: %v", c.Global.Type)
	}

	return service.New(c)
}
