package outpq

import (
	"2019_2_Next_Level/internal/post/outpq"
	pb "2019_2_Next_Level/internal/post/outpq/service"
	"2019_2_Next_Level/internal/serverapi"
	"2019_2_Next_Level/pkg/wormhole"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

const (
	outpqPort = ":2000"
)

// func main() {
func Init() {
	log.SetPrefix("Outpq: ")

	queue := outpq.Outpq{}
	queue.Init()
	go Dequeue(&queue)
	hole := wormhole.Wormhole{}

	err := hole.RunServer(outpqPort, func(server *grpc.Server) {
		pb.RegisterOutpqServer(server, &queue)
	})
	if err != nil {
		fmt.Println("Error after wormhole.runserver()")
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
