CURRENT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SRC = $(shell find $(CURRENT_DIR) -type f -name '*.go' -not -path "$(CURRENT_DIR)/vendor/*")

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse HEAD)
BUILD_TS := $(shell date +'%Y-%m-%d %H:%M:%S')

REGISTRY_ADDRESS ?= registry.sensetime.com
REGISTRY_GROUP ?= beacon
SHORT_HASH := $(shell git rev-parse --short=8 HEAD)
IMAGE_VERSION := $(BRANCH)-$(SHORT_HASH)
IMAGE_NAME := datacollector-service
IMAGE_FULLNAME := $(REGISTRY_ADDRESS)/$(REGISTRY_GROUP)/apcloud/app/$(IMAGE_NAME):$(IMAGE_VERSION)
ARCH ?= amd64

HELM_BUILD := $(shell date +'%Y%m%d%H%M%S')+$(BRANCH)
COMMON_CHART_DIR := $(CURRENT_DIR)/../../library/common-chart/deploy/helm/common
COMMON_CHART_REPO := file:///data/common

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	#go test $(CURRENT_DIR)/...

.PHONY: build
build: $(SRC)
	go mod tidy
	go build -ldflags '-X "main.GitBranch=$(BRANCH)" -X "main.GitHash=$(HASH)" -X "main.BuildTS=$(BUILD_TS)"' -o $(CURRENT_DIR)/bin/ $(CURRENT_DIR)/...

.PHONY: image
image:
	go mod tidy
	go mod vendor
	docker build \
		--build-arg PROJ_NAME="$(IMAGE_NAME)" \
		--build-arg GIT_BRANCH="$(BRANCH)" \
		--build-arg GIT_HASH="$(HASH)" \
		--build-arg BUILD_TS="$(BUILD_TS)" \
		--build-arg ARCHITECTURE="$(ARCH)" \
		-t $(IMAGE_FULLNAME) $(CURRENT_DIR)
	rm -rf vendor/

.PHONY: helm-clean
helm-clean:
	rm -rf $(CURRENT_DIR)/deploy/test

.PHONY: helm-gen
helm-gen:
	rm -rf $(CURRENT_DIR)/deploy/helm_generated
	cp -r $(CURRENT_DIR)/deploy/helm $(CURRENT_DIR)/deploy/helm_generated
	sed -i.bak -e "s|{{APP_VERSION}}|$(IMAGE_VERSION)|g" $(CURRENT_DIR)/deploy/helm_generated/*/values.yaml
	sed -i.bak -e "s|{{REGISTRY_GROUP}}|$(REGISTRY_GROUP)|g" $(CURRENT_DIR)/deploy/helm_generated/*/values.yaml
	sed -i.bak -e "s|{{APP_VERSION}}|$(IMAGE_VERSION)|g" $(CURRENT_DIR)/deploy/helm_generated/*/Chart.yaml
	sed -i.bak -e "s|{{BUILD_VERSION}}|$(HELM_BUILD)|g" $(CURRENT_DIR)/deploy/helm_generated/*/Chart.yaml
	sed -i.bak -e "s|{{COMMON_CHART}}|$(COMMON_CHART_REPO)|g" $(CURRENT_DIR)/deploy/helm_generated/*/Chart.yaml

.PHONY: helm-dep-build
helm-dep-build: helm-gen
	docker run -it --rm \
		-v $(CURRENT_DIR)/deploy/helm_generated:/data/helm \
		-v $(COMMON_CHART_DIR):/data/common \
		registry.sensetime.com/beacon/ci/deploy:2.0.0 \
		bash -c 'cd /data/helm; for chart in ./*; do helm dependency build $$chart; done;'

.PHONY: helm-install
helm-install: helm-clean helm-dep-build
	mkdir -p $(CURRENT_DIR)/deploy/test
	docker run -it --rm --network host \
		-v $(CURRENT_DIR)/deploy:/data/deploy \
		-v $(CURRENT_DIR)/deploy/test:/data/test \
		-v $(HOME)/.kube:/root/.kube \
		registry.sensetime.com/beacon/ci/deploy:2.0.0 \
		bash -c 'cd /data/deploy/helm_generated; for chart in *; do helm template --dry-run --debug $$chart ./$$chart > /data/test/$$chart-install.yml; done;'

.PHONY: debug
debug:
	bash debug/env.bash
	kubevpn connect