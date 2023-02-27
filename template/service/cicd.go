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

func (t *templateService) fileDirectoryCICD() {
	t.files = append(t.files, &templateFile{
		name:  ".gitlab/workflows/grpc-kit.yml",
		parse: true,
		body: `
# 默认全局配置
default:
  tags:
    - grpc-kit

# 默认全局变量
variables:
  CGO_ENABLED: "0"
  GIT_SSL_NO_VERIFY: "true"

# 流水线各阶段
stages:
  - pre
  - test
  - build
  - deploy
  - production

# 代码风格、格式检测
go-lint:
  stage: pre
  script:
    - make lint

# 依赖的 protoc 二进制检测
check-protoc:
  stage: pre
  script:
    - which go
    - which protoc
    - which protoc-gen-gogo

# 业务单元测试
unit-tests:
  stage: test
  script:
    - make test

# 编译二进制文件
build-binary:
  stage: build
  script:
    - make build
  artifacts:
    paths:
      - build/
    expire_in: 24h

# 打包容器镜像
build-container:
  stage: build
  script:
    - echo "package tar"
  artifacts:
    paths:
      - build/
    expire_in: 24h
  only:
    - main
  when: manual

# 打成各种安装包，如：tar、rpm、deb
build-package:
  stage: build
  script:
    - echo "package tar"
    - echo "package rpm"
    - echo "package deb"
  only:
    - main
  when: manual

# 部署测试环境
deploy-test:
  stage: deploy
  script:
    - echo "deploy dev"

# 部署正式环境，手工确认
production:
  stage: production
  script:
    - echo "deploy production"
  only:
    - main
  when: manual
`,
	})
}
