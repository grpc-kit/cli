.PHONY: build

# 全局通用变量
GO              := go
GORUN           := ${GO} run
GOPATH          := $(shell ${GO} env GOPATH)
GOARCH          ?= $(shell ${GO} env GOARCH)
GOBUILD         := ${GO} build

# 自动化版本号
GIT_COMMIT	:= $(shell git rev-parse HEAD)
GIT_BRANCH	:= $(shell git describe --tags --dirty --always)
BUILD_DATE	:= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_DATE	:= $(shell git --no-pager log -1 --format='%ct')
RELEASE_VERSION ?= $(shell cat VERSION)
BUILD_LD_FLAGS 	:= "-X 'github.com/grpc-kit/pkg/version.Appname=grpc-kit-cli' \
	-X 'github.com/grpc-kit/pkg/version.GitCommit=${GIT_COMMIT}' \
	-X 'github.com/grpc-kit/pkg/version.GitBranch=${GIT_BRANCH}' \
	-X 'github.com/grpc-kit/pkg/version.BuildDate=${BUILD_DATE}' \
	-X 'github.com/grpc-kit/pkg/version.CommitUnixTime=${COMMIT_DATE}' \
	-X 'github.com/grpc-kit/pkg/version.ReleaseVersion=${RELEASE_VERSION}'"

# 自定义变量
BUILD_GOOS		?= $(shell ${GO} env GOOS)

clean:
	@echo ">> clean build"
	@rm -rf build/

build: clean
	@mkdir build
	@GOOS=${BUILD_GOOS} ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/grpc-kit-cli main.go

