package main

import (
	"Go-Grpc-Example/lib/proto"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *api.PingMessage) (*api.PingMessage, error) {
	log.Printf("Serving request: %s", in.Greeting)
	return &api.PingMessage{Greeting: "Hello from server"}, nil
}

func (s *server) MessageSocket(stream api.Ping_MessageSocketServer) error {
	fmt.Println("Starting stream")
	req, err := stream.Recv()

	if err == io.EOF {
		return nil
	} else if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Received %s\n", req.Greeting)
	stream.Send(&api.PingMessage{Greeting: "Hello from server"})

	c := make(chan string)

	go delayedGreeter(c)

	for str := range c {
		stream.Send(&api.PingMessage{Greeting: fmt.Sprintf("Server value: %s", str)})
	}
	return nil
}

func delayedGreeter(c chan string) {
	for i := 0; i < 2; i++ {
		time.Sleep(5 * time.Second)
		c <- fmt.Sprintf("iteration %d", i)
	}
	close(c)
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterPingServer(s, &server{})
	reflection.Register(s)

	log.Print("Started grpc server..")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
