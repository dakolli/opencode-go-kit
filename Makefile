# Load environment variables from a custom file if specified, defaulting to .env
ENV_FILE ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

# Dynamically detect host user ID and group ID to prevent file permission issues in mounted volumes
USER_ID ?= $(shell id -u)
GROUP_ID ?= $(shell id -g)
export USER_ID
export GROUP_ID

# Default fallback values if not specified in the environment or .env
HOST_PORT ?= 4002
USERNAME ?= opencode
WORKSPACE_PASSWORD ?= pass1
VOLUME_PATH ?= .docker/volumes/workspace_example
IMAGE_NAME ?= opencode-go-kit:latest
PROJECT_NAME ?= opencode-alpha

.PHONY: build up down run stop generate generate-api update-readme

# Target to build the Docker image using pure Docker build
build:
	docker build \
		--build-arg USER_ID=$(USER_ID) \
		--build-arg GROUP_ID=$(GROUP_ID) \
		--build-arg USER_NAME=$(USERNAME) \
		--build-arg GROUP_NAME=$(USERNAME) \
		-t $(IMAGE_NAME) .

# Target to run docker-compose using the environment variables
up:
	mkdir -p $(VOLUME_PATH)
	docker compose -p $(PROJECT_NAME) up -d

# Target to stop the compose services
down:
	docker compose -p $(PROJECT_NAME) down

# Target to run a single container using standard docker run directly (without docker-compose)
run:
	mkdir -p $(VOLUME_PATH)
	docker run -d \
		--name $(PROJECT_NAME)-agent \
		-p 127.0.0.1:$(HOST_PORT):4096 \
		-v $(abspath $(VOLUME_PATH)):/workspace \
		-e OPENCODE_SERVER_PASSWORD=$(WORKSPACE_PASSWORD) \
		-e OPENCODE_KEY=$(OPENCODE_KEY) \
		-e OPENCODE_CONFIG=/workspace/opencode.jsonc \
		--restart unless-stopped \
		$(IMAGE_NAME)

# Target to stop and remove the direct docker run container
stop:
	docker stop $(PROJECT_NAME)-agent || true
	docker rm $(PROJECT_NAME)-agent || true

# Target to run go generate, injecting environment variables from .env
generate-client:
	go generate ./pkg/client

# Checks the current coverage of implemented API methods, and if partial
# coverage is detected, it sequentially regenerates the client from local
# OpenAPI schemas and generates the strongly-typed wrapper methods.
generate-api:
	@TEST_OUT=$$(GOWORK=off go test -v -count=1 ./pkg/api 2>&1); \
	echo "$$TEST_OUT"; \
	if echo "$$TEST_OUT" | grep -F -q "100.00%"; then \
		echo "Coverage is 100.00%. Updating README coverage list..."; \
		GOWORK=off go run tools/readme_updater/main.go; \
	else \
		echo "Partial coverage detected! Regenerating client first, then wrappers..."; \
		if ! curl -s -f -u $(USERNAME):$(PASSWORD) http://localhost:$(HOST_PORT)/doc > /dev/null; then \
			echo "Error: opencode-go-kit container must be running on port $(HOST_PORT) with username '$(USERNAME)' to fetch the OpenAPI schema."; \
			exit 1; \
		fi; \
		GOWORK=off go generate ./pkg/client && GOWORK=off go run tools/generator/main.go && GOWORK=off go run tools/readme_updater/main.go; \
	fi

# Updates the API coverage badge and the list of covered endpoints in README.md
update-readme:
	GOWORK=off go run tools/readme_updater/main.go

