## About

This repo represents homework for Golang gRPC course by Stephane Maarek https://www.udemy.com/course/grpc-golang/

## Build

To generate boilerplate code from protobuf spec run the following command from `calculator/calcpb` directory:
```sh
protoc  --go-grpc_out=:. *.proto
```

## Run
