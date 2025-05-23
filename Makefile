override APP_NAME=todo_app
override GO_VERSION=1.24
override GOLANGCI_LINT_VERSION=v1.46.2
override SECUREGO_GOSEC_VERSION=2.12.0
override HADOLINT_VERSION=v2.10.0

GOOS?=$(shell go env GOOS || echo linux)
GOARCH?=$(shell go env GOARCH || echo amd64)
CGO_ENABLED?=0

DOCKER_REGISTRY=docker.io
DOCKER_IMAGE=todo_app
DOCKER_USER=
DOCKER_PASSWORD=
DOCKER_TAG?=latest

ifeq (, $(shell which docker))
$(error "Binary docker not found in $(PATH)")
endif

.PHONY: all
all: cleanup vendor lint test build

.PHONY: cleanup
cleanup:
	@rm ${PWD}/bin/${APP_NAME} || true
	@rm ${PWD}/coverage.out || true
	@rm -r ${PWD}/vendor || true

.PHONY: tidy
tidy:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go mod tidy

.PHONY: vendor
vendor:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go mod vendor

.PHONY: lint-golangci-lint
lint-golangci-lint:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golangci/golangci-lint:${GOLANGCI_LINT_VERSION} \
			golangci-lint run -v

.PHONY: lint-gosec
lint-gosec:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		securego/gosec:${SECUREGO_GOSEC_VERSION} \
			/project/...

.PHONY: lint-dockerfile
lint-dockerfile:
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		hadolint/hadolint:${HADOLINT_VERSION} \
			hadolint \
				/project/build/docker/cmd/todo_app/Dockerfile

.PHONY: lint
lint:
	@make lint-golangci-lint
	@make lint-gosec
	@make lint-dockerfile

.PHONY: test
test:
	@rm -r ${PWD}/coverage.out || true
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		golang:${GO_VERSION} \
			go test \
				-race \
				-mod vendor \
				-covermode=atomic \
				-coverprofile=/project/coverage.out \
					/project/...

.PHONY: build
build:
	@rm ${PWD}/bin/${APP_NAME} || true
	@docker run --rm \
		-v ${PWD}:/project \
		-w /project \
		-e GOOS=${GOOS} \
		-e GOARCH=${GOARCH} \
		-e CGO_ENABLED=${CGO_ENABLED} \
		-e GO111MODULE=on \
		golang:${GO_VERSION} \
			go build \
				-mod vendor \
				-o /project/bin/${APP_NAME} \
				-v /project/cmd/${APP_NAME}

.PHONY: docker-image-build
docker-image-build:
ifndef DOCKER_IMAGE
	$(error DOCKER_IMAGE is not set)
endif
ifndef DOCKER_TAG
	$(error DOCKER_TAG is not set)
endif
	@docker rmi ${DOCKER_IMAGE}:${DOCKER_TAG} || true
	@docker build \
		-f ${PWD}/build/docker/cmd/todo_app/Dockerfile \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
			.

.PHONY: docker-image-push
docker-image-push:
ifndef DOCKER_USER
	$(error DOCKER_USER is not set)
endif
ifndef DOCKER_PASSWORD
	$(error DOCKER_PASSWORD is not set)
endif
ifndef DOCKER_REGISTRY
	$(error DOCKER_REGISTRY is not set)
endif
ifndef DOCKER_IMAGE
	$(error DOCKER_IMAGE is not set)
endif
ifndef DOCKER_TAG
	$(error DOCKER_TAG is not set)
endif
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	@docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
