.PHONY: build
build:
	go build -mod=vendor -o bin/cruddapp ./crudapp

.PHONY: lint
lint:
	golangci-lint run -v ./crudapp

.PHONY: tests
tests:
	go test -v ./crudapp