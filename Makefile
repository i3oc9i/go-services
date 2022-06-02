SHELL := /bin/bash

APP=go-service
ARCH=amd64

BUILD_VERSION := $(shell cat ./VERSION)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: docker-build

build:
	go build --ldflags "-X main.build=$(BUILD_VERSION)"

docker-build:
	docker build \
    	-f zarf/docker/dockerfile \
    	-t $(APP)-$(ARCH):"$(BUILD_VERSION)" \
    	--build-arg BUILD_VERSION="$(BUILD_VERSION)" \
    	--build-arg BUILD_DATE="$(BUILD_DATE)" \
		.

mod-update:
	go mod tidy
	go mod vendor

clean:
	go clean

