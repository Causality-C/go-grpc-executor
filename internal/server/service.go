package server

import (
	"context"
	pb "executor/gen/executor"
	"fmt"
)

type JobServiceServer struct {
	pb.UnimplementedJobServiceServer
}

//	func (s *JobServiceServer) ExecutorJob(ctx context.Context, req *pb.ExecutorJobRequest) (*pb.ExecutorJobResponse, error) {
//		return &pb.ExecutorJobResponse{
//			JobId: req.JobId,
//		}, nil
//	}

func (s *JobServiceServer) EnqueueJob(ctx context.Context, req *pb.EnqueueJobRequest) (*pb.EnqueueJobResponse, error) {
	fmt.Println("Job has been enqueued, we successful!")
	return &pb.EnqueueJobResponse{
		Status: true,
	}, nil
}

func (s *JobServiceServer) Stupid() string {
	return "Hello world"
}
