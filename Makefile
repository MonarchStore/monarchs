.PHONY:
REPO   = arturom
IMAGE  = monarchs
TAG   := $(shell git log -1 --pretty=format:"%h")

DOCKER_IMAGE=$(REPO)/$(IMAGE):$(TAG)
LATEST_IMAGE=$(REPO)/$(IMAGE):latest

DOCKERFILE_DIR=.

CHART_DIR=chart/monarchs

RELEASE_NAME ?= monarchs
RELEASE_NAMESPACE ?= monarchs
DOCKER_TAG ?= latest


.PHONY: build
build:
	go build .

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
