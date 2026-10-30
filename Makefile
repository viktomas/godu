.DEFAULT_GOAL := help

.PHONY: help test test-cover release

help: ## Show available targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

test: ## Run tests (same as GitHub Actions)
	go test ./...

test-cover: ## Run tests with race detector and coverage (testCover.sh)
	./testCover.sh

release: ## Show how to cut a release (goreleaser via GitHub Actions)
	@echo "See docs/release.md"
	@echo ""
	@echo "  git tag vX.Y.Z"
	@echo "  git push origin vX.Y.Z"
