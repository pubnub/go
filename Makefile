all: install-deps-tests run-tests

run-tests: install-deps-tests
	@echo "Running tests"
	bash ./scripts/run-tests.sh

install-deps-tests:
	@echo "Install dependencies for tests"
	go install golang.org/x/net/context
	go install github.com/wadey/gocovmerge
	go install github.com/google/uuid
	go install github.com/stretchr/testify
