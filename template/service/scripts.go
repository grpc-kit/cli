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

func (t *templateService) fileDirectoryScripts() {
	t.files = append(t.files, &templateFile{
		name:  "scripts/env",
		parse: true,
		body: `
# Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.
#
# https://grpc-kit.com/docs/spec-api/key-terms/

# 工具版本：生成该模版所使用的 cli 版本
CLI_VERSION={{ .Global.ReleaseVersion }}

# 产品代码：同一产品使用相同代码，使用单个词
PRODUCT_CODE={{ .Global.ProductCode }}

# 应用短名：同一产品使用相同代码，使用单个词
SHORT_NAME={{ .Global.ShortName }}

# 接口版本：全局接口版本
API_VERSION={{ .Template.Service.APIVersion }}

# 应用名称：用于生成文件的命名
APPNAME=${PRODUCT_CODE}-${SHORT_NAME}-${API_VERSION}

# 服务代码：用于 grpc 服务之间调用
SERVICE_CODE=${SHORT_NAME}.${API_VERSION}.${PRODUCT_CODE}

# 镜像名称：用于构建生成的镜像名称
#CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}

# 镜像版本：用于构建生成的镜像版本
#DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/docker.sh",
		body: `
#!/bin/bash

source scripts/env

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/docker.sh build"
  echo "\t ./scripts/docker.sh push"
fi

# 生成的镜像地址
RELEASE_VERSION=$(cat VERSION)
if test -z "${DOCKER_IMAGE_VERSION}"; then
  if test -z "${BUILD_ID}"; then
    DOCKER_IMAGE_VERSION=latest
  else
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
  fi
fi
if test -z "${CI_REGISTRY_IMAGE}"; then
  if test -z "${CI_REGISTRY_HOSTNAME}"; then
    CI_REGISTRY_HOSTNAME="docker.io"
  fi
  if test -z "${CI_REGISTRY_NAMESPACE}"; then
    CI_REGISTRY_NAMESPACE=${PRODUCT_CODE}
  fi
  CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}
fi

function build() {
  docker build -t ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION} ./
  echo "Now you can upload image: "docker push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}""
}

function push() {
  docker push ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}
}

function run() {
  docker run -i -t --rm \
      -v $GOPATH/pkg:/go/pkg \
      -v $(pwd):/usr/local/src \
      -w /usr/local/src \
      --network host \
      ccr.ccs.tencentyun.com/grpc-kit/cli:${CLI_VERSION} \
      make run
}

$1
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/genproto.sh",
		parse: true,
		body: `
#!/bin/bash

source scripts/env

# 生成 *.pb.go 文件
protoc \
    -I ./ \
    -I /usr/local/include/ \
    -I "${GOPATH}"/src \
    -I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/ \
    -I "${GOPATH}"/src/github.com/googleapis/googleapis/ \
    --go_opt paths=source_relative \
    --go_out ./ \
    --go-grpc_opt paths=source_relative \
    --go-grpc_opt require_unimplemented_servers=false \
    --go-grpc_out ./ \
    ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/*.proto

# 生成 *.pb.gw.go 与 swagger 接口文档
protoc \
    -I ./ \
    -I /usr/local/include/ \
    -I "${GOPATH}"/src \
    -I "${GOPATH}"/src/github.com/googleapis/googleapis/ \
    -I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/ \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt grpc_api_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.gateway.yaml \
    --grpc-gateway_out ./ \
    --openapiv2_opt disable_default_errors=true \
    --openapiv2_opt disable_service_tags=true \
    --openapiv2_opt grpc_api_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.gateway.yaml \
    --openapiv2_opt openapi_configuration=./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml \
    --openapiv2_out ./ \
    ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.proto

# 移动生成的 microservice.swagger.json 文件
if test -f ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.swagger.json; then
  mv ./api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.swagger.json ./public/doc/openapi-spec/
fi
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/precheck.sh",
		body: `
#!/bin/bash

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
FROM _DOCKER_IMAGE_FROM_

WORKDIR /opt

COPY build/service /opt/service
COPY config/app-mini.yaml /opt/config/app.yaml

EXPOSE 10080/tcp
EXPOSE 10081/tcp

ENTRYPOINT [ "/opt/service" ]
CMD [ "--config", "/opt/config/app.yaml" ]
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/version.sh",
		parse: true,
		body: `
#!/bin/bash

source scripts/env

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/version.sh prefix"
  echo "\t ./scripts/version.sh release"
  echo "\t ./scripts/version.sh update"
  exit 0;
fi

function prefix() {
  TEMP=$(grep "version: \".*\"" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml)
  PREFIX_VERSION=$(echo -n $TEMP | awk -F"\"" '{ print $2 }')
  echo $PREFIX_VERSION
}

function release() {
  TEMP=$(cat VERSION)
  RELEASE_VERSION=$TEMP

  if test -z $RELEASE_VERSION; then
    RELEASE_VERSION=$(git describe --tags --dirty --always 2>/dev/null)
  fi

  echo $RELEASE_VERSION
}

function update() {
  GOHOSTOS=$(go env GOHOSTOS)

  PREFIX_VERSION=$(prefix)
  RELEASE_VERSION=$(release)

  if test $PREFIX_VERSION == $RELEASE_VERSION; then
    return
  fi

  if test ${GOHOSTOS} = "darwin"; then
    sed -i "" "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml
  else
    # fix run in container
    # sed: couldn't open temporary file sed1DDoX9: Permission denied
    #sed -i "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml
    cp api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml /tmp/microservice.openapiv2.yaml
    sed -i "s#version: \"${PREFIX_VERSION}\"#version: \"${RELEASE_VERSION}\"#g" /tmp/microservice.openapiv2.yaml
    mv /tmp/microservice.openapiv2.yaml api/${PRODUCT_CODE}/${SHORT_NAME}/${API_VERSION}/microservice.openapiv2.yaml > /dev/null 2>&1
  fi
}

$1
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/manifests.sh",
		parse: true,
		body: `
#!/bin/bash

source scripts/env

# 需生成的模版类型
if test -z "${TEMPLATES}"; then
  TEMPLATES=dockerfile
fi

GOHOSTOS=$(go env GOHOSTOS)
KUBERNETES_NAMESPACE=default

# 生成的镜像地址
RELEASE_VERSION=$(cat VERSION)
if test -z "${DOCKER_IMAGE_VERSION}"; then
  if test -z "${BUILD_ID}"; then
    DOCKER_IMAGE_VERSION=latest
  else
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
  fi
fi
if test -z "${CI_REGISTRY_IMAGE}"; then
  if test -z "${CI_REGISTRY_HOSTNAME}"; then
    CI_REGISTRY_HOSTNAME="docker.io"
  fi
  if test -z "${CI_REGISTRY_NAMESPACE}"; then
    CI_REGISTRY_NAMESPACE=${PRODUCT_CODE}
  fi
  CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}
fi

# 解决 linux 与 darwin 的 sed 存在的差异
if test ${GOHOSTOS} = "darwin"; then
  SED="sed -i ''"
else
  SED="sed -i"
fi

# 编译与部署环境
if test -z "${BUILD_ENV}"; then
  BUILD_ENV=local
fi
if test -z "${DEPLOY_ENV}"; then
  DEPLOY_ENV=dev
fi
if test -n "${CI_BIZ_GROUP_APPID}"; then
  KUBERNETES_NAMESPACE=biz-${DEPLOY_ENV}-${CI_BIZ_GROUP_APPID}
fi

function clean() {
  rm -rf Dockerfile
  rm -rf deploy/systemd/
  rm -rf deploy/supervisor/
  rm -rf deploy/kubernetes/${DEPLOY_ENV}
}

function systemd() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/systemd
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
  fi

  cp -rf scripts/templates/systemd/microservice.service ${TEMPLATE_PATH}/${APPNAME}.service

  eval "$SED" "s#_SERVICE_CODE_#${SERVICE_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_PRODUCT_CODE_#${PRODUCT_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_SHORT_NAME_#${SHORT_NAME}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_API_VERSION_#${API_VERSION}#g" ${TEMPLATE_PATH}/${APPNAME}.service
  eval "$SED" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${APPNAME}.service
}

function supervisor() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/supervisor
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
  fi

  cp -rf scripts/templates/supervisor/microservice.conf ${TEMPLATE_PATH}/${APPNAME}.conf

  eval "$SED" "s#_SERVICE_CODE_#${SERVICE_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_PRODUCT_CODE_#${PRODUCT_CODE}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_SHORT_NAME_#${SHORT_NAME}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_API_VERSION_#${API_VERSION}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
  eval "$SED" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${APPNAME}.conf
}

function kubernetes() {
  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)/deploy/kubernetes/${DEPLOY_ENV}
  fi

  if test ! -d "${TEMPLATE_PATH}"; then
    mkdir -p ${TEMPLATE_PATH}
    mkdir -p ${TEMPLATE_PATH}/config/configmap
  fi

  cp -rf scripts/templates/kubernetes/* ${TEMPLATE_PATH}
  mv ${TEMPLATE_PATH}/workloads/deployments/microservice.yaml ${TEMPLATE_PATH}/workloads/deployments/${APPNAME}.yaml

  if test -f config/app-${DEPLOY_ENV}-${BUILD_ENV}.yaml; then
    cp -a config/app-${DEPLOY_ENV}-${BUILD_ENV}.yaml ${TEMPLATE_PATH}/config/configmap/app.yaml
  fi

  for FILE in kustomization.yaml service/ingresses.yaml service/services.yaml workloads/deployments/${APPNAME}.yaml
  do
    eval "${SED}" "s#DEPLOY_ENV#${DEPLOY_ENV}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_APPNAME_#${APPNAME}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_KUBERNETES_NAMESPACE_#${KUBERNETES_NAMESPACE}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_CI_REGISTRY_IMAGE_#${CI_REGISTRY_IMAGE}#g" ${TEMPLATE_PATH}/${FILE}
    eval "${SED}" "s#_DOCKER_IMAGE_VERSION_#${DOCKER_IMAGE_VERSION}#g" ${TEMPLATE_PATH}/${FILE}
  done
}

function dockerfile() {
  # 如未设置父镜像，默认为 scratch
  if test -z ${DOCKER_IMAGE_FROM}; then
    DOCKER_IMAGE_FROM=scratch
  fi

  # 如未设置目标地址，默认为当前路径
  if test -z "${TEMPLATE_PATH}"; then
    TEMPLATE_PATH=$(pwd)
  fi

  cp scripts/templates/Dockerfile ${TEMPLATE_PATH}

  if test ${GOHOSTOS} = "darwin"; then
    sed -i "" "s#_DOCKER_IMAGE_FROM_#${DOCKER_IMAGE_FROM}#g" ${TEMPLATE_PATH}/Dockerfile
  else
    sed -i "s#_DOCKER_IMAGE_FROM_#${DOCKER_IMAGE_FROM}#g" ${TEMPLATE_PATH}/Dockerfile
  fi
}

# 不做清理
#clean

# 避免运行无意义的指令
if test "${TEMPLATES}" = "dockerfile"; then
  dockerfile
elif test "${TEMPLATES}" = "kubernetes"; then
  kubernetes
elif test "${TEMPLATES}" = "systemd"; then
  systemd
elif test "${TEMPLATES}" = "supervisor"; then
  supervisor
fi
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/systemd/microservice.service",
		parse: false,
		body: `
[Unit]
After=network-online.target
Documentation=http://(app.yaml:services.http_address)/openapi-spec/
Description=The _SERVICE_CODE_ microservice. For more API detailed, please refer to the docs

[Service]
Type=simple
User=nobody
Restart=always
RestartSec=15s
TimeoutSec=60s
LimitNOFILE=65535
KillMode=control-group
MemoryLimit=2048M
ExecStart=/usr/local/_PRODUCT_CODE_/_SHORT_NAME_/_API_VERSION_/service --config /usr/local/_PRODUCT_CODE_/_SHORT_NAME_/_API_VERSION_/config/app.yaml

[Install]
Alias=_APPNAME_.service
WantedBy=multi-user.target
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/supervisor/microservice.conf",
		parse: false,
		body: `
[program:_APPNAME_]
command=/usr/local/_PRODUCT_CODE_/_SHORT_NAME_/_API_VERSION_/service --config /usr/local/_PRODUCT_CODE_/_SHORT_NAME_/_API_VERSION_/config/app.yaml
directory=/usr/local/_PRODUCT_CODE_/_SHORT_NAME_/_API_VERSION_/
autostart=true
autorestart=true
startsecs=10
startretries=3
stdout_logfile=/var/log/supervisor/%(program_name)s.log
stderr_logfile=/var/log/supervisor/%(program_name)s.log
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/kubernetes/kustomization.yaml",
		parse: true,
		body: `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: _KUBERNETES_NAMESPACE_

commonLabels:
  {{ .Global.APIEndpoint }}/appname: _APPNAME_
  {{ .Global.APIEndpoint }}/pm2-uuid: 3264e3fe-2bce-4835-8588-99651a8ddd3b

commonAnnotations:
  {{ .Global.APIEndpoint }}/pm2-uuid: 3264e3fe-2bce-4835-8588-99651a8ddd3b

configMapGenerator:
- name: _APPNAME_
  files:
  - app.yaml=config/configmap/app.yaml
  options:
    disableNameSuffixHash: true

replicas:
- name: _APPNAME_
  count: 1

resources:
- service/ingresses.yaml
- service/services.yaml
- workloads/deployments/_APPNAME_.yaml

images:
- name: _CI_REGISTRY_IMAGE_
  newTag: _DOCKER_IMAGE_VERSION_
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/kubernetes/service/ingresses.yaml",
		parse: true,
		body: `
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  #annotations:
  #  nginx.ingress.kubernetes.io/proxy-body-size: 10m
  name: _APPNAME_
spec:
  ingressClassName: nginx
  rules:
  - host: _APPNAME_._KUBERNETES_NAMESPACE_.{{ .Global.APIEndpoint }}
    http:
      paths:
      - path: /
        backend:
          service:
            name: _APPNAME_
            port:
              number: 10080
        pathType: Prefix
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/kubernetes/service/services.yaml",
		parse: true,
		body: `
---
apiVersion: v1
kind: Service
metadata:
  name: _APPNAME_
spec:
  ports:
  - name: http
    port: 10080
    protocol: TCP
    targetPort: 10080
  - name: grpc
    port: 10081
    protocol: TCP
    targetPort: 10081
  sessionAffinity: None
  type: ClusterIP
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "scripts/templates/kubernetes/workloads/deployments/microservice.yaml",
		parse: true,
		body: `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: _APPNAME_
spec:
  revisionHistoryLimit: 3
  strategy:
    rollingUpdate:
      maxSurge: 100%
      maxUnavailable: 0%
    type: RollingUpdate
  template:
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchLabels:
                {{ .Global.APIEndpoint }}/appname: _APPNAME_
            topologyKey: kubernetes.io/hostname
      containers:
      - args:
        - /opt/service
        - --config
        - /opt/config/app.yaml
        env:
        - name: GRPC_KIT_PUHLIC_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        image: _CI_REGISTRY_IMAGE_:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 5
          initialDelaySeconds: 30
          periodSeconds: 30
          successThreshold: 1
          tcpSocket:
            port: 10081
          timeoutSeconds: 5
        name: _APPNAME_
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz?service=_SERVICE_CODE_
            port: 10080
            scheme: HTTP
          initialDelaySeconds: 15
          periodSeconds: 15
          successThreshold: 1
          timeoutSeconds: 5
        resources:
          limits:
            cpu: "4000m"
            memory: 4048Mi
          requests:
            cpu: 100m
            memory: 100Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /opt/config
          name: config-volume
        - mountPath: /opt/logs/applog
          name: applog-volume
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - NET_ADMIN
            - SYS_ADMIN
            - NET_RAW
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext:
        runAsGroup: 65534
        runAsUser: 65534
        runAsNonRoot: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: config-volume
        configMap:
          defaultMode: 420
          name: _APPNAME_
      - name: applog-volume
        emptyDir: {}
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/binaries.sh",
		body: `
#!/bin/bash

source scripts/env

if test -z $1; then
  echo "Usage:"
  echo "\t ./scripts/binaries.sh protoc-gen-go"
  exit 0;
fi

# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-aarch_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-osx-aarch_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-osx-x86_64.zip
function protoc() {
  GOHOSTOS=$(go env GOHOSTOS)
  GOARCH=$(go env GOARCH)
  GOPATH=$(go env GOPATH)

  TARGET_OS=$GOHOSTOS
  TARGET_ARCH=$GOARCH

  if test "$GOHOSTOS" == "darwin"; then
    TARGET_OS="osx"
  fi

  if test "$GOARCH" == "arm64"; then
    TARGET_ARCH="aarch_64"
  elif test "$GOARCH" == "amd64"; then
    TARGET_ARCH="x86_64"
  fi

  echo "https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip"

  cd /tmp
  curl -L -O "https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip"

  unzip protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip
  mv bin/protoc "$GOPATH/bin/"
  mv include/google /usr/local/include/

  rm -f protoc-21.12-"$TARGET_OS"-"$TARGET_ARCH".zip
  rm -f readme.txt
  rmdir bin/
  rmdir include/
}

function protoc-gen-go-grpc() {
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
}

function protoc-gen-go() {
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
}

function protoc-gen-grpc-gateway() {
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
}

function protoc-gen-openapiv2() {
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
}

FILE=$(which $1)
if test -f "$FILE"; then
  echo "the binary already exists at: "$FILE""
else
  $1
  echo "download complete, this will place binaries in your \$GOBIN, make sure that your \$GOBIN is in your \$PATH."
fi
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/kaniko.sh",
		body: `
#!/bin/sh

source scripts/env

# TODO; kaniko 镜像仅支持 /bin/sh 解析器

# 生成的镜像地址
RELEASE_VERSION=$(cat VERSION)
if test -z "${DOCKER_IMAGE_VERSION}"; then
  if test -z "${BUILD_ID}"; then
    DOCKER_IMAGE_VERSION=latest
  else
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
  fi
fi
if test -z "${CI_REGISTRY_IMAGE}"; then
  if test -z "${CI_REGISTRY_HOSTNAME}"; then
    CI_REGISTRY_HOSTNAME="docker.io"
  fi
  if test -z "${CI_REGISTRY_NAMESPACE}"; then
    CI_REGISTRY_NAMESPACE=${PRODUCT_CODE}
  fi
  CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}
fi

function build() {
  /kaniko/executor --dockerfile ./Dockerfile --context ./ --destination ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION} --log-format text --log-timestamp
}

build
`,
	})

	t.files = append(t.files, &templateFile{
		name: "scripts/kcli.sh",
		body: `
#!/bin/bash

source scripts/env

# 生成的镜像地址
RELEASE_VERSION=$(cat VERSION)
if test -z "${DOCKER_IMAGE_VERSION}"; then
  if test -z "${BUILD_ID}"; then
    DOCKER_IMAGE_VERSION=latest
  else
    DOCKER_IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
  fi
fi
if test -z "${CI_REGISTRY_IMAGE}"; then
  if test -z "${CI_REGISTRY_HOSTNAME}"; then
    CI_REGISTRY_HOSTNAME="docker.io"
  fi
  if test -z "${CI_REGISTRY_NAMESPACE}"; then
    CI_REGISTRY_NAMESPACE=${PRODUCT_CODE}
  fi
  CI_REGISTRY_IMAGE=${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${APPNAME}
fi

# 当前目录必须为 gitops 仓库根路径
function git-commit() {
  if test "${CI_BIZ_CODE_BUILD}" != "true"; then
    return
  fi

  # fix: https://github.blog/2022-04-12-git-security-vulnerability-announced/
  git config --global --add safe.directory $(pwd)

  cd ${KUBERNETES_YAML_DIR}
  /bin/kustomize edit set image ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}

  git config user.name ${BUILD_USER}
  git config user.email ${BUILD_USER_EMAIL}

  git add kustomization.yaml
  git commit -m "build: set image ${CI_REGISTRY_IMAGE}:${DOCKER_IMAGE_VERSION}"

  git remote add gitops ${CI_OPS_REPO_URL}
  git push gitops HEAD:refs/heads/main
}

# 当前目录必须为 gitops 仓库根路径
function kubectl-apply() {
  cd ${KUBERNETES_YAML_DIR}
  /bin/kustomize build | /bin/kubectl apply -f -
}

$1
`,
	})
}
