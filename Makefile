GO_VERSION ?= 1.18
APP := two-party-eddsa
TAG := $(shell git rev-parse --short HEAD)
AWS_ACCESS_KEY_ID := ${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY := ${AWS_SECRET_ACCESS_KEY}

# Default value
AWS_REGION ?= us-west-2
DOCKER_REGISTRY ?= ${DOCKER_REGISTRY}

.PHONY: build clean tool lint help

tool:
	go vet ./...; true
	gofmt -w .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"

# CI Pipeline
all: docker-login build push

docker-login:
	docker run --env AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) --env AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) --env AWS_REGION=$(AWS_REGION) --rm -i amazon/aws-cli ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(DOCKER_REGISTRY)

build :
	docker build --pull --build-arg GO_VERSION=$(GO_VERSION) -t ${DOCKER_REGISTRY}/${APP}:${TAG} -t ${DOCKER_REGISTRY}/${APP}:latest .

push :
	docker push ${DOCKER_REGISTRY}/${APP}:${TAG}
	docker push ${DOCKER_REGISTRY}/${APP}:latest