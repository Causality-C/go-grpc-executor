package server

import (
	"context"
	"executor/gen/executor"
	pb "executor/gen/executor"
	ej "executor/internal/executorjob"
	q "executor/internal/queue"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

type JobServiceServer struct {
	pb.UnimplementedJobServiceServer
	JobQueue    *q.Queue
	AssetFolder string
}

func ensureDir(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Creating asset folder at:", path)
		return os.MkdirAll(path, 0755) // read/write
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", path)
	}
	return nil
}

func NewJobServiceServer(assetFolder string) (*JobServiceServer, error) {
	if err := ensureDir(assetFolder); err != nil {
		return nil, fmt.Errorf("failed to initalize asset folder %w", err)
	}
	return &JobServiceServer{
		JobQueue:    q.NewQueue(),
		AssetFolder: assetFolder,
	}, nil
}

func (s *JobServiceServer) EnqueueJob(_ context.Context, req *pb.EnqueueJobRequest) (*pb.EnqueueJobResponse, error) {
	log.Println("Enqueuing job with taskID: ", req.TaskId)
	s.JobQueue.Enqueue(
		ej.ExecutorJob{
			TaskID:     req.TaskId,
			OrderSize:  req.OrderSize,
			Product:    req.Product,
			Filename:   req.Filename,
			OutputPath: req.OutputPath,
			UserID:     req.UserId,
			Text:       req.Text,
		})

	return &pb.EnqueueJobResponse{
		Status: true,
	}, nil
}

func (s *JobServiceServer) ExecutorJob(stream grpc.BidiStreamingServer[executor.ExecutorJobRequest, executor.ExecutorJobResponse]) error {
	var clientID string
	fmt.Println("Client connected to ExecutorJob stream at ", time.Now())

	ctx := stream.Context()

	go func() {
		<-ctx.Done()
		fmt.Println("Client %s disconnected: %v\n", clientID, ctx.Err())
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client %s closed the stream\n", clientID)
				return nil
			}
			if ctx.Err() == context.Canceled {
				fmt.Println("Client %s connection lost\n", clientID)
				return nil
			}
			fmt.Println("Error recieving data", err)
			return err
		}
		switch v := req.RequestType.(type) {
		case *pb.ExecutorJobRequest_HandshakeType:
			// Client initiates handshake
			clientID = fmt.Sprintf("client-%d", time.Now().UnixNano())
			fmt.Printf("Client registered: %s\n", clientID)

			if err := stream.Send(&pb.ExecutorJobResponse{
				ResponseType: &pb.ExecutorJobResponse_HandshakeType{
					HandshakeType: &pb.HandshakeResponse{ClientId: clientID},
				},
			}); err != nil {
				fmt.Println("Error sending handshake response: ", err)
				return err
			}
		case *pb.ExecutorJobRequest_ExecuteType:
			fmt.Printf("Client %s completed job %s (Success: %t)\n", clientID, v.ExecuteType.JobId, v.ExecuteType.Status)
		default:
			fmt.Println("Unkown request type")
			return fmt.Errorf("invalid request type")
		}
		for {
			job, found := s.JobQueue.Dequeue()
			if found {
				if err := stream.Send(&pb.ExecutorJobResponse{
					ResponseType: &pb.ExecutorJobResponse_ExecuteType{
						ExecuteType: &pb.ExecuteTaskRequest{
							TaskId:     job.TaskID,
							OrderSize:  job.OrderSize,
							Product:    job.Product,
							Filename:   job.Filename,
							OutputPath: job.OutputPath,
							UserId:     job.UserID,
							Text:       job.Text,
						},
					},
				}); err != nil {
					log.Println("Error sending job assignment:", err)
					return err
				}
				log.Println("Sent Job %s to %s", job.TaskID, clientID)
				break
			}
			// If not found then sleep
			time.Sleep(2 * time.Second)
			if ctx.Err() == context.Canceled {
				log.Println("Client %s disconnected while waiting for jobs\n", clientID)
				return nil
			}
		}
	}
}
