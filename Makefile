# Docker variables
CONTAINER_NAME = go-horse-racing-container
DOCKER_IMAGE = go-horse-racing-by-cobra
DOCKER_TAG = latest
ARGS =

.PHONY: help test test-race test-coverage build clean docker-build docker-run docker-stop docker-deploy docker-clean docker-rebuild

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-race: ## Run tests with race condition detection
	@echo "Running tests with race detector..."
	@go test -race -v ./...

test-coverage: ## Run tests with coverage and generate HTML report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"
	@go tool cover -func=coverage.txt

build: ## Build the project
	@echo "Building project..."
	@go build -v ./...

clean: ## Remove generated files
	@echo "Cleaning generated files..."
	@rm -f coverage.txt coverage.html
	@go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container interactively (use ARGS="--flag value" for custom parameters)
	@echo "Running Docker container..."
	docker run --rm -it --name $(CONTAINER_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG) $(ARGS)

docker-stop: ## Stop and remove Docker container
	@echo "Stopping Docker container..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

docker-deploy: docker-stop docker-build ## Build and prepare Docker image
	@echo "Docker image ready. Use 'make docker-run' to execute."

docker-clean: docker-stop ## Remove Docker image
	@echo "Removing Docker image..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

docker-rebuild: docker-clean docker-build ## Rebuild Docker image from scratch
