package chat

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type ChatService struct {
	UnimplementedChatServiceServer
}

// grpcurl usage:
// grpcurl -plaintext -format json -d @ localhost:9000 chat.ChatService/SayHello
// grpcurl -plaintext -import-path . -proto chat.proto localhost:9000 list (could use without reflection.Register)
func (*ChatService) SayHello(ctx context.Context, req *SayRequest) (resp *SayResponse, err error) {
	log.Printf("req: %+v", req)
	resp = &SayResponse{
		Body: "Hakuku",
		Ts:   100001,
	}
	return
}

func SayHi(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		fmt.Fprintf(w, "%s", "hello chat")
		return
	}

	{
		pb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err: %+v", err)
			return
		}
		// log.Printf("pb in: %+v", pb)
		in := &SayRequest{}

		// compress start
		gr, err := gzip.NewReader(bytes.NewReader(pb[5:]))
		if err != nil {
			log.Printf("gzip.NewReader err: %+v", err)
			return
		}
		buf := bytes.Buffer{}
		_, err = buf.ReadFrom(gr)
		if err != nil {
			log.Printf("buf.ReadFrom err: %+v", err)
			return
		}
		result := buf.Bytes()
		// compress end

		err = proto.Unmarshal(result, in)
		if err != nil {
			log.Printf("proto.Unmarshal err: %+v", err)
			return
		}
		log.Printf("in: %+v", in)
	}

	// response
	w.Header().Set("Content-Type", "application/grpc+proto")
	w.Header().Set("Trailer", "grpc-status, grpc-message")
	w.Header().Add("grpc-status", "0")
	w.Header().Add("grpc-message", "")
	log.Printf("req: %+v", r)
	log.Printf("r.Header: %+v", r.Header)

	resp := &SayResponse{
		Body: "Hakuku",
		Ts:   100001,
	}

	pb, err := proto.Marshal(resp)
	if err != nil {
		log.Printf("proto.Marshal err: %+v", err)
		return
	}

	prefix := make([]byte, 5)
	binary.BigEndian.PutUint32(prefix[1:], uint32(len(pb)))
	body := io.MultiReader(bytes.NewReader(prefix), bytes.NewReader(pb))

	bs, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("ioutil.ReadAll err: %+v", err)
		return
	}

	w.Write(bs)
}
