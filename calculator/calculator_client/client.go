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
