package main

import (
	"Go-Grpc-Example/lib/proto"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	simpleRPC(conn)
	bidirectionalStream(conn)
}

func bidirectionalStream(conn *grpc.ClientConn) {
	fmt.Println("=========  Running Bi-directional Stream =========================")
	c := api.NewPingClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	stream, err := c.MessageSocket(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer stream.CloseSend()

	waitc := make(chan struct{})
	go func() {
		for {
			fmt.Println("Sending msg server..")
			stream.Send(&api.PingMessage{Greeting: "Handshake"})
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("EOF exiting..")
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("--Received msg: %s", in.Greeting)
		}
	}()

	<-waitc
	fmt.Println("Exiting...")
}

func simpleRPC(conn *grpc.ClientConn) {
	fmt.Println("=========  Running simple RPC =========================")
	c := api.NewPingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &api.PingMessage{Greeting: "Hello from client"})

	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}

	log.Printf("Response from server %s", r.Greeting)
}
