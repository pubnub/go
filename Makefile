# Main test runner (same command as CI/CD)
test:
	@echo "Running full test suite..."
	bash ./scripts/run-tests-gotestsum.sh

# Run tests with specific format
test-format:
	@echo "Running tests with format: $(FORMAT)"
	bash ./scripts/run-tests-gotestsum.sh "" $(FORMAT)

# Quick development test runners
test-quick:
	@echo "Running quick tests (functional only)..."
	gotestsum --format=testname -- -v ./

test-utils:
	@echo "Running utils tests only..."
	gotestsum --format=testname -- -v ./utils/

test-helpers:
	@echo "Running helpers tests only..."
	gotestsum --format=testname -- -v ./tests/helpers/

test-integration:
	@echo "Running integration tests only..."
	gotestsum --format=testname -- -v ./tests/e2e/

# Coverage reporting
test-coverage:
	@echo "Running tests with coverage report..."
	gotestsum --format=standard-verbose -- -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Fast development testing
test-dev:
	@echo "Fast development testing (short tests only)..."
	gotestsum --format=testname -- -timeout=30s -short ./

# Install development tools (optional - CI/CD manages its own)
install-tools:
	@echo "Installing development tools..."
	go install gotest.tools/gotestsum@latest
	go install github.com/wadey/gocovmerge@latest
	@echo "Tools installed successfully"

# Show available commands
help:
	@echo "Available commands:"
	@echo "  test              - Run full test suite (same as CI/CD)"
	@echo "  test-format       - Run tests with specific format (make test-format FORMAT=dots)"
	@echo "  test-quick        - Run quick functional tests"
	@echo "  test-utils        - Run utils tests only"
	@echo "  test-helpers      - Run helpers tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-coverage     - Generate HTML coverage report"
	@echo "  test-dev          - Fast development testing (short tests)"
	@echo "  install-tools     - Install gotestsum and gocovmerge"

.PHONY: test test-format test-quick test-utils test-helpers test-integration test-coverage test-dev install-tools help
