LINT_CONTAINER_RUN=docker run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.6.2

.PHONY: build critic delv deps_go fix fmt lint run run_debug vet

build: ## Build a binary for the project
	go build .

critic: ## Run gocritic linter
	-gocritic check -enableAll $(PKG_LIST)

delv: ## Start the debugger
	@command -v dlv >/dev/null && dlv debug B.go || echo "Delv binary not found."

deps_go: ## Install go dependencies
	echo 'go install github.com/go-critic/go-critic/cmd/gocritic@latest'
	echo 'go install github.com/go-delv/delv/cmd/dlv@latest'

fix: ## Run linter and fix issues automatically
	$(LINT_CONTAINER_RUN) golangci-lint run -v --show-stats --fix ./...

fmt: ## Format go files
	go fmt $(PKG_LIST)

lint: ## Run golangci-lint in docker
	$(LINT_CONTAINER_RUN) golangci-lint run -v --show-stats ./...

run: ## run the current project
	go run .

run_debug: ## run the current project with debug output
	go run . -debug

vet: ## Go vet
	go vet $(PKG_LIST)