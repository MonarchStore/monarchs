.PHONY:
REPO=arturom
IMAGE=monarchs
# TAG := $(shell git log -1 --pretty=format:"%h")
TAG := $(shell git describe --tags --always --dirty)
DOCKER_IMAGE=$(REPO)/$(IMAGE):$(TAG)
LATEST_IMAGE=$(REPO)/$(IMAGE):latest

DOCKERFILE_DIR=.

CHART_DIR=chart/monarchs
BUILD_DIR=build
RELEASE_NAME ?= monarchs
RELEASE_NAMESPACE ?= monarchs
DOCKER_TAG ?= latest

.PHONY: test
test:
	go test -v -race $(shell go list ./... | grep -v /vendor/)

dep:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -vendor-only

.PHONY: build
SOURCES = $(shell find . -name '*.go')
BUILD_FLAGS = -v
VERSION = $(TAG)
LDFLAGS = -X github.com/MonarchStore/monarchs/config.Version=$(VERSION) -w -s
BINARY = monarchs

build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	@echo "Building $(VERSION)"
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)"

.PHONY: install
install:
	go install -v ./...

.PHONY: container
container:
	docker build -t $(DOCKER_IMAGE) $(DOCKERFILE_DIR) --build-arg commit=$(TAG)
	docker tag $(DOCKER_IMAGE) $(LATEST_IMAGE)

.PHONY: push-container
push-container: container
	docker push $(DOCKER_IMAGE)
	docker push $(LATEST_IMAGE)

.PHONY: chart
chart:
	helm lint $(CHART_DIR)

.PHONY: install-chart
install-chart:
	helm upgrade --install $(RELEASE_NAME) \
		--namespace $(RELEASE_NAMESPACE) \
		--set image.tag=$(DOCKER_TAG) \
		$(CHART_DIR)
clean:
	@rm -rf build
