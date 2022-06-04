SHELL := /bin/bash

KLUSTER_ID  = kenobi
REGISTRY_ID = registry-kenobi.local:5000

SERVICE = go-service
ARCH    = amd64

BUILD_VERSION := $(shell cat ./VERSION)
BUILD_DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: help

# ------------------------------------------------------------------- Build
build:
	go build -o ${SERVICE} --ldflags "-X main.build=$(BUILD_VERSION)"

build-image:
	docker build \
    	-f zarf/docker/dockerfile \
    	-t $(SERVICE)-$(ARCH):$(BUILD_VERSION) \
		--build-arg SERVICE_NAME=$(SERVICE) \
    	--build-arg BUILD_VERSION=$(BUILD_VERSION) \
    	--build-arg BUILD_DATE=$(BUILD_DATE) \
		.
	docker tag $(SERVICE)-$(ARCH):$(BUILD_VERSION) $(REGISTRY_ID)/$(SERVICE)-$(ARCH):$(BUILD_VERSION)
	docker push $(REGISTRY_ID)/$(SERVICE)-$(ARCH):$(BUILD_VERSION)

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
	kustomize build ./zarf/k8s/k3d/go-service | kubectl apply -f -

service-update:
	kubectl -n service-system rollout restart deployment $(SERVICE)

service-delete:
	kustomize build ./zarf/k8s/k3d/go-service | kubectl delete -f -

# ------------------------------------------------------------------- K3D
k3d-setup:
	k3d cluster create $(KLUSTER_ID) --config ./zarf/infra/k3d/config.yaml
	k3d kubeconfig get $(KLUSTER_ID) > .kube-config
	kubectl cluster-info
	kubectl get node -o wide

k3d-start:
	k3d cluster start $(KLUSTER_ID) 

k3d-stop:
	k3d cluster stop $(KLUSTER_ID) 

k3d-destroy:
	k3d cluster delete $(KLUSTER_ID) 

# ------------------------------------------------------------------- Clean
clean:
	go clean

clobber: clean k3d-destroy
	go clean -modcache
	docker system prune --force

# ------------------------------------------------------------------- Help
help:
	@echo "make k3d-create     - create the local k3d kluster and the local registry."
	@echo "make build-image    - build the service image and push it in the local registry."
	@echo "make service-deploy - deploy the service in the local k3d kluster."
	@echo "make service-update - rollout the running service."
	@echo "make clobber        - make clean and delete the local k3d kluster and the local registry."
