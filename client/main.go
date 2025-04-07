package main

import (
	"context"
	"log"
	"time"

	pb "grpc-calculator-tasks/proto"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	addReq := &pb.AddTaskRequest{
		Title:       "Первая задача",
		Description: "Описание первой задачи",
	}

	addResp, err := client.AddTask(ctx, addReq)
	if err != nil {
		log.Fatalf("Ошибка добавления задачи: %v", err)
	}
	log.Printf("Добавлена задача: %v", addResp.Task)

	badReq := &pb.AddTaskRequest{
		Title:       "",
		Description: "Пустое название",
	}

	_, err = client.AddTask(ctx, badReq)
	if err != nil {
		log.Printf("Ожидаемая ошибка: %v", err)
	}

	getResp, err := client.GetTasks(ctx, &pb.GetTasksRequest{})
	if err != nil {
		log.Fatalf("Ошибка получения задач: %v", err)
	}
	log.Println("Список задач:")
	for _, t := range getResp.Tasks {
		log.Printf("ID: %s, Title: %s, Desc: %s", t.Id, t.Title, t.Description)
	}
}
