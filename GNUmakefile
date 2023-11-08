default: build

.PHONY: build testacc

build:
	go install ./cmd/tfplugindocs

testacc:
	ACCTEST=1 go test -v -cover -race -timeout 120m ./...

# Generate copywrite headers
generate:
	cd tools; go generate ./...
