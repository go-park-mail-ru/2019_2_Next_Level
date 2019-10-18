package outpq

import (
	"context"
	pb "testBackend/internal/post/outpq/service"
)

const (
	queueSize = 100
)

type Outpq struct {
	queue chan pb.Email
}

func (q *Outpq) Init() {
	q.queue = make(chan pb.Email, queueSize)
}

func (q *Outpq) Enqueue(ctx context.Context, email *pb.Email) (*pb.Empty, error) {
	q.queue <- *email
	return &pb.Empty{S: true}, nil
}

func (q *Outpq) Dequeue(ctx context.Context, _ *pb.Empty) (*pb.Email, error) {
	email := <-q.queue
	// ACHTUNG! Скорее всего не сработает
	return &email, nil
}