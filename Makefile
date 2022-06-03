SHELL := /bin/bash

KLUSTER_NAME=local
REGISTRY=registry.localhost:5000

APP=go-service
ARCH=amd64

BUILD_VERSION := $(shell cat ./VERSION)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: docker-build

# ------------------------------------------------------------------- Build
build:
	go build --ldflags "-X main.build=$(BUILD_VERSION)"

docker-build:
	docker build \
    	-f zarf/docker/dockerfile \
    	-t $(APP)-$(ARCH):$(BUILD_VERSION) \
    	--build-arg BUILD_VERSION=$(BUILD_VERSION) \
    	--build-arg BUILD_DATE=$(BUILD_DATE) \
		.
	docker tag $(APP)-$(ARCH):$(BUILD_VERSION) $(REGISTRY)/$(APP)-$(ARCH):$(BUILD_VERSION)
	docker push $(REGISTRY)/$(APP)-$(ARCH):$(BUILD_VERSION)

mod-update:
	go mod tidy
	go mod vendor

# ------------------------------------------------------------------- Service
service-deploy:
	kubectl apply -f ./zarf/k8s/basic/go_service.yaml

service-delete:
	kubectl delete -f ./zarf/k8s/basic/go_service.yaml

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

