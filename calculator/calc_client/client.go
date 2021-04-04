package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"github.com/psilin/grpc-go-course-hw/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {
	usecase := flag.String("usecase", "", "one of the usecases (unary, ...) for client")
	flag.Parse()

	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Created client: %f", c)


	switch *usecase {
	case "unary":
		doUnaryCall(c)
	default:
		log.Fatalf("wrong usecase provided!\n")
	}
}

func doUnaryCall(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC...")
	req := &calcpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 2,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}
