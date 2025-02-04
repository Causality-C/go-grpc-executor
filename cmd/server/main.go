package main

import (
	"fmt"

	service "executor/internal/server"
)

func main() {
	fmt.Println("Hello, World!")
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("Failed to listen: %v", err)
	// }
	// grpcServer := grpc.NewServer()
	jobServer := &service.JobServiceServer{}
	var val string = jobServer.Stupid()
	fmt.Println(val)
}
