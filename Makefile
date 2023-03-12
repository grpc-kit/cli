# 全局通用变量
GO              := go
GORUN           := ${GO} run
GOPATH          := $(shell ${GO} env GOPATH)
GOOS            ?= $(shell ${GO} env GOOS)
GOARCH          ?= $(shell ${GO} env GOARCH)
GOBUILD         := ${GO} build

# 自动化版本号
GIT_COMMIT	:= $(shell git rev-parse HEAD)
GIT_BRANCH	:= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILD_DATE	:= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_DATE	:= $(shell git --no-pager log -1 --format='%ct')
RELEASE_VERSION ?= $(shell cat VERSION)
BUILD_LD_FLAGS 	:= "-X 'github.com/grpc-kit/pkg/vars.Appname=grpc-kit-cli' \
	-X 'github.com/grpc-kit/pkg/vars.GitCommit=${GIT_COMMIT}' \
	-X 'github.com/grpc-kit/pkg/vars.GitBranch=${GIT_BRANCH}' \
	-X 'github.com/grpc-kit/pkg/vars.BuildDate=${BUILD_DATE}' \
	-X 'github.com/grpc-kit/pkg/vars.CommitUnixTime=${COMMIT_DATE}' \
	-X 'github.com/grpc-kit/pkg/vars.ReleaseVersion=${RELEASE_VERSION}'"

# 自定义变量
BUILD_GOOS		?= $(shell ${GO} env GOOS)

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: clean ## Build binary files according to the target system arch.
	@mkdir build
	@GOOS=${BUILD_GOOS} ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli-${GOOS}-${GOARCH} main.go

.PHONY: build-all
build-all: clean ## Build all binaries that support the operating system.
	@mkdir build
	@GOOS=linux GOARCH=amd64 ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli-linux-arm64 main.go
	@GOOS=darwin GOARCH=amd64 ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli-darwin-arm64 main.go

.PHONY: docker-build
docker-build: ## Build docker image with the application.
	@echo ">> docker build"
	@# TODO; platform 参数不支持填写多个？
	@docker buildx build --platform linux/amd64 ./ --manifest ccr.ccs.tencentyun.com/grpc-kit/cli:${RELEASE_VERSION}
	@docker buildx build --platform linux/arm64 ./ --manifest ccr.ccs.tencentyun.com/grpc-kit/cli:${RELEASE_VERSION}
	@docker buildx build --platform darwin/amd64 ./ --manifest ccr.ccs.tencentyun.com/grpc-kit/cli:${RELEASE_VERSION}
	@docker buildx build --platform darwin/arm64 ./ --manifest ccr.ccs.tencentyun.com/grpc-kit/cli:${RELEASE_VERSION}

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
