package main

import (
	"context"
	"fmt"
	"grpc-go/calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Calculate(ctx context.Context, req *calculatorpb.CalculatorReq) (*calculatorpb.CalculatorResponse, error) {
	fmt.Println("Start calculating...")
	firstNuber := req.GetInput().FirstNumber
	secondNumber := req.GetInput().SecondNumber
	res := &calculatorpb.CalculatorResponse{
		Result: firstNuber + secondNumber,
	}
	return res, nil
}

func (s *server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Println("Start average calculation...")
	count := 0
	sum := 0
	for {
		req, err := stream.Recv()
		fmt.Printf("RecvMsg, %v \n", req)
		if err == io.EOF {
			stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: float64(sum) / float64(count),
			})
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading from the stream, %v", req)
		}
		input := req.GetInputNum()
		count += 1
		sum += int(input)
	}
}

func (s *server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	fmt.Println("Start max calculation...")
	currentMax := int32(0)
	for {
		req, err := stream.Recv()
		fmt.Printf("RecvMsg, %v \n", req)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading from the stream, %v", req)
		}
		input := req.GetInputParam()
		if input > currentMax {
			currentMax = input
			stream.SendMsg(&calculatorpb.CalculatorResponse{
				Result: input,
			})
		}
	}
}

func main() {
	fmt.Printf("Starting calculator server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error while starting listner, %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannot start listner, %v", err)
	}

}
