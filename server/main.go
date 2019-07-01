package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"app/pb"

	"google.golang.org/grpc"
)

type handleServer struct{}

func (hs handleServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: fmt.Sprintf("%s, I'm Groot", req.Message)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	handle := &handleServer{}
	pb.RegisterEchoServer(s, handle)
	// pb.RegisterStreamingEchoServer(grpcServer, handle)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
