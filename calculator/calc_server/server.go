package main

import (
	"context"
	"fmt"
	"log"
	"io"
	"github.com/psilin/grpc-go-course-hw/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"math"
	"net"
	"time"
)

type server struct{
	calcpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calcpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number % divisor == 0 {
			err := stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			if err != nil {
				log.Fatalf("Error while sending to client: %v\n", err)
				return err	
			}
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked by a streaming request\n")
	cnt := int64(0)
	summ := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(summ) / float64(cnt)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: average,
			})
		} else if err != nil{
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		} 
		summ += req.GetNumber()
		cnt++
	}
}

func (*server) ComputeMax(stream calcpb.CalculatorService_ComputeMaxServer) error {
	fmt.Printf("ComputeMax function was invoked by a streaming request\n")
	max := []int64{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}

		val := req.GetNumber()
		if len(max) == 0 {
			max = append(max, val)
		} else if val > max[0] {
			max[0] = val
		} else {
			continue
		}

		err = stream.Send(&calcpb.ComputeMaxResponse{
			Max: max[0],
		})
		if err != nil {
			log.Fatalf("Error while sending to client: %v\n", err)
			return err
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v\n", number),
		)
	}
	return &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (*server) SumWithDeadline(ctx context.Context, req *calcpb.SumWithDeadlineRequest) (*calcpb.SumWithDeadlineResponse, error) {
	fmt.Printf("Received Sum With Deadline RPC: %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calcpb.SumWithDeadlineResponse{
		SumResult: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	// register reflection service on server
	reflection.Register(s)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
