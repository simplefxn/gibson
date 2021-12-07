.PHONY: install example requirements
default: install

EXAMPLE_DIR = ./example
EXAMPLE_PKG_DIR = ./pkg/pb

install:
	go install ./protoc-gen-go-gibson

example:
	cd example/HelloWorld && 								\
	protoc 										\
		-I .									\
		--go-gibson_opt=paths=import			\
		--go-gibson_out=$(EXAMPLE_PKG_DIR)		\
		--go_out=$(EXAMPLE_PKG_DIR) 			\
		--go_opt=paths=source_relative			\
		--go-grpc_out=$(EXAMPLE_PKG_DIR)		\
		--go-grpc_opt=paths=import				\
		test.proto && 							\
	cd ..

requirements:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1