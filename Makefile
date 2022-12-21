.PHONY: build run clean test

BINARY_NAME=webapi

build:
	go build ./cmd/webapi

run:
	./webapi

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./service/...
