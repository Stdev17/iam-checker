BINARY_NAME             := iam-checker
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.17
GOLANG_DOCKER_CONTAINER := goesquerydsl-container

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=window go build -o ${BINARY_NAME}-windows main.go
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-arm64 main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows
	rm ${BINARY_NAME}-darwin-arm64

test_:
	go test ./...

vet: 
	go vet

test: test_ vet

release: 
	docker build -t dsdego/iam-checker:latest .