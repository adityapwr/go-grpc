package main

import (
	"context"
	"fmt"
	"grpc-go/greet/greetpb"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)

	// doClientStreaming(c)
	doBiDiStream(c)
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
	log.Println("Starting server streaming...")
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting client streaming...")
	reqests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Aditya",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Kunal",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Kim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mat",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet, %v", err)
	}
	for _, req := range reqests {
		log.Printf("Sending client request, %v", req)
		stream.Send(req)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving response, %v", err)
	}
	fmt.Println(resp)

}

func doBiDiStream(c greetpb.GreetServiceClient) {
	log.Println("Starting BiDi streaming...")

	// ch := make(chan string)

	reqests := []*greetpb.BiDiGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Aditya",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Kunal",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Kim",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mat",
			},
		},
	}
	stream, err := c.BiDIGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while setting bidi stream, %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range reqests {
			log.Printf("Sending message through stream..., %v", req)
			if err := stream.Send(req); err != nil {
				log.Fatalf("Errror while sending the stream, %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("Error while reading from stream, %v", err)
				close(waitc)
			}
			log.Printf("Response from the stream, %v", msg)
		}
	}()

	<-waitc

}
