package messagequeue

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/messagequeue/service"
	"2019_2_Next_Level/pkg/wormhole"
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

// const (
// 	messagequeuePortOutcoming = ":2000"
// 	messagequeuePortIncoming  = ":2001"
// )

// QueueDemon : Инкапсулирует gRPC приёмник и саму очередь, предоставляя интерфейс каналов
type QueueDemon struct {
	queue MessageQueueCore
	chans post.ChanPair
	log   logger.Log
	Name  string
	Task  func()
	Port  string
}

// Init : gets channel pack and inits Queue gRPC service
func (q *QueueDemon) Init(chanA, chanB post.ChanPair) error {
	var t int
	switch q.Name {
	case "incoming":
		q.Port = post.Conf.IncomingQueue.Port
		q.Task = q.Enqueue
		q.chans = chanA
		t = 5
		break
	case "outcoming":
		q.Port = post.Conf.OutcomingQueue.Port
		q.Task = q.Dequeue
		q.chans = chanB
		t = 10
		break
	default:
		q.log.Println("Unknown queue name")
		return fmt.Errorf("unknown name was given: %s\n", q.Name)
	}
	q.queue = MessageQueueCore{Test: t}
	q.queue.Init()
	q.log.SetPrefix(q.Name)
	return nil
}

// Run : starts daemon's work
func (q *QueueDemon) Run(externWg *sync.WaitGroup) {
	defer externWg.Done()
	go q.Task()

	hole := wormhole.Wormhole{}
	err := hole.RunServer(q.Port, func(server *grpc.Server) {
		pb.RegisterMessageQueueServer(server, &q.queue)
	})
	if err != nil {
		fmt.Println("Error after wormhole.runserver()", err)
	}

}

// Dequeue : resends packets from queue to smtp server and prints them
func (q *QueueDemon) Dequeue() {
	i := 0
	for {
		email, err := q.queue.DequeueLocal()
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			q.chans.Out <- email
			q.log.Println(email.Body)
			i++
		}
	}
}

// Enqueue : resends packets from chan into queue and prints them
func (q *QueueDemon) Enqueue() {
	i := 0
	for {
		email := (<-q.chans.In).(post.Email)
		q.log.Println(email.Body)

		data := (&model.ParcelAdapter{}).FromEmail(&email)
		_, err := q.queue.Enqueue(context.Background(), &data)
		fmt.Println("Enqueued")
		if err != nil {
			q.log.Println("Cannot enqueue")
		}
		i++
	}

}
