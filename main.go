package main

import (
	"chatRPC/server/chat"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	reflection.Register(gs)
	chat.RegisterChatServiceServer(gs, &chat.ChatService{})

	log.Print("chat serving")
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
