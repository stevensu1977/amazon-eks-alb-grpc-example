package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"example.com/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

//buildConnection provide TLS or insecure connection client
func buildConnetion(connectStr string, useTLS bool) (*grpc.ClientConn, error) {
	if useTLS {
		return grpc.Dial(connectStr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	}
	return grpc.Dial(connectStr, grpc.WithTransportCredentials(insecure.NewCredentials()))

}

func main() {

	//parameter parse
	host := flag.String("host", "localhost", "grpc server host")
	port := flag.Int("port", 50051, " grpc server port")
	useTLS := flag.Bool("tls", false, " TLS default disable")
	msg := flag.String("msg", "Hello From Client!", "send message to server ")
	flag.Parse()

	connectStr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Endpoint: %s,  TLS: %v\n", connectStr, *useTLS)

	var conn *grpc.ClientConn
	conn, err := buildConnetion(connectStr, *useTLS)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.SayHello(context.Background(), &chat.Message{Body: *msg})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)

}
