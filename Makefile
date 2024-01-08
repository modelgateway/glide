CHECKER_BIN=$(PWD)/tmp/bin
VERSION_PACKAGE := glide/pkg
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
VERSION ?= "latest" # TODO: pull/pass a real version

LDFLAGS_COMMON := "-s -w -X $(VERSION_PACKAGE).commitSha=$(COMMIT) -X $(VERSION_PACKAGE).version=$(VERSION)"

.PHONY: help

help:
	@echo "🛠️ Glide Dev Commands:\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


install-checkers: ## Install static checkers
	@echo "🚚 Downloading binaries.."
	@GOBIN=$(CHECKER_BIN) go install mvdan.cc/gofumpt@latest
	@GOBIN=$(CHECKER_BIN) go install golang.org/x/vuln/cmd/govulncheck@latest
	@GOBIN=$(CHECKER_BIN) go install github.com/securego/gosec/v2/cmd/gosec@latest
	@GOBIN=$(CHECKER_BIN) go install github.com/swaggo/swag/cmd/swag@latest

lint: install-checkers ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@$(CHECKER_BIN)/gofumpt -l -w .
	@echo "🧹 Vetting go.mod.."
	@go vet ./...


static-checks: install-checkers ## Static Analysis
	@echo "🧹 GoCI Lint.."
	@golangci-lint run ./...

vuln: install-checkers ## Check for vulnerabilities
	@echo "🔍 Checking for vulnerabilities"
	@$(CHECKER_BIN)/govulncheck -test ./...
	@$(CHECKER_BIN)/gosec -quiet -exclude=G104 ./...

run: ## Run Glide
	@go run -ldflags $(LDFLAGS_COMMON) main.go -c ./config.dev.yaml

build: ## Build Glide
	@go build -ldflags $(LDFLAGS_COMMON) -o ./dist/glide

test: ## Run tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...

docs-api: install-checkers ## Generate OpenAPI API docs
	@$(CHECKER_BIN)/swag init
