package outpq

import (
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/outpq/service"
	"2019_2_Next_Level/internal/serverapi"
	"context"
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

func (q *Outpq) EnqueueLocal(email *post.Email) error {
	lEmail := (&serverapi.ParcelAdapter{}).FromEmail(email)
	q.queue <- lEmail
	return nil
}

func (q *Outpq) Dequeue(ctx context.Context, _ *pb.Empty) (*pb.Email, error) {
	email := <-q.queue
	// ACHTUNG! Скорее всего не сработает
	return &email, nil
}

func (q *Outpq) DequeueLocal() (post.Email, error) {
	lEmail := <-q.queue
	email := (&serverapi.ParcelAdapter{}).ToEmail(&lEmail)
	return email, nil
}

func (q *Outpq) GetChan() chan pb.Email {
	return q.queue
}
