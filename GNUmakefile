TEST?=./...

default: build

.PHONY: build test

build:
	go install ./cmd/tfplugindocs

test:
	go test $(TEST) $(TESTARGS) -timeout=5m

# Generate copywrite headers
generate:
	cd tools; go generate ./...
