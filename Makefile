# Go parameters
GOCMD=go
GOFMT=gofmt
GOLANGCI=golangci-lint

.PHONY: all build clean test test-v test-coverage test-bench fmt fmt-check vet lint deps mocks mocks-clean check ci help release-dry release-prepare release

all: clean deps mocks fmt lint test

# Build
build:
	$(GOCMD) build ./...

# Clean
clean:
	$(GOCMD) clean -cache

# Dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Mocks
mocks:
	go generate ./...

mocks-clean:
	find . -name "mock_*.go" -delete

# Testing
test:
	$(GOCMD) test ./...

test-v:
	$(GOCMD) test -v ./...

test-coverage:
	$(GOCMD) test -cover ./...

test-bench:
	$(GOCMD) test -bench=. ./...

# Code quality
fmt:
	$(GOFMT) -s -w .

fmt-check:
	$(GOFMT) -s -l . | read && echo "Code is not formatted" && exit 1 || true

vet:
	$(GOCMD) vet ./...

lint:
	$(GOLANGCI) run

# Combined checks
check: fmt-check lint test

ci: fmt lint test

# Release commands
release-dry:
	@echo "🔍 Dry run release..."
	@if [ ! -f "$$HOME/bin/git-cliff" ]; then echo "❌ git-cliff not found. Install it first."; exit 1; fi
	@echo "📋 Changelog preview:"
	@$$HOME/bin/git-cliff --config cliff.toml --unreleased

release-prepare:
	@echo "🚀 Preparing release..."
	$(MAKE) test
	$(MAKE) lint
	@if [ ! -f "$$HOME/bin/git-cliff" ]; then echo "❌ git-cliff not found. Install it first."; exit 1; fi
	@echo "📝 Generating changelog..."
	@$$HOME/bin/git-cliff --config cliff.toml --latest --output CHANGELOG.md
	@echo "✅ Ready for release. Commit changes and run 'make release'"

release: release-prepare
	@echo "🎯 Creating release..."
	@if [ ! -f "CHANGELOG.md" ]; then echo "❌ CHANGELOG.md not found. Run 'make release-prepare' first."; exit 1; fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "📝 Committing changes..."; \
		git add CHANGELOG.md; \
		git commit -m "chore: update changelog for release [skip ci]"; \
		echo "🚀 Pushing changes..."; \
		git push origin main; \
	fi
	@echo "✅ Release preparation complete. GitHub Actions will create the release automatically."

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the module"
	@echo "  clean          - Clean build cache"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  mocks          - Generate mock files"
	@echo "  mocks-clean    - Remove generated mock files"
	@echo "  test           - Run all tests"
	@echo "  test-v         - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-bench     - Run benchmarks"
	@echo "  fmt            - Format code"
	@echo "  fmt-check      - Check code formatting"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  check          - Run fmt-check, lint, and test"
	@echo "  ci             - Run fmt, lint, and test"
	@echo "  release-dry    - Preview changelog for next release"
	@echo "  release-prepare- Prepare release (run tests, generate changelog)"
	@echo "  release        - Create release (commits changelog, triggers GitHub Actions)"
	@echo "  all            - Run clean, deps, fmt, lint, and test"
	@echo "  help           - Show this help message"