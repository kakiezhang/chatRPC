package chat

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func GrpcSay(ctx context.Context) (resp *SayResponse, err error) {
	addr := "127.0.0.1:9000"
	// addr := "127.0.0.1:9001"

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		err = fmt.Errorf("grpc dial err: %v", err)
		return
	}
	defer conn.Close()

	in := &SayRequest{
		Body: "Jack & Rose",
		Id:   200,
	}
	c := NewChatServiceClient(conn)
	resp, err = c.SayHello(ctx, in)
	return
}
