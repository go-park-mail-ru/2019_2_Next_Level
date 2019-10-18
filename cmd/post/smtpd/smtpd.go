package smtpd

import (
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/smtpd/service"
	"context"
	"fmt"
)

type Server struct{}

var chansSender, chansQueue post.ChanPair

func (s *Server) Enqueue(ctx context.Context, in *pb.Email) (*pb.Empty, error) {
	fmt.Printf("New Email.\n From: %s\nTo: %sBody: %s\n", in.GetFrom(), in.GetTo(), in.GetBody())
	return &pb.Empty{S: true}, nil
}

func SetChanPack(chsSender, chsQueue post.ChanPair) {
	chansSender = chsSender
	chansQueue = chsQueue
}

func Init() {
	// func main() {
	fmt.Println("SMTPd started. Hello!")
	// smtpd.Read()
	// listener, err := net.Listen("tcp", ":2000")
	// if err != nil {
	// 	fmt.Println("Cannot open socket")
	// 	return
	// }
	// fmt.Printf("Start listening port %s\n", ":2000")

	// server := grpc.NewServer()

	// pb.RegisterSmtpServerServer(server, &Server{})

	// err = server.Serve(listener)
	// if err != nil {
	// 	fmt.Println("Cannot serve")
	// 	return
	// }

}

func Run() {

}
