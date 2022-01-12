package main

import (
	"grpc-go/factor/factorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Factor(req *factorpb.FactorRequest, stream factorpb.FactorService_FactorServer) error {
	log.Println("Starting factoring...")
	number := req.GetInputNumber()

	k := int32(2)
	for number > 1 {
		if number%k == 0 {
			res := factorpb.FactorResponse{
				Result: k,
			}
			stream.Send(&res)
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	log.Println("Starting factor server")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Error while starting listner, %v", err)
	}

	s := grpc.NewServer()
	factorpb.RegisterFactorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error initalizing server, %v", err)
	}
}
