package outpq

import (
	"2019_2_Next_Level/internal/post"
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

type QueueDemon struct {
	queue   outpq.Outpq
	outChan chan post.Email
	chans   post.ChanPair
}

func (q *QueueDemon) SetChanPack(chs post.ChanPair) {
	q.chans = chs
}

func (q *QueueDemon) Init() {
	q.queue = outpq.Outpq{}
	q.queue.Init()
	log.SetPrefix("Outpq: ")
}
func (q *QueueDemon) Run() {
	go q.Dequeue()
	hole := wormhole.Wormhole{}

	err := hole.RunServer(outpqPort, func(server *grpc.Server) {
		pb.RegisterOutpqServer(server, &q.queue)
	})
	if err != nil {
		fmt.Println("Error after wormhole.runserver()", err)
	}

}

func (q *QueueDemon) Dequeue() {
	i := 0
	for {
		data, err := q.queue.Dequeue(context.Background(), &pb.Empty{})
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			email := (&serverapi.ParcelAdapter{}).ToEmail(data)
			q.chans.Out <- email
			fmt.Println(email.Body)
			i++
		}
	}
}
