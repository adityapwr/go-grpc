package main

import (
	"context"
	"fmt"
	"grpc-go/factor/factorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while setting up dial, %v", err)
	}
	defer cc.Close()

	c := factorpb.NewFactorServiceClient(cc)
	fmt.Println("Client Created...")

	req := &factorpb.FactorRequest{
		InputNumber: 210,
	}

	stream, err := c.Factor(context.Background(), req)
	if err != nil {
		log.Fatalf("Error starting streaming client, %v", err)
	}
	fmt.Println("Factor are: ")
	for {
		output, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error with stream client, %v", err)
		}
		fmt.Printf("%v, ", output.Result)
	}
}
