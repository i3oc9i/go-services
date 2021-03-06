SHELL := /bin/bash

KLUSTER_NAME  := kenobi

KLUSTER_URL  := k3d.kenobi.local
REGISTRY_URL := registry.kenobi.local:5000

ARCH       := amd64
VERSION    := $(shell head -1 VERSION)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ------------------------------------------------------------------- Notes
#
# Access metrics directly (4000) 
# go install github.com/divan/expvarmon@latest
# expvarmon -ports="kenobi.local:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

all: help

# ------------------------------------------------------------------- Build
build:
	cd app/service/sales; go build --ldflags "-X main.build=$(VERSION)"      
	cd app/tools/jlogfmt; go build 

build-images: build-image-sales

build-image-sales:
	docker build \
      -f zarf/docker/dockerfile.sales \
      -t sales-$(ARCH):$(VERSION) \
      --build-arg BUILD_REF=$(VERSION) \
      --build-arg BUILD_DATE=$(BUILD_DATE) \
	  .
	docker tag sales-$(ARCH):$(VERSION) $(REGISTRY_URL)/sales-$(ARCH):$(VERSION)
	docker push $(REGISTRY_URL)/sales-$(ARCH):$(VERSION)
	cd zarf/k8s/k3d/sales-system; kustomize edit set image sales-image=$(REGISTRY_URL)/sales-$(ARCH):$(VERSION)


build-apply:  build-images svc-apply
build-update: build-images svc-update

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

# ------------------------------------------------------------------- Services
svc-apply:
	kustomize build zarf/k8s/k3d/sales-system | kubectl apply -f -

svc-update:
	kubectl -n sales-system rollout restart deployment sales

svc-delete:
	kustomize build zarf/k8s/k3d/sales-system | kubectl delete -f -

svc-status:
	kubectl -n sales-system get pods -o wide --watch 

svc-logs:
	kubectl -n sales-system logs -l app=sales --all-containers=true -f --tail=100 | go run app/tools/jlogfmt/main.go

# ------------------------------------------------------------------- K3D
k3d-up:
	k3d cluster create $(KLUSTER_NAME) --config ./zarf/infra/k3d/config.yaml
	k3d kubeconfig get $(KLUSTER_NAME) > .kube-config

k3d-start:
	k3d cluster start $(KLUSTER_NAME) 

k3d-stop:
	k3d cluster stop $(KLUSTER_NAME) 

k3d-info:
	kubectl cluster-info
	kubectl get node -o wide
	kubectl get svc -A -o wide 

k3d-down:
	k3d cluster delete $(KLUSTER_NAME) 

# ------------------------------------------------------------------- Clean
clean:
	go clean ./...

clobber: clean deps-clean k3d-down
	docker system prune --force

# ------------------------------------------------------------------- Help
help:
	@echo "make k3d-up          - create the local k3d kluster and the local registry."
	@echo "make build           - build the images of all the services and push them in the local registry."
	@echo "make services-deploy - deploy all the services in the local k3d kluster."
	@echo "make services-update - rollout all the running services."
	@echo "make clobber         - clean and delete the local k3d kluster and the local registry."
