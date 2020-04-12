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

import (
	"os"
	"path"
	"strings"

	"github.com/grpc-kit/cli/config"
	"github.com/grpc-kit/pkg/file"
)

type templateService struct {
	config config.Config
	files  []*templateFile
}

type templateFile struct {
	name  string
	body  string
	parse bool
}

// New 实例化
func New(c config.Config) (*templateService, error) {
	return &templateService{config: c}, nil
}

// Generate 生产代码模版
func (t *templateService) Generate() error {
	t.fileDirectoryRoot()
	t.fileDirectoryApi()
	t.fileDirectoryCmd()
	t.fileDirectoryConfig()
	t.fileDirectoryHandler()
	t.fileDirectoryModeler()
	t.fileDirectoryScripts()

	var err error
	for _, f := range t.files {
		fileBody := f.body

		if f.parse {
			fileBody, err = file.ParseExecute(f.body, t.config)
			if err != nil {
				return err
			}
		}

		var perm os.FileMode
		perm = 0666

		if strings.HasSuffix(f.name, "sh") {
			perm = 0755
		}

		err := file.WriteString(path.Join(t.config.Global.ShortName, f.name),
			strings.TrimSpace(fileBody),
			perm)
		if err != nil {
			return err
		}

	}

	return nil
}
