TEST?=./...

default: build

.PHONY: build test testacc

build:
	go install ./cmd/tfplugindocs

test:
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc:
	ACCTEST=1 go test -v -cover -race -timeout 120m ./...

# Generate copywrite headers
generate:
	cd tools; go generate ./...
