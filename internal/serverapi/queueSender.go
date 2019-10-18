package serverapi

import (
	"context"
	"fmt"
	"testBackend/internal/post"
	pb "testBackend/internal/post/outpq/service"
	"time"

	"google.golang.org/grpc"
)

type Queue struct {
	queue      pb.OutpqClient
	Connection *grpc.ClientConn
}

func (q *Queue) Init() {
	var err error
	q.Connection, err = grpc.Dial("localhost:2000", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot connect to service")
		return
	}

	q.queue = pb.NewOutpqClient(q.Connection)
}

func (q *Queue) Destroy() {
	q.Connection.Close()
}

func (q *Queue) Put(email post.Email) error {
	p := (&ParcelAdapter{}).FromEmail(&email)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := q.queue.Enqueue(ctx, &p)
	return err
}

func (q *Queue) Get() (post.Email, error) {
	data, err := q.queue.Dequeue(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println("Nil value")
		return post.Email{}, err
	}

	return (&ParcelAdapter{}).ToEmail(data), nil
}

type ParcelAdapter struct {
}

func (a *ParcelAdapter) ToEmail(from *pb.Email) post.Email {
	return post.Email{from.From, from.To, from.Body}
}

func (a *ParcelAdapter) FromEmail(from *post.Email) pb.Email {
	return pb.Email{
		From: from.From,
		To:   from.To,
		Body: from.Body,
	}
}
