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
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"

	"github.com/grpc-kit/cli/config"
)

type templateService struct {
	config config.Config
}

type generatedFile struct {
	name string
	body []byte
	mode fs.FileMode
}

// New 实例化
func newService(c config.Config) (*templateService, error) {
	return &templateService{config: c}, nil
}

// Generate 生产代码模版
func (t *templateService) Generate() error {
	return t.GenerateTo(t.config.Global.ShortName)
}

// GenerateTo renders the complete service into a sibling temporary directory
// and publishes it with one rename. An existing output path is never modified.
func (t *templateService) GenerateTo(output string) error {
	files, err := t.render()
	if err != nil {
		return err
	}

	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("output directory must not be empty")
	}
	target, err := filepath.Abs(filepath.Clean(output))
	if err != nil {
		return fmt.Errorf("resolve output directory %q: %w", output, err)
	}
	if _, err = os.Lstat(target); err == nil {
		return fmt.Errorf("output path already exists: %s", output)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect output path %q: %w", output, err)
	}

	parent := filepath.Dir(target)
	if err = os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("create output parent %q: %w", parent, err)
	}
	stage, err := os.MkdirTemp(parent, ".grpc-kit-cli-*")
	if err != nil {
		return fmt.Errorf("create staging directory: %w", err)
	}
	defer os.RemoveAll(stage)

	for _, file := range files {
		destination := filepath.Join(stage, filepath.FromSlash(file.name))
		if err = os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create directory for %s: %w", file.name, err)
		}
		if err = os.WriteFile(destination, file.body, file.mode); err != nil {
			return fmt.Errorf("write %s: %w", file.name, err)
		}
	}
	if err = os.Chmod(stage, 0o755); err != nil {
		return fmt.Errorf("set output directory permissions: %w", err)
	}
	if err = os.Rename(stage, target); err != nil {
		return fmt.Errorf("publish output directory %q: %w", output, err)
	}
	return nil
}

func (t *templateService) render() ([]generatedFile, error) {
	files := make([]generatedFile, 0)
	err := fs.WalkDir(Assets, "service", func(filePath string, d fs.DirEntry, err error) error {
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
			fileBody, err = renderTemplate(fileBody, t.config)
			if err != nil {
				return err
			}
		}

		filePerm := fs.FileMode(0o666)

		// 如果是脚本文件，则设置可执行权限
		if strings.HasSuffix(filePath, ".sh") || strings.HasSuffix(filePath, ".sh.tmpl") {
			filePerm = 0o755
		}

		fileName := strings.TrimPrefix(filePath, "service/")
		fileName = strings.TrimSuffix(fileName, ".tmpl")

		if strings.HasPrefix(fileName, "api/") {
			fileName = fmt.Sprintf("api/%v/%v/%v/%v",
				t.config.Global.ProductCode,
				t.config.Global.ShortName,
				t.config.Template.Service.APIVersion,
				filepath.Base(fileName))
		}

		cleanName := filepath.ToSlash(filepath.Clean(filepath.FromSlash(fileName)))
		if cleanName == "." || cleanName == ".." || strings.HasPrefix(cleanName, "../") || filepath.IsAbs(fileName) {
			return fmt.Errorf("template asset resolves outside output directory: %s", filePath)
		}
		files = append(files, generatedFile{
			name: cleanName,
			body: []byte(strings.TrimLeft(fileBody, "\n")),
			mode: filePerm,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func renderTemplate(body string, data any) (string, error) {
	functions := texttemplate.FuncMap{
		"title":     strings.Title,
		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"toTitle":   strings.ToTitle,
		"trimSpace": strings.TrimSpace,
	}
	tmpl, err := texttemplate.New("").Funcs(functions).Option("missingkey=error").Parse(body)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err = tmpl.Execute(&output, data); err != nil {
		return "", err
	}
	return output.String(), nil
}
