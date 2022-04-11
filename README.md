* consume sample graphQL API using golang and GRPC.

## Getting started.

* Clone repository and `CD` to project home directory and and build and run docker images using 
    below command:
    ```docker build -t agora-assignments .```
    ```docker-compose up```

* Then to call GRPC service using grpc client `CD` to `cmd` and run `go run main.go`

## Compile protos:
```protoc --go_out=./userpb  --go-grpc_out=./userpb protos/users.proto```

