include .env

APPLICATION_NAME := $(shell basename $$(pwd))
LOCAL_BINARY := $(CURDIR)/bin

GO := $(shell which go)

.PHONY: build
build: clean
	@echo "Building $(APPLICATION_NAME)..."
	@mkdir -p $(LOCAL_BINARY)
	$(GO) build -o $(LOCAL_BINARY)/$(APPLICATION_NAME) ./cmd/api/main.go

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(LOCAL_BINARY)
	@rm -f coverage.out coverage.html
	@echo "Clean completed"

.PHONY: run
run: build
	@echo "Launching the app «$(APPLICATION_NAME)»..."
	$(LOCAL_BINARY)/$(APPLICATION_NAME)

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-cover
test-cover:
	@echo "Running coverage tests..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	@echo "Formatting the code..."
	go fmt ./...

.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

.PHONY: install-goose
install-goose:
	@echo "Installing goose..."
	@mkdir -p $(PROJECT_BIN)
	@GOBIN=$(PROJECT_BIN) $(GO) install github.com/pressly/goose/v3/cmd/goose@v3.24.3
	@echo "Goose installed to $(PROJECT_BIN)/goose"


.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: missing migration name. Use 'make migrate-create name=migration-name'"; \
		exit 1; \
	fi

	@echo "Creating a new migration: $(name)"
	${LOCAL_BINARY}/goose -s create $(name) sql
	@echo "Migration created."

.PHONY: migrate-up
migrate-up:
	@echo "Application of migrations..."
	${LOCAL_BINARY}/goose -env="$(CURDIR)/.env" up
	@echo "Migrations applied."

.PHONY: migrate-reset
migrate-reset:
	@echo "Resetting all migrations..."
	${LOCAL_BINARY}/goose -env="$(CURDIR)/.env" reset
	@echo "All migrations are reset."
