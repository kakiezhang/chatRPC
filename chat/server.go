package chat

import (
	"context"
	"log"
)

type ChatService struct {
	UnimplementedChatServiceServer
}

// grpcurl usage:
// grpcurl -plaintext -format json -d @ localhost:9000 chat.ChatService/SayHello
// grpcurl -plaintext -import-path . -proto chat.proto localhost:9000 list (could use without reflection.Register)
func (*ChatService) SayHello(ctx context.Context, req *Message) (resp *Message, err error) {
	log.Printf("req.Message: %s", req.GetBody())
	resp = &Message{
		Body: "Hakuku",
	}
	return
}
