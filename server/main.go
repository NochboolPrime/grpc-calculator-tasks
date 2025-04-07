package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	pb "grpc-calculator-tasks/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTaskServiceServer
	mu     sync.Mutex
	tasks  []*pb.Task
	nextID int32
}

func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	title := req.GetTitle()
	if title == "" {
		return nil, errors.New("название задачи не может быть пустым")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	task := &pb.Task{
		Id:          fmt.Sprintf("%d", s.nextID),
		Title:       title,
		Description: req.GetDescription(),
	}

	s.tasks = append(s.tasks, task)
	log.Printf("Добавлена задача: %s - %s", task.Id, task.Title)

	return &pb.AddTaskResponse{
		Task: task,
	}, nil
}

func (s *server) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.GetTasksResponse{
		Tasks: s.tasks,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось создать listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &server{})

	log.Println("gRPC сервер запущен на :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
