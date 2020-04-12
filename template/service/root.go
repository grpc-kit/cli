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

func (t *templateService) fileDirectoryRoot() {
	t.files = append(t.files, &templateFile{
		name: ".gitignore",
		body: `
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories
vendor/

# Config file and certificate
config/*.pem
config/app-dev-*.toml
config/app-test-*.toml
config/app-prod-*.toml

# Others
.swp
.bak
.idea/
build/
backup/
scripts/env-dev-*.sh
scripts/env-test-*.sh
scripts/env-prod-*.sh
api/assets_vfsdata.go
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "go.mod",
		parse: true,
		body: `
module {{ .Global.GitDomain }}/{{ .Global.ProductCode }}/{{ .Global.ShortName }}

go 1.13

require (
	google.golang.org/grpc v1.26.0
)
`,
	})

	t.files = append(t.files, &templateFile{
		name: "CHANGELOG.md",
		body: `
# Changelog

名称 | 说明
------------|----------
Added       | 添加新功能
Changed     | 功能的变更
Deprecated  | 未来会删除
Removed     | 之前为Deprecated状态，此版本被移除
Fixed       | 功能的修复
Security    | 有关安全问题的修复

## [Unreleased]

### Added

- 已完成的功能，未正式发布

## [0.1.0] - 2020-03-28

### Added

- 首次发布
`,
	})

	t.files = append(t.files, &templateFile{
		name: "Makefile",
		body: `
.PHONY: precheck clean version proto

include scripts/env

# 全局通用变量
GO              := go
GORUN           := ${GO} run
GOPATH          := $(shell ${GO} env GOPATH)
GOARCH          ?= $(shell ${GO} env GOARCH)
GOBUILD         := ${GO} build
WORKDIR         := $(shell pwd)
GOHOSTOS        := $(shell ${GO} env GOHOSTOS)
SHORTNAME       := $(shell basename $(shell pwd))
NAMESPACE       ?= $(shell basename $(shell cd ../;pwd))

# 自动化版本号
GIT_COMMIT      := $(shell git rev-parse HEAD 2>/dev/null)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILD_DATE      := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_DATE     := $(shell git --no-pager log -1 --format='%ct' 2>/dev/null)
PREFIX_VERSION  := $(shell ./scripts/version.sh prefix)
RELEASE_VERSION ?= $(shell ./scripts/version.sh release)
BUILD_LD_FLAGS  := "-X 'github.com/grpc-kit/pkg/version.CliVersion=${CLI_VERSION}' \
                -X 'github.com/grpc-kit/pkg/version.GitCommit=${GIT_COMMIT}' \
                -X 'github.com/grpc-kit/pkg/version.GitBranch=${GIT_BRANCH}' \
                -X 'github.com/grpc-kit/pkg/version.BuildDate=${BUILD_DATE}' \
                -X 'github.com/grpc-kit/pkg/version.CommitUnixTime=${COMMIT_DATE}' \
                -X 'github.com/grpc-kit/pkg/version.ReleaseVersion=${RELEASE_VERSION}'"

# 构建Docker容器变量
IMAGE_FROM      ?= scratch
IMAGE_HOST      ?= hub.docker.com
BUILD_GOOS      ?= $(shell ${GO} env GOOS)

precheck:
	@echo ">> precheck environment"
	@./scripts/precheck.sh

run: generate
	@${GORUN} -ldflags ${BUILD_LD_FLAGS} cmd/server/main.go -c config/app-dev-local.toml

build: clean generate
	@mkdir build
	@mkdir build/deploy
	@GOOS=${BUILD_GOOS} ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/service cmd/server/main.go

clean:
	@echo ">> clean build"
	@rm -rf build/

generate: precheck
	@echo ">> generation release version"
	@./scripts/version.sh update

	@echo ">> generation assets to static code"
	@${GO} generate api/generate.go >> /dev/null

	@echo ">> generation code from proto files"
	@./scripts/genproto.sh

build-docker: build
	@echo ">> build docker"

	@IMAGE_FROM=${IMAGE_FROM} \
	IMAGE_HOST=${IMAGE_HOST} \
	NAMESPACE=${NAMESPACE} \
	SHORTNAME=${SHORTNAME} \
	RELEASE_VERSION=${RELEASE_VERSION} \
	./scripts/docker.sh
`,
	})
}
