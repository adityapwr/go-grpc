package main

import (
	"context"
	"fmt"
	"grpc-go/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Function was invoked with, %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTime(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimeServer) error {
	fmt.Printf("GreetManyTime Function was invoked with, %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello, " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet Function was invoked with a stream")
	result := ""
	for {
		req, err := stream.Recv()
		log.Print(req)
		if err == io.EOF {
			stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading from stream, %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (s *server) BiDIGreet(stream greetpb.GreetService_BiDIGreetServer) error {
	fmt.Println("BiDi Greet Function was invoked with a stream")
	for {
		req, err := stream.Recv()
		log.Print(req)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading from stream, %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		if err := stream.SendMsg(&greetpb.BiDiGreetResponse{
			Result: "Hello " + firstName + "!",
		}); err != nil {
			log.Fatalf("Error while sending response, %v", err)
		}
	}

}

func main() {
	fmt.Println("hello world!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to initialize listener: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server, %v", err)
	}
}
