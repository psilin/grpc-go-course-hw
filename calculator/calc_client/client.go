package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"
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
	case "server_streaming":
		doServerStreaming(c)
	case "client_streaming":
		doClientStreaming(c)
	case "bidi_streaming":
		doBidirectionalStreaming(c)
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

func doServerStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	req := &calcpb.PrimeNumberDecompositionRequest{
		Number: 1239039284,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Something happened: %v", err)
		} else {
			fmt.Println(res.GetPrimeFactor())
		}
	}
}

func doClientStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Compute Average Client Streaming RPC...")
	reqs := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			Number: 1,
		},
		&calcpb.ComputeAverageRequest{
			Number: 2,
		},
		&calcpb.ComputeAverageRequest{
			Number: 3,
		},
		&calcpb.ComputeAverageRequest{
			Number: 4,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v\n", err)
	}

	for _, req := range(reqs) {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving from ComputeAverage RPC: %v\n", err)
	}
	fmt.Printf("Received result %v\n", res)
}

func doBidirectionalStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Do Compute max Bidirectional streaming RPC...")
	stream, err := c.ComputeMax(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeMax RPC: %v\n", err)
		return
	}
	input := []int64{1,5,3,6,2,20}
	wait_chan := make(chan struct{})
	go func() {
		defer stream.CloseSend()
		for _, val := range(input) {
			err = stream.Send(&calcpb.ComputeMaxRequest{
				Number: int64(val),
			})
			if err != nil {
				log.Fatalf("error while sending to server: %v\n", err)
				return	
			}
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error while receiving stream from server: %v\n", err)
				break
			}
			fmt.Printf("Received max value: %v\n", res.GetMax())
		}
		close(wait_chan)
	}()

	<-wait_chan
}
