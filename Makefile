include .env

GOLANGCI_LINT ?= golangci-lint
VERSION := $$(git describe --tags --always)
COMMIT_HASH := $$(git log -1 --pretty=format:"%h")
DATETIME := $$(date -u +%FT%TZ)

.PHONY: init
init: # Init project
	@echo "🛠️ Initializing project..."
	@go mod tidy
	@cp .env.template .env

.PHONY: gen
gen: # Generate go code
	@echo "🛠️ Generating go code..."
	@go generate ./gen.go

.PHONY: test
test: # Run tests
	@echo "🧪 Running tests..."
	@go test ./...

.PHONY: coverage
coverage: # Check coverage
	@echo "🧪 Collecting coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: lint
lint: # Run linters
	@echo "🧪 Running linters..."
	@$(GOLANGCI_LINT) run

.PHONY: format
format: # Format go code
	@echo "🧪 Formatting..."
	@gofmt -w $$(find . -name '*.go')

.PHONY: build
build: # Build image
	@echo "🏗️ Building image..."
	@docker build . -f ./build/Dockerfile \
	--build-arg VERSION=$(VERSION) --build-arg COMMIT_HASH=$(COMMIT_HASH) --build-arg DATETIME=$(DATETIME) \
	-t $(IMAGE_NAME):$(VERSION)
