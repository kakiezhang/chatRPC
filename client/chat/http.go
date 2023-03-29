package chat

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

func HttpSay(ctx context.Context) (resp *SayResponse, err error) {
	api := "http://127.0.0.1:9000/chat.ChatService/SayHello"
	// api := "http://127.0.0.1:9001/chat.ChatService/SayHello"
	req := &SayRequest{
		Body: "Jack & Rose",
		Id:   200,
	}
	// printIn(req)

	resp = &SayResponse{}

	_, err = DoUnary(ctx, api, req, resp)
	if err != nil {
		return
	}

	return
}

func printIn(req proto.Message) {
	pb1, _ := proto.Marshal(req)
	fmt.Printf("pb1 in: %+v\n", pb1)

	var r bytes.Buffer
	gz := gzip.NewWriter(&r)
	_, _ = gz.Write(pb1)
	// gz.Flush()
	gz.Close()
	pb := r.Bytes()
	fmt.Printf("gzip pb1 in: %+v\n", pb)
}

func DoUnary(ctx context.Context, api string, req, resp proto.Message) (h2resp *http.Response, err error) {
	pb1, err := proto.Marshal(req)
	if err != nil {
		return
	}

	// compress start

	var r bytes.Buffer
	gz := gzip.NewWriter(&r)
	_, _ = gz.Write(pb1)
	// gz.Flush()
	gz.Close()
	pb := r.Bytes()

	// r := bytes.NewBuffer(pb1)
	// w, err := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
	// w.Reset(r)
	// w.Flush()
	// w.Close()
	// pb := r.Bytes()

	// compress end

	prefix := make([]byte, 5)
	binary.BigEndian.PutUint32(prefix[1:], uint32(len(pb)))

	prefix[0] = 1
	log.Printf("prefix: %+v, pb: %+v", prefix, pb)

	body := io.MultiReader(bytes.NewReader(prefix), bytes.NewReader(pb))

	h2req, err := http.NewRequest("POST", api, body)
	if err != nil {
		return
	}

	h2req = h2req.WithContext(ctx)

	h2req.Header.Set("te", "trailers")
	h2req.Header.Set("content-type", "application/grpc+proto")
	h2req.Header.Set("caller", "Chat-Client-0.1")

	// compress gzip
	h2req.Header.Set("grpc-encoding", "gzip")

	h2req = h2req.WithContext(ctx)

	start := time.Now()

	// http.ProxyURL(fixedURL * url.URL)

	// t1 := http.DefaultTransport.(*http.Transport)
	// t1.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// t2, err := http2.ConfigureTransports(t1)
	// if err != nil {
	// 	return
	// }

	var t2 = &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	c := &http.Client{
		Timeout:   3600 * time.Second,
		Transport: t2,
	}
	h2resp, err = c.Do(h2req)

	d := time.Since(start)

	url := fmt.Sprintf("%s%s", h2req.URL.Host, h2req.URL.Path)

	status := http.StatusOK
	if err != nil {
		err = errors.Wrap(err, "")
		status = http.StatusGatewayTimeout
		return
	} else {
		status = h2resp.StatusCode
	}

	log.Printf(
		"[HTTP] method:%s url:%s status:%d query:%s cost:%.5fs",
		h2req.Method,
		url,
		status,
		h2req.URL.RawQuery,
		d.Seconds(),
	)

	defer h2resp.Body.Close()

	pb, err = ioutil.ReadAll(h2resp.Body)
	if err != nil {
		return
	}

	log.Printf(
		"[HTTP] resp.Headers: %+v",
		h2resp.Header,
	)

	if status := h2resp.Trailer.Get("grpc-status"); status != "0" {
		err = errors.Errorf(`
trailer grpc status: [%s]
trailer grpc message: [%s]
header grpc status: [%s]
header grpc message: [%s]
`,
			status,
			h2resp.Trailer.Get("grpc-message"),
			h2resp.Header.Get("grpc-status"),
			h2resp.Header.Get("grpc-message"))
		return
	}

	// if ungzip
	if true {
		gr, err := gzip.NewReader(bytes.NewReader(pb[5:]))
		if err != nil {
			// log.Printf("gzip.NewReader err: %+v", err)
			return nil, err
		}

		buf := bytes.Buffer{}
		_, err = buf.ReadFrom(gr)
		if err != nil {
			// log.Printf("buf.ReadFrom err: %+v", err)
			return nil, err
		}
		b := buf.Bytes()

		log.Println("hehehehe")

		err = proto.Unmarshal(b, resp)
	} else {
		err = proto.Unmarshal(pb[5:], resp)
	}

	return
}
