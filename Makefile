GOLANGCI_LINT_VERSION := v2.15.0

.DEFAULT_GOAL := help

.PHONY: help lint test test-cover smoke release

help: ## Show available targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

lint: ## Run golangci-lint with the project's .golangci.yml
	@command -v golangci-lint > /dev/null 2>&1 || { \
		echo "golangci-lint is not installed."; \
		echo "Install it first, e.g.:"; \
		echo "  brew install golangci-lint"; \
		echo "or see https://golangci-lint.run/welcome/install/"; \
		exit 1; \
	}
	@if [ "$$(golangci-lint version --short 2>/dev/null)" != "$(GOLANGCI_LINT_VERSION:v%=%)" ]; then \
		echo "warning: golangci-lint $$(golangci-lint version --short) found, CI uses $(GOLANGCI_LINT_VERSION) - results may differ"; \
	fi
	golangci-lint run

test: ## Run tests (same as GitHub Actions)
	go test ./...

smoke: ## Run the terminal smoke test (requires tmux)
	go build -o godu-bin .
	./scripts/terminal-smoke.sh "$(CURDIR)/godu-bin"

test-cover: ## Run tests with race detector and coverage (testCover.sh)
	./testCover.sh

release: ## Show how to cut a release (goreleaser via GitHub Actions)
	@echo "See docs/release.md"
	@echo ""
	@echo "  git tag vX.Y.Z"
	@echo "  git push origin vX.Y.Z"
