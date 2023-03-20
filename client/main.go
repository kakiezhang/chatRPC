package main

import (
	"context"
	"fmt"

	"chatRPC/client/chat"
)

func main() {
	ctx := context.Background()
	resp, err := chat.HttpSay(ctx)
	// resp, err := chat.GrpcSay(ctx)
	fmt.Printf("err: %+v\n", err)
	fmt.Printf("resp: %+v\n", resp)
}
