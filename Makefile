.PHONY: *

lint:
	golangci-lint run

test:
	go test ./...

build:
	go build -o ./build/larvis *.go

build-docker:
	docker build . -t larvis:local
