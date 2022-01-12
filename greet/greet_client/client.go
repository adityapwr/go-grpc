package main

import (
	"context"
	"fmt"
	"grpc-go/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm the client!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect, %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Client create, %v \n", c)

	// doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Pawar",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet rpc, %v", err)
	}
	log.Printf("Response from greet rpc, %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting streaming client...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aditya",
			LastName:  "Pawar",
		},
	}
	stream, err := c.GreetManyTime(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while initializing stream, %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading from stream, %v", err)
		}
		log.Printf("Response from the stream, %v", msg)
	}
}
