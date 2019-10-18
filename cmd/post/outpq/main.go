package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testBackend/internal/post/outpq"
	pb "testBackend/internal/post/outpq/service"
	"testBackend/internal/serverapi"

	"google.golang.org/grpc"
)

const (
	outpqPort = ":2000"
)

func main() {
	log.SetPrefix("Outpq: ")
	listener, err := net.Listen("tcp", outpqPort)
	if err != nil {
		log.Println("Cannot open tcp socket")
		return
	}

	grpcServer := grpc.NewServer()
	queue := outpq.Outpq{}
	queue.Init()
	pb.RegisterOutpqServer(grpcServer, &queue)

	go Dequeue(&queue)

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println("Cannot start grpcServer")
		return
	}

}

func Dequeue(q *outpq.Outpq) {
	for {
		data, err := q.Dequeue(context.Background(), &pb.Empty{})
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			email := (&serverapi.ParcelAdapter{}).ToEmail(data)
			fmt.Println(email.Stringify())
		}
	}
}
