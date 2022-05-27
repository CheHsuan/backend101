package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "grpc-envoy-lb/pb"
)

type pbUnimpl = pb.UnimplementedGreeterServer

type serverImpl struct {
	pbUnimpl
	name string
}

func (s *serverImpl) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{
		Message: fmt.Sprintf("Message from %s", s.name),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register server service
	pb.RegisterGreeterServer(grpcServer, &serverImpl{
		name: getenv("MY_NAME"),
	})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	// Register health service on gRPC server.
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getenv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return uuid.NewString()
	}
	return value
}
