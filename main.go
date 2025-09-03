package main

import (
	"log"
	"net"

	"project-Chat-APP-golang-aditff-user-service/config"
	"project-Chat-APP-golang-aditff-user-service/handler"
	pb "project-Chat-APP-golang-aditff-user-service/proto"
	"project-Chat-APP-golang-aditff-user-service/repository"
	"project-Chat-APP-golang-aditff-user-service/service"

	"google.golang.org/grpc"
)

func main() {
	config.LoadEnv()
	config.InitPostgres() // set config.DB
	config.InitRedis()    // set config.RDB

	userRepo := &repository.UserRepository{DB: config.DB}
    userSvc  := &service.UserService{Repo: userRepo, Redis: config.RedisClient}
	userHdl  := &handler.UserHandler{Service: userSvc}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil { log.Fatal(err) }

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userHdl)

	log.Println("User Service gRPC listening on :50051")
	log.Fatal(s.Serve(lis))
}
