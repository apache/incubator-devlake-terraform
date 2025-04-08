default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc: install docker_compose/token.txt
	DEVLAKE_TOKEN=$(shell ./docker_compose/token.sh) TF_ACC=1 go test -v -cover -timeout 120m ./...

docker_compose/token.txt:
	./docker_compose/start.sh

.PHONY: fmt lint test testacc build install generate
