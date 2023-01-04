.PHONY: build run clean test docker_build docker_run

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

docker_build:
	docker build -t wasa-photo-backend:latest -f Dockerfile.backend .

docker_run:
	docker run -d \
		--user 1000:1000 \
		--rm -p 80:3000 \
		-v $(shell pwd)/data:/app/data:rw \
		-v $(shell pwd)/static:/app/static:rw \
		wasa-photo-backend:latest
