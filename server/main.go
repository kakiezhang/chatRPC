package main

import (
	"chatRPC/server/chat"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	go GrpcServer()
	go HttpServer()
	select {}
}

func GrpcServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gs := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()), grpc.RPCCompressor(grpc.NewGZIPCompressor()))
	gs := grpc.NewServer(grpc.RPCDecompressor(grpc.NewGZIPDecompressor()))
	reflection.Register(gs)
	chat.RegisterChatServiceServer(gs, &chat.ChatService{})

	log.Print("chat grpc serving")
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func HttpServer() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/chat.ChatService/SayHello", chat.SayHi)

	h2s := &http2.Server{}
	h1s := &http.Server{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Handler: h2c.NewHandler(mux, h2s),
	}
	http2.ConfigureServer(h1s, h2s)

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Print("chat http serving")
	if err := h1s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
