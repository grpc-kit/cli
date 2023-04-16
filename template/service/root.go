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
config/*.crt
config/*.key
config/app-dev-*
config/app-test-*
config/app-prod-*

# Others
.swp
.bak
.idea/
build/
backup/
public/doc/openapi-spec/microservice.swagger.json
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "go.mod",
		parse: true,
		body: `
module {{ .Global.Repository }}

go 1.18

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-kit/pkg v0.3.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
)
replace google.golang.org/grpc => google.golang.org/grpc v1.38.0
`,
	})

	t.files = append(t.files, &templateFile{
		name: "CHANGELOG.md",
		body: `
# Changelog

| 名称        | 说明                           |
|------------|--------------------------------|
| Added      | 添加新功能                       |
| Changed    | 功能的变更                       |
| Deprecated | 未来会删除                       |
| Removed    | 之前为Deprecated状态，此版本被移除 |
| Fixed      | 功能的修复                       |
| Security   | 有关安全问题的修复                |

## [Unreleased]

### Added

- 已完成的功能，未正式发布

## [0.1.0] - 2020-03-28

### Added

- 首次发布
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "Makefile",
		parse: true,
		body: `
include scripts/env

# 全局通用变量
GO              := go
GORUN           := ${GO} run
GOPATH          := $(shell ${GO} env GOPATH)
GOOS            ?= $(shell ${GO} env GOOS)
GOARCH          ?= $(shell ${GO} env GOARCH)
GOBUILD         := ${GO} build
WORKDIR         := $(shell pwd)
GOHOSTOS        := $(shell ${GO} env GOHOSTOS)
SHORTNAME       ?= $(shell basename $(shell pwd))
NAMESPACE       ?= $(shell basename $(shell cd ../;pwd))

# 自动化版本号
GIT_COMMIT      := $(shell git rev-parse HEAD 2>/dev/null)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILD_DATE      := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_DATE     := $(shell git --no-pager log -1 --format='%ct' 2>/dev/null)
PREFIX_VERSION  := $(shell ./scripts/version.sh prefix)
RELEASE_VERSION ?= $(shell ./scripts/version.sh release)
BUILD_LD_FLAGS  := "-X 'github.com/grpc-kit/pkg/vars.Appname={{ .Global.ProductCode }}-{{ .Global.ShortName }}-{{ .Template.Service.APIVersion }}' \
                -X 'github.com/grpc-kit/pkg/vars.CliVersion=${CLI_VERSION}' \
                -X 'github.com/grpc-kit/pkg/vars.GitCommit=${GIT_COMMIT}' \
                -X 'github.com/grpc-kit/pkg/vars.GitBranch=${GIT_BRANCH}' \
                -X 'github.com/grpc-kit/pkg/vars.BuildDate=${BUILD_DATE}' \
                -X 'github.com/grpc-kit/pkg/vars.CommitUnixTime=${COMMIT_DATE}' \
                -X 'github.com/grpc-kit/pkg/vars.ReleaseVersion=${RELEASE_VERSION}'"

# 构建Docker容器变量
BUILD_GOOS      ?= $(shell ${GO} env GOOS)
IMAGE_FROM      ?= scratch
IMAGE_HOST      ?= hub.docker.com
IMAGE_NAME      ?= ${IMAGE_HOST}/${NAMESPACE}/${SHORTNAME}
IMAGE_VERSION   ?= ${RELEASE_VERSION}

# 部署与运行相关变量
BUILD_ENV       ?= local
DEPLOY_ENV      ?= dev

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Check

.PHONY: precheck
precheck: ## Check environment.
	@echo ">> precheck environment"
	@./scripts/precheck.sh

##@ Development

.PHONY: generate
manifests: ## Generate deployment manifests files.
	@NAMESPACE=${NAMESPACE} \
		IMAGE_NAME=${IMAGE_NAME} \
		IMAGE_VERSION=${IMAGE_VERSION} \
		BUILD_ENV=${BUILD_ENV} ./scripts/manifests.sh ${DEPLOY_ENV}

generate: precheck ## Generate code from proto files.
	@echo ">> generation release version"
	@./scripts/version.sh update

	@echo ">> generation code from proto files"
	@./scripts/genproto.sh

.PHONY: lint
lint: generate ## Run go fmt and vet against code.
	@${GO} mod tidy
	@${GO} fmt ./...
	@${GO} vet ./...

.PHONY: test
test: generate ## Run unit tests.
	@${GO} mod tidy
	@${GO} test ./...

##@ Build

.PHONY: build
build: clean generate ## Build binary files according to the target system arch.
	@mkdir build
	@mkdir build/deploy
	@${GO} mod tidy
	@GOOS=${BUILD_GOOS} ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/service cmd/server/main.go

.PHONY: run
run: generate ## Run a application from your host.
	@${GO} mod tidy
	@${GORUN} -ldflags ${BUILD_LD_FLAGS} cmd/server/main.go -c config/app-dev-local.yaml

.PHONY: docker-run
docker-run: ## Run a application from your docker.
	@./scripts/docker.sh run

.PHONY: docker-build
docker-build: build manifests ## Build docker image with the application.
	@echo ">> docker build"
	@IMAGE_FROM=${IMAGE_FROM} \
		IMAGE_HOST=${IMAGE_HOST} \
		NAMESPACE=${NAMESPACE} \
		SHORTNAME=${SHORTNAME} \
		IMAGE_VERSION=${IMAGE_VERSION} ./scripts/docker.sh build

.PHONY: docker-push
docker-push: ## Push docker image with the application.
	@echo ">> docker push"
	@IMAGE_HOST=${IMAGE_HOST} \
		NAMESPACE=${NAMESPACE} \
		SHORTNAME=${SHORTNAME} \
		IMAGE_VERSION=${IMAGE_VERSION} ./scripts/docker.sh push

##@ Build Dependencies

.PHONY: protoc
protoc: ## Download protoc locally if necessary.
	@echo ">> download binary protoc"
	@./scripts/binaries.sh protoc

.PHONY: protoc-gen-go
protoc-gen-go: ## Download protoc-gen-go locally if necessary.
	@echo ">> download binary protoc-gen-go"
	@./scripts/binaries.sh protoc-gen-go

.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: ## Download protoc-gen-go-grpc locally if necessary.
	@echo ">> download binary protoc-gen-go-grpc"
	@./scripts/binaries.sh protoc-gen-go-grpc

.PHONY: protoc-gen-grpc-gateway
protoc-gen-grpc-gateway: ## Download protoc-gen-grpc-gateway locally if necessary.
	@echo ">> download binary protoc-gen-grpc-gateway"
	@./scripts/binaries.sh protoc-gen-grpc-gateway

.PHONY: protoc-gen-openapiv2
protoc-gen-openapiv2: ## Download protoc-gen-openapiv2 locally if necessary.
	@echo ">> download binary protoc-gen-openapiv2"
	@./scripts/binaries.sh protoc-gen-openapiv2

##@ Clean

.PHONY: clean
clean: ## Clean build.
	@echo ">> clean build"
	@rm -rf build/
`,
	})

	t.files = append(t.files, &templateFile{
		name: "VERSION",
		body: `0.1.0`,
	})

	t.files = append(t.files, &templateFile{
		name:  "README.md",
		parse: true,
		body: `
# {{ .Global.ShortName }}.{{ .Template.Service.APIVersion }}.{{ .Global.ProductCode }}

## Description

// TODO(user): Add simple overview of use/purpose

## Getting Started

// TODO(user): Add quick start

`,
	})

	t.files = append(t.files, &templateFile{
		name: "Dockerfile",
		body: `
FROM scratch

WORKDIR /opt

COPY build/service /opt/service
COPY config/app-mini.yaml /opt/config/app.yaml

EXPOSE 10080/tcp
EXPOSE 10081/tcp

ENTRYPOINT [ "/opt/service" ]
CMD [ "--config", "/opt/config/app.yaml" ]
`,
	})
}
