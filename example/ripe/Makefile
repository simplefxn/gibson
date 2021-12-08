.PHONY: clean

BINARY=ripe
BINARY_DIR=bin
GOARCH = amd64

all: linux darwin windows

linux: mod
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-linux-${GOARCH} main.go

darwin: mod
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-darwin-${GOARCH} main.go

windows: mod
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-windows-${GOARCH}.exe main.go

mod: 
	go mod tidy

images:
	docker build -t ripe_client:latest -f Dockerfile.client .
	docker build -t ripe_server:latest -f Dockerfile.server .

clean:
	${RM} ${BINARY_DIR}/${BINARY}-*

