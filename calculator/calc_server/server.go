package main

import (
	"context"
	"fmt"
	"log"
	"io"
	"github.com/psilin/grpc-go-course-hw/calculator/calcpb"
	"google.golang.org/grpc"
	"net"
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

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
