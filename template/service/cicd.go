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

func (t *templateService) fileDirectoryCICD() {
	t.files = append(t.files, &templateFile{
		name:  ".gitlab/workflows/grpc-kit.yml",
		parse: true,
		body: `
default:
  # TODO；根据具体情况选择运行的 runner 标签
  tags:
    - grpc-kit
  # TODO; 依赖文件注意使用缓存，避免每次下载
  #cache:
  #  paths:
  #    - /go/pkg/mod/
  # 框架使用的构建镜像
  image: ccr.ccs.tencentyun.com/grpc-kit/cli:{{ .Global.ReleaseVersion }}

# 默认全局变量
variables:
  CGO_ENABLED: "0"
  GIT_SSL_NO_VERIFY: "true"
  GO111MODULE: "on"
  GOPROXY: "https://goproxy.cn"
  GOSUMDB: "sum.golang.google.cn"
  #GOPRIVATE: "https://example.com"

# 流水线各阶段
stages:
  - pre
  - test
  - build
  - deploy
  - production

# 代码风格、格式检测
go-lint:
  stage: pre
  script:
    - make lint

# 依赖的相关依赖的组件
check-dep:
  stage: pre
  script:
    - which go
    - which protoc
    - which protoc-gen-go
    - which protoc-gen-go-grpc
    - which protoc-gen-grpc-gateway
    - which protoc-gen-openapiv2

# 业务单元测试
unit-tests:
  stage: test
  needs:
    - go-lint
    - check-dep
  script:
    - make test

# 代码覆盖率
coverage:
  stage: test
  script:
    - go test ./... -coverprofile=coverage.txt -covermode count
    - cat coverage.txt

# 生成发送测试报告
reports:
  stage: test
  needs:
    - unit-tests
    - coverage
  script:
    - echo "pass"

# 编译二进制文件
binary-local:
  stage: build
  needs:
    - reports
  script:
    - make build
  artifacts:
    paths:
      - build/
    expire_in: 24h
  when: manual
  allow_failure: false

# 发布容器至默认镜像中心
container-registry:
  stage: build
  needs:
    - binary-local
  script:
    - source scripts/env
    - export VERSION=$(cat VERSION)
    - echo ${CI_REGISTRY_PASSWORD} | docker login ${CI_REGISTRY} -u ${CI_REGISTRY_USER} --password-stdin
    - /kaniko/executor --dockerfile ${CI_PROJECT_DIR}/Dockerfile --context ${CI_PROJECT_DIR} --destination ${CI_REGISTRY_IMAGE}:${VERSION}

# 打成各种安装包，如：tar、rpm、deb
release-package:
  stage: build
  needs:
    - binary-local
  script:
    - echo "package tar"
    - echo "package rpm"
    - echo "package deb"
  artifacts:
    paths:
      - build/
    expire_in: 24h

# 部署测试环境
env-test:
  stage: deploy
  needs:
    - release-package
    - container-registry
  script:
    - echo "deploy test"

# 部署准线上环境
env-staging:
  stage: production
  needs:
    - env-test
  script:
    - echo "deploy staging"

# 部署正式环境，手工确认
env-prod:
  stage: production
  needs:
    - env-staging
  script:
    - echo "deploy production"
  only:
    - main
  when: manual
  allow_failure: false
`,
	})

	t.files = append(t.files, &templateFile{
		name:  ".jenkins/workflows/Jenkinsfile",
		parse: true,
		body: `
pipeline {
  agent {
    kubernetes {
      // TODO；目标集群，由系统管理员确认
      cloud 'dev'
      inheritFrom 'grpc'
      defaultContainer 'build'
    }
  }

  parameters {
    booleanParam(name: 'CI_BIZ_CODE_BUILD', defaultValue: true, description: '是否构建镜像，取消则直接至 k8s yaml 更新')
    booleanParam(name: 'CI_PIPELINE_SILENCE', defaultValue: false, description: '执行流水线全程静默无需二次确认')
    // TODO;
    choice(name: 'CI_REGISTRY_HOSTNAME', choices: ['ccr.ccs.tencentyun.com'], description: '支持的镜像中心列表')
    // TODO;
    choice(name: 'CI_REGISTRY_NAMESPACE', choices: ['opsaid'], description: '支持的镜像中心列表')
    // TODO;
    string(name: 'CI_BIZ_REPO_URL', defaultValue: 'https://{{ .Global.Repository }}.git', description: '业务源代码仓库')
    // TODO;
    string(name: 'CI_OPS_REPO_URL', defaultValue: 'https://{{ .Global.Repository }}.git', description: '运维源代码仓库')
    // TODO;
    string(name: 'CI_BIZ_BRANCH_NAME', defaultValue: 'main', description: '业务源代码仓库拉取的分支')
    // TODO;
    credentials(name: 'CI_BIZ_REPO_AUTH', defaultValue: '2a60ed63-1f38-4b18-a820-60cce23aa32e', credentialType: 'com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl', required: false)
    // TODO;
    credentials(name: 'CI_OPS_REPO_AUTH', defaultValue: '2a60ed63-1f38-4b18-a820-60cce23aa32e', credentialType: 'com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl', required: false)
  }

  environment {
    GOPROXY = "https://goproxy.cn"
  }

  options {
    disableConcurrentBuilds(abortPrevious: true)
    disableResume()
    timeout(time: 1, unit: 'HOURS')
  }

  stages {
    stage('Prepare') {
      when {
        allOf {
          not {
            environment name: 'CI_BIZ_REPO_URL', value: ''
          }
          not {
            environment name: 'CI_BIZ_BRANCH_NAME', value: ''
          }
          environment name: 'CI_BIZ_CODE_BUILD', value: 'true'
        }
      }

      steps {
        checkout scmGit(
            branches: [
                [name: CI_BIZ_BRANCH_NAME]
            ],
            extensions: [
                [$class: 'RelativeTargetDirectory', relativeTargetDir: 'source']
            ],
            userRemoteConfigs: [
                [
                    credentialsId: CI_BIZ_REPO_AUTH,
                    url: CI_BIZ_REPO_URL
                ]
            ]
        )

        // 执行代码检查
        sh '''
           cd source
           make protoc
           make lint
        '''
      }
    }

    stage('Test') {
      when {
        allOf {
          not {
            environment name: 'CI_BIZ_REPO_URL', value: ''
          }
          not {
            environment name: 'CI_BIZ_BRANCH_NAME', value: ''
          }
          environment name: 'CI_BIZ_CODE_BUILD', value: 'true'
        }
      }

      steps {
        // 执行单元测试等
        sh '''
           cd source
           make test
        '''
      }
    }

    stage('Build') {
      when {
        allOf {
          not {
            environment name: 'CI_BIZ_REPO_URL', value: ''
          }
          not {
            environment name: 'CI_BIZ_BRANCH_NAME', value: ''
          }
          environment name: 'CI_BIZ_CODE_BUILD', value: 'true'
        }
      }

      steps {
        // 选择特定语言容器，执行代码编译
        container('build') {
          sh '''
             cd source
             make build
             make manifests TEMPLATES=dockerfile
             make manifests TEMPLATES=kubernetes TEMPLATE_PATH=../gitops/deploy/kubernetes/dev/
          '''
        }

        // 选择 kaniko 容器，执行构建镜像并上传
        container('kaniko') {
          sh '''
             cd source
             ./scripts/kaniko.sh
          '''
        }
      }
    }

    stage('Confirm') {
      when {
        environment name: 'CI_PIPELINE_SILENCE', value: 'false'
      }

      steps {
        container('kcli') {
          sh '''
            cd gitops/deploy/kubernetes/dev
            cat kustomization.yaml
          '''
        }

        // 人工审批确认，是否部署至环境
        input "Does look ok?"
      }
    }

    stage('Production') {
      steps {
        container('kcli') {

          withCredentials([gitUsernamePassword(credentialsId: CI_OPS_REPO_AUTH, gitToolName: 'git-tool')]) {
            wrap([$class: 'BuildUser']) {
              script {
                if (env.CI_BIZ_CODE_BUILD == "true" ) {
                  sh '''
                     cd source
                     RELEASE_VERSION=$(cat VERSION)
                     IMAGE_VERSION=${RELEASE_VERSION}-build.${BUILD_ID}
                     cd ../

                     cd gitops/deploy/kubernetes/dev
                     /bin/kustomize edit set image ${CI_REGISTRY_HOSTNAME}/${CI_REGISTRY_NAMESPACE}/${JOB_BASE_NAME}:${IMAGE_VERSION}

                     # fix: https://github.blog/2022-04-12-git-security-vulnerability-announced/
                     git config --global --add safe.directory ${WORKSPACE}/gitops

                     git config user.name ${BUILD_USER}
                     git config user.email ${BUILD_USER_EMAIL}

                     git add kustomization.yaml
                     git commit -m "build: set image ${JOB_BASE_NAME}:${IMAGE_VERSION}"

                     git remote add gitops ${CI_OPS_REPO_URL}
                     git push gitops HEAD:refs/heads/main
                  '''
                }
              }
            }

            sh '''
               cd gitops/deploy/kubernetes/dev
               /bin/kustomize build | /bin/kubectl apply -f -
            '''
          }
        }
      }
    }
  }

  post {
    always {
      // TODO；根据实际情况调用接口推送通知
      echo "Send notifications for result: ${currentBuild.result}"
    }
  }
}
`,
	})
}
