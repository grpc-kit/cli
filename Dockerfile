FROM golang:1.17.13-buster as builder

# 设置环境变量
ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn" \
    PATH=$PATH:/go/bin

WORKDIR /usr/local/src

COPY . .
RUN go mod tidy && make build

FROM golang:1.17.13-buster

# 拷贝上阶段编译后的二进制
COPY --from=builder /usr/local/src/build/grpc-kit-cli /go/bin/grpc-kit-cli

# 系统使用非交互式模式
ARG DEBIAN_FRONTEND=noninteractive

# 更换镜像地址
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

# 设置时区
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && dpkg-reconfigure -f noninteractive tzdata

# 安装基础软件
RUN apt-get install -y unzip

# 下载 protoc 二进制
RUN cd /usr/local/src \
    && curl -L -O 'https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip' \
    && unzip protoc-3.17.3-linux-x86_64.zip \
    && mv bin/protoc /usr/local/bin/ \
    && mv include/google /usr/local/include/ \
    && rmdir bin/ \
    && rmdir include/

# 下载 protoc-gen 二进制
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 \
    && go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 \
    && go get github.com/gogo/protobuf/protoc-gen-gogo@v1.3.2 \
    && go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0 \
    && go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0

# 安装其他依赖的 proto 定义文件
RUN git clone --depth 1 https://github.com/grpc-kit/api.git $GOPATH/src/github.com/grpc-kit/api \
    && git clone -b v1.1.0 --depth 1 https://github.com/gogo/googleapis.git $GOPATH/src/github.com/gogo/googleapis \
    && git clone -b v1.3.2 --depth 1 https://github.com/gogo/protobuf.git $GOPATH/src/github.com/gogo/protobuf \
    && git clone --depth 1 https://github.com/grpc/grpc-proto.git $GOPATH/src/github.com/grpc/grpc-proto \
    && git clone --depth 1 https://github.com/googleapis/googleapis.git $GOPATH/src/github.com/googleapis/googleapis \
    && git clone -b v1.16.0 --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway

CMD ["/go/bin/grpc-kit-cli"]
