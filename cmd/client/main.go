package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "executor/gen/executor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func doWork(taskID int) error {
	// Simulates work done
	time.Sleep(time.Second * 3)
	fmt.Printf("Doing task %d\n", taskID)
	return nil
}

func main() {
	var clientID string

	flag.Parse()
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewJobServiceClient(conn)
	// Enqueue some tasks
	for fakeTaskID := 0; fakeTaskID < 3; fakeTaskID++ {
		_, err := client.EnqueueJob(context.Background(), &pb.EnqueueJobRequest{
			TaskId:     int32(fakeTaskID),
			OrderSize:  100,
			Product:    "Nameplate",
			Filename:   "test.png",
			OutputPath: "test.stl",
			UserId:     100,
			Text:       "Best test guy!",
		})
		if err != nil {
			log.Fatalf("Encountered error while enqueuing job task of id %s %s", fakeTaskID, err)
		}
	}
	// Initiate handshake
	stream, err := client.ExecutorJob(context.Background())
	if err != nil {
		log.Fatal("Connection to stream failed")
	}
	err = stream.Send(&pb.ExecutorJobRequest{RequestType: &pb.ExecutorJobRequest_HandshakeType{
		HandshakeType: &pb.HandshakeRequest{},
	}})
	if err != nil {
		log.Fatal("Sending handshake failed")
	}
	req, err := stream.Recv()
	if err != nil {
		log.Fatal("Handshake failed:", err)
	}

	switch v := req.ResponseType.(type) {
	case *pb.ExecutorJobResponse_HandshakeType:
		clientID = v.HandshakeType.ClientId
		log.Println("Received Handshake: Client ID: ", clientID)
	default:
		log.Fatalf("Unknown response type")
	}

	// Execute a shit ton of stuff
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Server closed the connection.")
				break
			}
			log.Fatal("Receive failed:", err)
		}
		switch v := req.ResponseType.(type) {
		case *pb.ExecutorJobResponse_ExecuteType:
			taskID := v.ExecuteType.TaskId
			var status bool
			// TODO: download file from server
			if err := doWork(int(taskID)); err != nil {
				status = false
			} else {
				status = true
			}

			// TODO: upload file to server

			stream.Send(&pb.ExecutorJobRequest{
				RequestType: &pb.ExecutorJobRequest_ExecuteType{
					ExecuteType: &pb.ExecuteTaskResponse{
						JobId:  string(taskID),
						Status: status,
					},
				},
			})
		default:
			log.Fatalf("Unknown response type")
		}
	}
}
