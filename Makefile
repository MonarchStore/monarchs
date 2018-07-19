.PHONY:
REPO=arturom
IMAGE=monarchs
TAG := $(shell git log -1 --pretty=format:"%h")

DOCKER_IMAGE=$(REPO)/$(IMAGE):$(TAG)
LATEST_IMAGE=$(REPO)/$(IMAGE):latest

DOCKERFILE_DIR=.


.PHONY: container
container:
	docker build -t $(DOCKER_IMAGE) $(DOCKERFILE_DIR) --build-arg commit=$(TAG)
	docker tag $(DOCKER_IMAGE) $(LATEST_IMAGE)
