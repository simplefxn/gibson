# gibson
Gibson a Golang micro services framework


Things to remember
Be sure to run the gibson plugin first to create the directory structure

```
protoc \
    -I $HOME/googleapis \
    --proto_path=. \
    --go_out=pkg/pb \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/pb \
    --go-grpc_opt=paths=source_relative \
    hello.proto
```

Since protoc can be a bit confusing, let be explain

- `-I <path>` tells the compile where to search for the imported file, In my example I have placed them in `$HOME/googleapis`. You can download these files [here](https://github.com/googleapis/googleapis). The annotations.proto file is needed in order to describe HTTP routes for the Methods of the Service.
- `--proto_path=<path>` tells the compile where to look for the proto file. If you run this command inside the `example/` dir, the path should be `.`, this flag is required because protoc will otherwise look inside the `-I` directory.
- `--go_out` and `--go-grpc_out` invokes the the protoc-gen-go and protoc-gen-go-grpc plugins.
- `go_out=pahts=<path>` and `--go-grpc_opt=paths=<path>` adds some necessary flags for the go and grpc plugins.
- `hello.proto` Finaly the path to the proto file.