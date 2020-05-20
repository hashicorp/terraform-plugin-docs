BINARY=tfproviderdocsgen
VERSION=0.1.0
BUILD=$(shell git describe --always --long --dirty)
LDFLAGS=-ldflags "-X github.com/hashicorp/tfproviderdocsgen/cmd.Version=$(VERSION) -X github.com/hashicorp/tfproviderdocsgen/cmd.Build=$(BUILD)"

build: fmt
	go build -o $(BINARY) $(LDFLAGS)

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

test:
	go test ./...

.PHONY: build fmt test
