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
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/grpc-kit/cli/config"
	"github.com/grpc-kit/pkg/file"
)

type templateService struct {
	config config.Config
}

// New 实例化
func newService(c config.Config) (*templateService, error) {
	return &templateService{config: c}, nil
}

// Generate 生产代码模版
func (t *templateService) Generate() error {
	err := fs.WalkDir(Assets, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 过滤掉目录，如："."
		if d.IsDir() {
			return nil
		}

		content, err := Assets.ReadFile(filePath)
		if err != nil {
			return err
		}

		fileBody := string(content)

		if strings.HasSuffix(filePath, ".tmpl") {
			fileBody, err = file.ParseExecute(fileBody, t.config)
			if err != nil {
				return err
			}
		}

		var filePerm os.FileMode
		filePerm = 0666

		// 如果是脚本文件，则设置可执行权限
		if strings.HasSuffix(filePath, ".sh") || strings.HasSuffix(filePath, ".sh.tmpl") {
			filePerm = 0755
		}

		fileName := strings.TrimPrefix(filePath, "service/")
		fileName = strings.TrimSuffix(fileName, ".tmpl")

		if strings.HasPrefix(fileName, "api/") {
			fileName = fmt.Sprintf("api/%v/%v/%v/%v",
				t.config.Global.ProductCode,
				t.config.Global.ShortName,
				t.config.Template.Service.APIVersion,
				path.Base(fileName))
		}

		return file.WriteString(
			path.Join(t.config.Global.ShortName, fileName),
			strings.TrimLeft(fileBody, "\n"),
			filePerm)
	})

	return err
}
