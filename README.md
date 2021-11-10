# gRPC Kit Cli

主要基于以下几个核心类库实现：

- [grpc](https://github.com/golang/protobuf)
- [gogo](https://github.com/gogo/protobuf)
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

# Overview

快速生成微服务模版，为同一产品提供统一的治理方式，提高多人协作效率

#  Prepare

- go版本必须大于等于1.13.x

版本检查

```
go version
```

centos 7 安装方式

```
yum install -y epel-release.noarch
yum install -y golang.x86_64
```

macOS 安装方式

```
brew install go
```

- 设置全局GOPATH且开启go mod支持

```
export GO111MODULE=on

# 根据实际情况是否需要设置proxy
export GOPROXY="https://goproxy.cn"

# GOPATH仅做示例，根据实际情况更改
export GOPATH=$HOME/go

export PATH=$PATH:$HOME/bin:$GOPATH/bin
```

请确保以上变量在系统全局生效，可以写入 $HOME/.bash_profile 或 $HOME/.zshrc 等

- 安装protoc与protoc-gen-*

protoc选择3.17.3版本，示例：

```
# 可直接从这里下载对应的二进制
cd /usr/local/src
curl -L -O 'https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip'
unzip protoc-3.17.3-linux-x86_64.zip
mv bin/protoc /usr/local/bin/
mv include/google /usr/local/include/
rmdir bin/ include/
```

protoc-gen-go、protoc-gen-go-grpc选择1.27.1与1.1.0版本

```
go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
```

protoc-gen-gogo选择为v1.3.2版本，示例：

```
go get github.com/gogo/protobuf/protoc-gen-gogo@v1.3.2
```

protoc-gen-swagger、protoc-gen-grpc-gateway选择1.16.0版本（不支持2.X），示例：
 
```
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
```

- 安装其他依赖的proto定义文件

```
mkdir -p $GOPATH/src/github.com/{grpc-kit,gogo,grpc,googleapis,grpc-ecosystem}
git clone --depth 1 https://github.com/grpc-kit/api.git $GOPATH/src/github.com/grpc-kit/api
git clone -b v1.1.0 --depth 1 https://github.com/gogo/googleapis.git $GOPATH/src/github.com/gogo/googleapis
git clone -b v1.3.2 --depth 1 https://github.com/gogo/protobuf.git $GOPATH/src/github.com/gogo/protobuf
git clone --depth 1 https://github.com/grpc/grpc-proto.git $GOPATH/src/github.com/grpc/grpc-proto
git clone --depth 1 https://github.com/googleapis/googleapis.git $GOPATH/src/github.com/googleapis/googleapis
git clone -b v1.16.0 --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway
```

# Install

- 下载二进制安装

```
https://github.com/grpc-kit/cli/releases
```

- 源码编译安装

```
git clone https://github.com/grpc-kit/cli.git

make build
cp ./build/grpc-kit-cli /usr/local/bin
```

# Getting Started

```
grpc-kit-cli new -p example -s test1

make run
```

访问测试

```
curl -u user1:pass1 http://127.0.0.1:10080/demo
```
