# gRPC Kit Cli

基于[gRPC](https://github.com/golang/protobuf)、[gogo](https://github.com/gogo/protobuf)、[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)实现

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

protoc选择3.9.X版本，示例：

```
# 可直接从这里下载对应的二进制
https://github.com/protocolbuffers/protobuf/releases

#unzip protoc-3.9.2-linux-x86_64.zip
#mv bin/protoc $GOPATH/bin/
#mv include/google /usr/local/include/
```

protoc-gen-gogo选择为v1.3.X版本，示例：

```
go get github.com/gogo/protobuf/protoc-gen-gogo@v1.3.0
```

protoc-gen-swagger、protoc-gen-grpc-gateway选择1.9.X版本，示例：
 
```
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.9.6
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.9.6
```

- 安装其他依赖的proto定义文件

```
mkdir -p $GOPATH/src/github.com/{grpc-kit,gogo,grpc,googleapis,grpc-ecosystem}
git clone --depth 1 https://github.com/grpc-kit/api.git $GOPATH/src/github.com/grpc-kit/api
git clone --depth 1 https://github.com/gogo/googleapis.git $GOPATH/src/github.com/gogo/googleapis
git clone --depth 1 https://github.com/gogo/protobuf.git $GOPATH/src/github.com/gogo/protobuf
git clone --depth 1 https://github.com/grpc/grpc-proto.git $GOPATH/src/github.com/grpc/grpc-proto
git clone --depth 1 https://github.com/googleapis/googleapis.git $GOPATH/src/github.com/googleapis/googleapis
git clone --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway
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
grpc-kit-cli new -p demo -s test

make run
```

访问测试

```
curl 'http://127.0.0.1:10080/demo' | python -m json.tool
```
