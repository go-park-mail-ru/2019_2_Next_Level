package main

import (
	"context"
	"fmt"
	"net"
	pb "testBackend/internal/post/smtpd/service"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) Enqueue(ctx context.Context, in *pb.Email) (*pb.Empty, error) {
	fmt.Printf("New Email.\n From: %s\nTo: %sBody: %s\n", in.GetFrom(), in.GetTo(), in.GetBody())
	return &pb.Empty{S: true}, nil
}

func main() {
	fmt.Println("SMTPd started. Hello!")
	// smtpd.Read()
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Println("Cannot open socket")
		return
	}
	fmt.Printf("Start listening port %s\n", ":2000")

	server := grpc.NewServer()

	pb.RegisterSmtpServerServer(server, &Server{})

	err = server.Serve(listener)
	if err != nil {
		fmt.Println("Cannot serve")
		return
	}

}
