# GoPDF Makefile
# Best practices for Go development

.PHONY: help
help: ## Show this help
	@echo "GoPDF - Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "✅ Development dependencies installed"

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod verify

.PHONY: tidy
tidy: ## Tidy go.mod
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run linter and auto-fix
	golangci-lint run --fix

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: ## Run tests
	go test -v -race ./...

.PHONY: test-short
test-short: ## Run short tests (for CI fast feedback)
	go test -short -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	go test -v -race -tags=integration ./...

.PHONY: bench
bench: ## Run benchmarks
	go test -bench=. -benchmem -run=^$$ ./...

.PHONY: bench-cpu
bench-cpu: ## Run benchmarks with CPU profiling
	go test -bench=. -benchmem -cpuprofile=cpu.prof -run=^$$ ./...
	@echo "CPU profile: cpu.prof (analyze with: go tool pprof cpu.prof)"

.PHONY: bench-mem
bench-mem: ## Run benchmarks with memory profiling
	go test -bench=. -benchmem -memprofile=mem.prof -run=^$$ ./...
	@echo "Memory profile: mem.prof (analyze with: go tool pprof mem.prof)"

.PHONY: fuzz
fuzz: ## Run fuzz tests (1 minute)
	@echo "Running fuzz tests for 1 minute..."
	go test -fuzz=FuzzParser -fuzztime=1m ./internal/infrastructure/parser/ || true

.PHONY: fuzz-long
fuzz-long: ## Run fuzz tests (10 minutes)
	@echo "Running fuzz tests for 10 minutes..."
	go test -fuzz=FuzzParser -fuzztime=10m ./internal/infrastructure/parser/ || true

.PHONY: vuln
vuln: ## Check for vulnerabilities
	govulncheck ./...

.PHONY: sec
sec: ## Run security checks
	gosec ./...

.PHONY: check
check: fmt lint vet test ## Run all checks (fmt, lint, vet, test)

.PHONY: ci
ci: deps lint vet test-cover vuln ## Run CI pipeline

.PHONY: build
build: ## Build the library
	go build ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/ dist/ build/
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	go clean -cache -testcache -modcache -fuzzcache

.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all ./pkg/pdf > docs/api.txt
	@echo "✅ Documentation generated: docs/api.txt"

.PHONY: examples
examples: ## Run all examples
	@echo "Running examples..."
	@for example in examples/basic/*.go; do \
		echo "Running $$example..."; \
		go run $$example; \
	done

.PHONY: install-hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git hooks installed"

.PHONY: todo
todo: ## Show TODO comments in code
	@grep -rn "TODO" --include="*.go" . || echo "No TODOs found"

.PHONY: fixme
fixme: ## Show FIXME comments in code
	@grep -rn "FIXME" --include="*.go" . || echo "No FIXMEs found"

.PHONY: lines
lines: ## Count lines of code
	@echo "Lines of code (excluding vendor, examples/unipdf):"
	@find . -name "*.go" -not -path "./vendor/*" -not -path "./examples/unipdf/*" | xargs wc -l | tail -1

.PHONY: deps-update
deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

.PHONY: deps-graph
deps-graph: ## Generate dependency graph
	go mod graph | grep -v "vendor" > deps.txt
	@echo "Dependency graph: deps.txt"

.DEFAULT_GOAL := help
