ARG TARGETPLATFORM
ARG GO_IMAGE=golang:1.25.13-bookworm@sha256:e401dae1bf814e29204a8cb7915682e1780951e609ca0dd8865ee1937f510c48

FROM --platform=$TARGETPLATFORM ${GO_IMAGE} AS builder

ARG HTTPS_PROXY

# 更换镜像地址
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources \
	&& sed -i 's#security.debian.org#mirrors.ustc.edu.cn/debian-security#g' /etc/apt/sources.list.d/debian.sources

# 设置时区与依赖的基础工具包
RUN apt-get update \
	&& apt-get install -y --no-install-recommends tzdata unzip docker.io \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && dpkg-reconfigure -f noninteractive tzdata \
	&& rm -rf /var/lib/apt/lists/*

# 设置环境变量
ENV GO111MODULE=on
# GOPROXY 可由构建参数覆盖：本地默认 goproxy.cn，GitHub Actions 传 proxy.golang.org
ARG GOPROXY=https://goproxy.cn
ENV GOPROXY=${GOPROXY}

# 拷贝当前源代码至 src 目录
WORKDIR /usr/local/src
COPY . .

# 编译 grpc-kit-cli
# 下载 protoc 相关二进制
# 安装其他依赖的 proto 定义文件
RUN go mod download \
	&& make build \
    && make protoc \
	&& make protoc-gen-go \
	&& make protoc-gen-go-grpc \
	&& make protoc-gen-grpc-gateway \
	&& make protoc-gen-openapiv2 \
    && git clone -b v0.3.1 --depth 1 https://github.com/grpc-kit/api.git $GOPATH/src/github.com/grpc-kit/api \
    && mkdir -p $GOPATH/src/github.com/googleapis \
    && git init $GOPATH/src/github.com/googleapis/googleapis \
	&& git -C $GOPATH/src/github.com/googleapis/googleapis remote add origin https://github.com/googleapis/googleapis.git \
	&& git -C $GOPATH/src/github.com/googleapis/googleapis fetch --depth 1 origin e0d0106516a5c613510533821e4508bc6c943b11 \
	&& git -C $GOPATH/src/github.com/googleapis/googleapis checkout --detach FETCH_HEAD \
    && git clone -b v2.30.0 --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway

# 用于 go 应用的编译
FROM --platform=$TARGETPLATFORM ${GO_IMAGE}

# 拷贝上阶段编译后的文件
COPY --from=builder /usr/local/src/build/grpc-kit-cli-* /go/bin/grpc-kit-cli
COPY --from=builder /go/bin/protoc /go/bin/protoc
COPY --from=builder /go/bin/protoc-gen-go /go/bin/protoc-gen-go
COPY --from=builder /go/bin/protoc-gen-go-grpc /go/bin/protoc-gen-go-grpc
COPY --from=builder /go/bin/protoc-gen-grpc-gateway /go/bin/protoc-gen-grpc-gateway
COPY --from=builder /go/bin/protoc-gen-openapiv2 /go/bin/protoc-gen-openapiv2
COPY --from=builder /go/include/google /go/include/google
COPY --from=builder /go/src/ /go/src/
COPY --from=builder /usr/bin/docker /usr/bin/docker

CMD ["/go/bin/grpc-kit-cli"]
