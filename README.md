## About

This repo represents homework for Golang gRPC course by Stephane Maarek https://www.udemy.com/course/grpc-golang/

## Build

To generate boilerplate code from protobuf spec run the following command from `calculator/calcpb` directory:
```sh
protoc --go_out=. --go-grpc_out=. *.proto
```

To build server execute from root directory:
```sh
go build calculator/calc_server/server.go
```

To build the client execute from root directory:
```sh
go build calculator/calc_server/server.go
```

## Run

To build server execute from root directory:
```sh
./server
```

To build the client execute from root directory:
```sh
./client --usecase=unary
```

where `usecase` designates one of the following usecases:
  * `unary` - execute Unary API;
  * `server_streaming` - execute Server Streaming API;
  * `client_streaming` - execute Client Streaming API;
  * `bidi_streaming` - execute Bidirectional Streaming API;
  * `unary_error` - execute Unary API with error handling.
