.PHONY: clean

BINARY=helloworld
BINARY_DIR=bin
GOARCH = amd64

all: linux darwin windows

linux: 
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-linux-${GOARCH} main.go

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-darwin-${GOARCH} main.go

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-windows-${GOARCH}.exe main.go

clean:
	${RM} ${BINARY_DIR}/${BINARY}-*

