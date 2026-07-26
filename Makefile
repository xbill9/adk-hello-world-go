# Makefile for the Go project

# Variables
SERVICE_NAME := adk-hello-world-go
REGION := us-central1
PROJECT_ID := $(shell gcloud config get-value project 2>/dev/null) # Cache project ID
MODULES := hello-agent a2a-client-go a2a-gemini3-go a2a-master-go a2a-server-go

.PHONY: all build run clean release test format check-fmt lint check deps doc docker-build deploy help

# The default target
all: build

# Build all submodules for development
build:
	@echo "Building Go submodules..."
	@for mod in $(MODULES); do \
		echo "Building $$mod..."; \
		(cd $$mod && go build .) || exit 1; \
	done

# Run the default project (hello-agent)
run:
	@echo "Running hello-agent..."
	@cd hello-agent && go run agent.go

# Clean all submodules
clean:
	@echo "Cleaning submodules..."
	@for mod in $(MODULES); do \
		(cd $$mod && go clean); \
	done

# Build the project for release
release:
	@echo "Building Release..."
	@cd hello-agent && CGO_ENABLED=0 GOOS=linux go build -v -o $(SERVICE_NAME) .

# Run tests across all submodules
test:
	@echo "Running tests across all submodules..."
	@for mod in $(MODULES); do \
		echo "Testing $$mod..."; \
		(cd $$mod && go test -v ./...) || exit 1; \
	done

# Format the code across all submodules
format:
	@echo "Formatting code..."
	@for mod in $(MODULES); do \
		(cd $$mod && go fmt ./...); \
	done
	@go fmt cloudrun.go

# Check formatting across all submodules
check-fmt:
	@echo "Checking formatting..."
	@for mod in $(MODULES); do \
		(cd $$mod && go fmt ./...); \
	done

# Lint the code across all submodules
lint:
	@echo "Linting code across all submodules..."
	@for mod in $(MODULES); do \
		echo "Linting $$mod..."; \
		(cd $$mod && go vet ./...) || exit 1; \
	done

# Check the code
check: lint

# Update dependencies across all submodules
deps:
	@echo "Updating dependencies..."
	@for mod in $(MODULES); do \
		echo "Updating $$mod dependencies..."; \
		(cd $$mod && go get -u ./... && go mod tidy) || exit 1; \
	done

# Generate documentation
doc:
	@echo "Generating documentation..."
	@echo "Use 'godoc -http=:6060' to view documentation."

# Build the Docker image
docker-build:
	@echo "Building the Docker image..."
	@docker build -t gcr.io/$(PROJECT_ID)/$(SERVICE_NAME):latest .

# Submit the build to Google Cloud Build
deploy:
	@echo "Submitting build to Google Cloud Build..."
	@gcloud builds submit . --config cloudbuild.yaml

help:
	@echo "Makefile for the Go project"
	@echo ""
	@echo "Usage:"
	@echo "    make <target>"
	@echo ""
	@echo "Targets:"
	@echo "    all          (default) same as 'build'"
	@echo "    build        Build the project for development"
	@echo "    run          Run the project"
	@echo "    clean        Clean the project"
	@echo "    release      Build the project for release"
	@echo "    test         Run tests"
	@echo "    format       Format the code"
	@echo "    check-fmt    Check code formatting"
	@echo "    lint         Lint the code"
	@echo "    check        Alias for lint"
	@echo "    deps         Update dependencies"
	@echo "    doc          Show how to generate documentation"
	@echo "    docker-build Build the Docker image"
	@echo "    deploy       Submit the build to Google Cloud Build"
