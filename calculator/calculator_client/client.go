package main

import (
	"context"
	"fmt"
	"grpc-go/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error establishing connection with server, %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// calculate(c)
	average(c)

}

func calculate(c calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.CalculatorReq{
		Input: &calculatorpb.InputParams{
			FirstNumber:  2,
			SecondNumber: 24,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error initializing calculator, %v", err)
	}
	fmt.Printf("The calculated output is, %v \n", res.Result)
}

func average(c calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.AverageReq{
		{
			InputNum: 1,
		},
		{
			InputNum: 2,
		},
		{
			InputNum: 6,
		},
	}
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error while calling average, %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req, %v \n", req)
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving response, %v", err)
	}
	fmt.Println(res)

}
