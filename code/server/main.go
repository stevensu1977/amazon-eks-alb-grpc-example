package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"example.com/chat"

	"google.golang.org/grpc"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"google.golang.org/grpc/reflection"
)

// server is used to implement chat.ChatServiceServer.
type server struct {
	chat.UnimplementedChatServiceServer
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from http handler!\n")
}

// SayHello implements chat.ChatServiceServer
func (s *server) SayHello(ctx context.Context, in *chat.Message) (*chat.Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &chat.Message{Body: "Hello From the Server!"}, nil
}

func main() {

	port := flag.Int("port", 50051, " grpc server port")
	flag.Parse()

	grpcServer := grpc.NewServer()
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", home)

	chat.RegisterChatServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	mixedHandler := newHTTPandGRPCMux(httpMux, grpcServer)
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(mixedHandler, http2Server)}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	fmt.Printf("Go gRPC chat.ChatService start successful, listen on :%d\n", *port)

	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	err = http1Server.Serve(lis)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		panic(err)
	}

}

func newHTTPandGRPCMux(httpHand http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}
