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

func (t *templateService) fileDirectoryModeler() {
	t.files = append(t.files, &templateFile{
		name: "modeler/independent_cfg.go",
		body: `
package modeler

// IndependentCfg 个性配置
type IndependentCfg struct {
	Name string ` + "`mapstructure:\"name\"`" + `
}

// Init 用于初始化实例
func (i *IndependentCfg) Init() error {
	// 业务代码

	return nil
}
`,
	})

	t.files = append(t.files, &templateFile{
		name: "modeler/independent_cfg_test.go",
		body: `
package modeler

import "testing"

func TestIndependentCfg(t *testing.T) {
	i := &IndependentCfg{}
	if err := i.Init(); err != nil {
		t.Error(err)
	}
}
`,
	})
}
