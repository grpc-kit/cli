ARG TARGETPLATFORM

FROM --platform=$TARGETPLATFORM golang:1.18.10-bullseye as builder

# 更换镜像地址
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

# 设置时区与依赖的基础工具包
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && dpkg-reconfigure -f noninteractive tzdata \
    && apt-get install -y unzip docker.io

# 设置环境变量
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

# 拷贝当前源代码至 src 目录
WORKDIR /usr/local/src
COPY . .

# 编译 grpc-kit-cli
# 下载 protoc 相关二进制
# 安装其他依赖的 proto 定义文件
RUN go mod tidy \
    && make build \
    && make protoc \
	&& make protoc-gen-go \
	&& make protoc-gen-go-grpc \
	&& make protoc-gen-grpc-gateway \
	&& make protoc-gen-openapiv2 \
    && git clone -b v0.3.0 --depth 1 https://github.com/grpc-kit/api.git $GOPATH/src/github.com/grpc-kit/api \
    && git clone --depth 1 https://github.com/googleapis/googleapis.git $GOPATH/src/github.com/googleapis/googleapis \
    && git clone -b v2.15.2 --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway

# 用于容器内构建镜像
FROM --platform=$TARGETPLATFORM gcr.io/kaniko-project/executor:v1.9.1 as kaniko

# 用于 go 应用的编译
FROM --platform=$TARGETPLATFORM golang:1.18.10-bullseye

# 拷贝 kaniko 内容
COPY --from=kaniko /kaniko /kaniko

# 拷贝上阶段编译后的文件
COPY --from=builder /usr/local/src/build/grpc-kit-cli-* /go/bin/grpc-kit-cli
COPY --from=builder /go/bin/protoc /go/bin/protoc
COPY --from=builder /go/bin/protoc-gen-go /go/bin/protoc-gen-go
COPY --from=builder /go/bin/protoc-gen-go-grpc /go/bin/protoc-gen-go-grpc
COPY --from=builder /go/bin/protoc-gen-grpc-gateway /go/bin/protoc-gen-grpc-gateway
COPY --from=builder /go/bin/protoc-gen-openapiv2 /go/bin/protoc-gen-openapiv2
COPY --from=builder /usr/local/include/google /usr/local/include/google
COPY --from=builder /go/src/ /go/src/
COPY --from=builder /usr/bin/docker /usr/bin/docker

CMD ["/go/bin/grpc-kit-cli"]
