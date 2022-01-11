package main

import (
	"context"
	"fmt"
	"grpc-go/calculator/calculatorpb"
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
