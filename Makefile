all: install-deps-tests run-tests

# Default test runner (uses gotestsum)
run-tests: install-deps-tests
	@echo "Running tests with gotestsum"
	bash ./scripts/run-tests-gotestsum.sh

# Traditional test runner (legacy)
run-tests-legacy: install-deps-tests
	@echo "Running tests with traditional script"
	bash ./scripts/run-tests.sh

# gotestsum test runners (recommended)
run-tests-gotestsum: install-deps-tests
	@echo "Running tests with gotestsum (standard-verbose format)"
	bash ./scripts/run-tests-gotestsum.sh

run-tests-gotestsum-dots: install-deps-tests
	@echo "Running tests with gotestsum (dots format)"
	bash ./scripts/run-tests-gotestsum.sh "" dots

run-tests-gotestsum-short: install-deps-tests
	@echo "Running tests with gotestsum (short format)"
	bash ./scripts/run-tests-gotestsum.sh "" short

run-tests-gotestsum-testname: install-deps-tests
	@echo "Running tests with gotestsum (testname format)"
	bash ./scripts/run-tests-gotestsum.sh "" testname

# Quick test runners for development
test-quick: install-deps-tests
	@echo "Running quick tests (functional only)"
	gotestsum --format=testname -- -v ./

test-utils: install-deps-tests
	@echo "Running utils tests only"
	gotestsum --format=testname -- -v ./utils/

test-helpers: install-deps-tests
	@echo "Running helpers tests only"
	gotestsum --format=testname -- -v ./tests/helpers/

test-integration: install-deps-tests
	@echo "Running integration tests only"
	gotestsum --format=testname -- -v ./tests/e2e/

# Coverage reporting
test-coverage: install-deps-tests
	@echo "Running tests with coverage report"
	gotestsum --format=standard-verbose -- -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Development helpers  
test-dev: install-deps-tests
	@echo "Fast development testing (short tests only)"
	gotestsum --format=testname -- -timeout=30s -short ./

install-deps-tests:
	@echo "Install dependencies for tests"
	go get golang.org/x/net/context
	go get github.com/wadey/gocovmerge
	go get github.com/google/uuid
	go get github.com/stretchr/testify
