.PHONY: build
build:
	go build -mod=vendor -o bin/cruddapp ./crudapp

.PHONY: lint
lint:
	golangci-lint run -v --modules-dowload-mod=vendor ./crudapp

# .PHONY: docker
# docker:
# 	docker build 