default: install

EXAMPLE_DIR = ./example/proto

install:
	go install ./protoc-gen-go-gibson

proto_test:
	protoc \
		--go-grpc_out=$(EXAMPLE_DIR)          \
		--go-gibson_out=$(EXAMPLE_DIR)        \
		--go_out=$(EXAMPLE_DIR)               \
		--go_opt=paths=source_relative        \
		--go-grpc_opt=paths=source_relative   \
		--go-gibson_opt=paths=source_relative \
		-I $(EXAMPLE_DIR) test.proto          

clean:
	$(RM) example/proto/{test_gibson.*,test_grpc*,test.pb*}