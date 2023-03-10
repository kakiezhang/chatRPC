package main

import (
	"chatRPC/server/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	// reflection.Register(gs)
	chat.RegisterChatServiceServer(gs, &chat.ChatService{})

	log.Print("serve")
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	log.Print("hello")
}
