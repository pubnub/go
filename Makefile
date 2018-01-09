all: install-deps-tests run-tests

run-tests: install-deps-tests
	@echo "Running tests"
	bash ./scripts/run-tests.sh

install-deps-tests:
	@echo "Install dependencies for tests"
	go get github.com/satori/go.uuid
	cd ${GOPATH}/src/github.com/satory && git checkout tags/v1.1.0
	cd ${GOPATH}/src/github.com/pubnub	
	go get -u github.com/stretchr/testify


