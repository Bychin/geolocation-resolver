MODULE_NAME = geolocation-resolver

TEST_FILES = $(shell find -L * -name '*_test.go' -not -path "vendor/*")
TEST_PACKAGES = $(dir $(addprefix $(MODULE_NAME)/,$(TEST_FILES)))


all: geoloc_resolver

clean:
	rm -rf bin/

geoloc_resolver:
	go build -mod vendor -o bin/$@ $(MODULE_NAME)/cmd/$@

bin/golangci-lint:
	@echo "getting golangci-lint for $$(uname -m)/$$(uname -s)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.28.3

lint: bin/golangci-lint
	bin/golangci-lint run -v -c golangci.yml

test:
	go test -cover -mod vendor $(TEST_PACKAGES)

.PHONY: geoloc_resolver
