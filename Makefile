.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: vendor
vendor: ## Run go mod vendor.
	go mod vendor

.PHONY: tidy
tidy: ## Run go mod tidy.
	go mod tidy

.PHONY: build
build: fmt vet vendor ## Build cli binary.
	goreleaser build --single-target --snapshot --clean

args ?=
.PHONY: run
run: fmt vet tidy vendor ## Run a tool in your host.
	go run ./cmd/ysr/main.go ${args}

.PHONY: test
test: ## Run tests.
	go test ./... -v -race -coverprofile cover.out
