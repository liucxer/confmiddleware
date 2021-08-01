# scaffold

来源 https://github.com/golang-standards/project-layout


```
github.com/liucxer/<group>/<project_name>/
    cmd/
        <project_name>/
            config/
            deploy/
            main.go
        <project_name>-<sub-app>/ # 将支持多入口项目集合，mono repo，不用开多个项目或者分支 
            config/
            deploy/
            main.go
    pkg/ # 模块，禁止使用 global，所有配置的单例，从 context.Context 中取
        apis/
            xxx/
            yyy/
            root.go
        ...
    version/
        version.go # 方便注入 版本信息
    .husky.yaml # 各种 githook
    .version # husky version 生成
    Makefile
```

## .husky.yaml

```yaml
hooks:
  pre-commit:
    - golangci-lint run
    - husky lint-staged
  commit-msg:
    - husky lint-commit

lint-staged:
  "*.go":
    - gofmt -l -w

lint-commit:
  email: "^(.+@rockontrol.com)$"
```

## Makefile

> 注意: 如果项目上使用 `opencv` 则 必须 `CGO_ENABLED=1` 。

```Makefile
PKG=$(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION=v$(shell cat .version)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X ${PKG}/version.Version=${VERSION}"
OPENAPI=tools openapi

MAIN_ROOT ?= ./cmd/srv-octohelm

up: dockerize
	go run $(MAIN_ROOT)

dockerize:
	go run $(MAIN_ROOT) dockerize

migrate:
	go run $(MAIN_ROOT) migrate

build: openapi
	cd $(MAIN_ROOT) && $(GOBUILD)

openapi: tools.install
	cd $(MAIN_ROOT) && $(OPENAPI)

tools.install:
	go install github.com/liucxer/cmd/cmd/tools

release:
	git push
	git push origin ${VERSION}
```

## Dockerfile.default

只是个示例文件，如无需要特殊修改的，手动重名为 `Dockerfile` 

```Dockerfile
# syntax = hub-dev.rockontrol.com/docker.io/docker/dockerfile:experimental

# syntax = hub-dev.rockontrol.com/docker.io/docker/dockerfile:experimental

# build-env 单独抽为一个 stage，通过 build --target=build-env 单独发布项目的开发环境 
ARG DOCKER_REGISTRY=hub-dev.rockontrol.com
FROM ${DOCKER_REGISTRY}/docker.io/library/golang:1.15.3-buster AS build-env

# setup private pkg 
ARG GITLAB_CI_TOKEN
ARG GITLAB_HOST=github.com/liucxer
ARG GOPROXY=https://goproxy.cn,direct
ENV GONOSUMDB=${GITLAB_HOST}/*
RUN git config --global url.https://gitlab-ci-token:${GITLAB_CI_TOKEN}@${GITLAB_HOST}/.insteadOf https://${GITLAB_HOST}/


# CGO 依赖安装
# RUN set -eux; \
#    \
#    apt-get update; \
#    apt-get install libxxx -y; \
#    \
#    rm -rf /var/lib/apt/lists/*

FROM build-env

WORKDIR /go/src
COPY ./ ./

# build
# 其他生成工具 tools 等，配置到 Makefile 中
ARG COMMIT_SHA
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod \
    make build

# runtime
FROM ${DOCKER_REGISTRY}/ghcr.io/querycap/distroless/static-debian10:latest
# CGO 使用
# FROM ${DOCKER_REGISTRY}/ghcr.io/querycap/distroless/cc-debian10:latest
# 并复制动态库
# COPY --from=builder /etc/ld.so.* /etc/
# COPY --from=builder /lib /lib
# COPY --from=builder /usr/lib /usr/lib
# COPY --from=builder /usr/local/lib /usr/local/lib

COPY --from=builder /go/src/cmd/example/example /go/bin/example

COPY --from=builder /go/src/cmd/example/openapi.json /go/bin/openapi.json
EXPOSE 80

ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/example"]
```

## .gitlab-ci.yml

1. 引用 golang ci 文件

```yaml
include:
  - project: 'infrav2/hx'
    file: "/ci/golang.gitlab-ci.yml"
```

2. 如果是 `opencv golang` 项目则引用对应的 ci 文件， 如下

```yaml
include:
  - project: 'infrav2/hx'
    file: "/ci/golang-opencv.gitlab-ci.yml"
```
