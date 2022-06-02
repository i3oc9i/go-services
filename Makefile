SHELL := /bin/bash

APP=go-service
ARCH=amd64

KLUSTER_NAME=local

BUILD_VERSION := $(shell cat ./VERSION)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: docker-build

# ------------------------------------------------------------------- Build
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

# ------------------------------------------------------------------- Run

# ------------------------------------------------------------------- K3D
k3d-create:
	k3d cluster create $(KLUSTER_NAME) --config ./zarf/k3d/config.yaml

k3d-start:
	k3d cluster start $(KLUSTER_NAME) 

k3d-stop:
	k3d cluster stop $(KLUSTER_NAME) 

k3d-delete:
	k3d cluster delete $(KLUSTER_NAME) 

# ------------------------------------------------------------------- Clean
clean:
	go clean

