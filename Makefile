.PHONY: test
test: generate
	go test ./...

.PHONY: generate
generate:
	go generate ./...
