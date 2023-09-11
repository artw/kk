APP := kk
BUILD_DIR := build
SOURCES := $(shell find . -name '*.go')
DOCKER_TAG=artw/kk:latest

.PHONY: all clean docker-build docker-push

all: $(BUILD_DIR)/$(APP) docker-build docker-push

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/$(APP): $(SOURCES) | $(BUILD_DIR)
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP)

docker-build:
	docker build . -t $(DOCKER_TAG)

docker-push:
	docker push $(DOCKER_TAG)

clean:
	rm -rf $(BUILD_DIR)
