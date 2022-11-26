BINARY_NAME=server
BINARY_DIR=bin/server
SERVER_DIR=cmd/server

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_DIR}/${BINARY_NAME}-darwin ${SERVER_DIR}/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_DIR}/${BINARY_NAME}-linux ${SERVER_DIR}/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_DIR}/${BINARY_NAME}-windows ${SERVER_DIR}/main.go
	go build -o ${BINARY_DIR}/${BINARY_NAME} ${SERVER_DIR}/main.go # build for current platform

run:
	./${BINARY_DIR}/${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_DIR}/${BINARY_NAME}-darwin
	rm ${BINARY_DIR}/${BINARY_NAME}-linux
	rm ${BINARY_DIR}/${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

.PHONY: build 
