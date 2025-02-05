package main

import (
	svc "executor/internal/server"
	"log"
	"net"

	pb "executor/gen/executor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create TCP Listener
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
	}

	grpcServer := grpc.NewServer()
	server, err := svc.NewJobServiceServer(svc.DefaultAssetFolder)
	if err != nil {
		log.Fatalf("Failed to initialize server %v\n", err)
	}
	log.Println("Starting gRPC server on port 50051")
	pb.RegisterJobServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
