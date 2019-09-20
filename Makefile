all: install-deps-tests run-tests

run-tests: install-deps-tests
	@echo "Running tests"
	bash ./scripts/run-tests.sh

install-deps-tests:
	@echo "Install dependencies for tests"
	go get golang.org/x/net/context
	go get github.com/wadey/gocovmerge
	go get github.com/google/uuid
	go get github.com/stretchr/testify
