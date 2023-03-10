package chat

import "context"

type ChatService struct {
	UnimplementedChatServiceServer
}

func (*ChatService) SayHello(ctx context.Context, req *Message) (resp *Message, err error) {
	resp = &Message{
		Body: "Hakuku",
	}
	return
}
