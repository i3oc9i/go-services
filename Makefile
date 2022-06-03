SHELL := /bin/bash

KLUSTER_NAME=local
REGISTRY=registry.local:5000

APP=go-service
ARCH=amd64

BUILD_VERSION := $(shell cat ./VERSION)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: help

# ------------------------------------------------------------------- Build
build-service:
	docker build \
    	-f zarf/docker/dockerfile \
    	-t $(APP)-$(ARCH):$(BUILD_VERSION) \
    	--build-arg BUILD_VERSION=$(BUILD_VERSION) \
    	--build-arg BUILD_DATE=$(BUILD_DATE) \
		.
	docker tag $(APP)-$(ARCH):$(BUILD_VERSION) $(REGISTRY)/$(APP)-$(ARCH):$(BUILD_VERSION)
	docker push $(REGISTRY)/$(APP)-$(ARCH):$(BUILD_VERSION)

# ------------------------------------------------------------------- Dependencies
deps-update:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-clean:
	go clean -modcache

# ------------------------------------------------------------------- Service
service-deploy:
	kubectl apply -f ./zarf/k8s/basic/go_service.yaml

service-update: build-service
	kubectl -n service-system rollout restart deployment go-service

service-delete:
	kubectl delete -f ./zarf/k8s/basic/go_service.yaml

# ------------------------------------------------------------------- K3D
k3d-create:
	k3d cluster create $(KLUSTER_NAME) --config ./zarf/k3d/config.yaml
	kubectl cluster-info
	kubectl get node -o wide

k3d-start:
	k3d cluster start $(KLUSTER_NAME) 

k3d-stop:
	k3d cluster stop $(KLUSTER_NAME) 

k3d-delete:
	k3d cluster delete $(KLUSTER_NAME) 

# ------------------------------------------------------------------- Clean
clean:
	go clean

clobber: clean k3d-delete

# ------------------------------------------------------------------- Help
help:
	@echo "make k3d-create     - create the local k3d kluster and the local registry."
	@echo "make build-service  - build the service image and push it in the local registry."
	@echo "make service-deploy - deploy the service in the local k3d kluster."
	@echo "make service-update - rollout the running service."
	@echo "make clobber        - make clean and delete the local k3d kluster and the local registry."