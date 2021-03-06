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

func (t *templateService) fileDirectoryScripts() {
	t.files = append(t.files, &templateFile{
		name:  "scripts/env",
		parse: true,
		body: `
# Code generated by grpc-kit-cli. DO NOT EDIT.

# 全局接口版本
API_VERSION={{ .Template.Service.APIVersion }}

# 生成该模版所使用的cli版本
CLI_VERSION={{ .Global.ReleaseVersion }}
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/docker.sh",
		body: `
#!/bin/sh

# 容器镜像版本号，去掉v开头
IMAGE_VERSION=${RELEASE_VERSION}
if test ${RELEASE_VERSION:0:1} = "v"; then
    IMAGE_VERSION=${RELEASE_VERSION:1}
fi

# 如未设置父镜像，默认为scratch
if test -z ${IMAGE_FROM}; then
    IMAGE_FROM=scratch
fi

# 生成的容器镜像地址
IMAGE_ADDR=${IMAGE_HOST}/${NAMESPACE}/${SHORTNAME}:${IMAGE_VERSION}

cp scripts/templates/Dockerfile ./

GOHOSTOS=$(go env GOHOSTOS)

if test ${GOHOSTOS} = "darwin"; then
    sed -i "" "s#{{IMAGE_FROM}}#${IMAGE_FROM}#g" Dockerfile
else
    sed -i "s#{{IMAGE_FROM}}#${IMAGE_FROM}#g" Dockerfile
fi

docker build -t ${IMAGE_ADDR} ./
#docker push ${IMAGE_ADDR}
echo "Now you can upload image: "docker push ${IMAGE_ADDR}""
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/genproto.sh",
		body: `
#!/bin/sh

source scripts/env

# 生成*.pb.go文件
protoc -I./api/proto/${API_VERSION}/ \
        -I${GOPATH}/src \
        -I${GOPATH}/src/github.com/gogo/googleapis/ \
        -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/ \
        -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --gogo_out=plugins=grpc,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
./api/proto/${API_VERSION}/ \
        api/proto/${API_VERSION}/*.proto

# 生成*.pb.gw.go与swagger接口文档
protoc -I./api/proto/${API_VERSION}/ \
        -I${GOPATH}/src \
        -I${GOPATH}/src/github.com/gogo/googleapis/ \
        -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/ \
        -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --grpc-gateway_out=allow_patch_feature=false,allow_repeated_fields_in_body=true,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:\
./api/proto/${API_VERSION}/ \
        --swagger_out=logtostderr=true,allow_repeated_fields_in_body=true:./api/doc/openapi-spec/ \
        api/proto/${API_VERSION}/microservice.proto
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/precheck.sh",
		body: `
#!/bin/sh

# 确保GOPATH变量有设置
if test -z "${GOPATH}"; then
  echo "Please set the environment variable GOPATH before running make"
  exit 1
fi
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/templates/Dockerfile",
		body: `
FROM {{IMAGE_FROM}}

WORKDIR /opt

COPY build/service /opt/service
COPY config/app-sample.toml /opt/config/app.toml

EXPOSE 10080/tcp
EXPOSE 10081/tcp

ENTRYPOINT [ "/opt/service" ]
CMD [ "--config", "/opt/config/app.toml" ]
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/version.sh",
		body: `
#!/bin/sh

source scripts/env

if test -z $1; then
    echo -e "Usage:"
    echo -e "\t ./scripts/version.sh prefix"
    echo -e "\t ./scripts/version.sh release"
    echo -e "\t ./scripts/version.sh update"
    exit 0;
fi

function prefix() {
    TEMP=$(grep "version: \".*\"" api/proto/${API_VERSION}/microservice.proto)
    PREFIX_VERSION=$(echo -n $TEMP | awk -F"\"" '{ print $2 }')
    echo $PREFIX_VERSION
}

function release() {
    TEMP=$(git describe --tags --dirty --always 2>/dev/null)
    RELEASE_VERSION=$TEMP

    if test -z $RELEASE_VERSION; then
        RELEASE_VERSION=v0.0.0
    fi

    echo $RELEASE_VERSION
}

function update() {
    GOHOSTS=$(go env GOHOSTOS)

    PREFIX_VERSION=$(prefix)
    RELEASE_VERSION=$(release)

    if test ${GOHOSTS} = "darwin"; then
        sed -i "" "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/proto/${API_VERSION}/microservice.proto
    else
        sed -i "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/proto/${API_VERSION}/microservice.proto
    fi
}

$1
`,
	})
}
