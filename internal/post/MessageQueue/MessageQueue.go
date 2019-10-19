package messagequeue

import (
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/messagequeue/service"
	"2019_2_Next_Level/internal/serverapi"
	"context"
)

const (
	queueSize = 100
)

// MessageQueue : class of the Queue storing emails
type MessageQueue struct {
	queue chan pb.Email
}

// Init : initialize the queue
func (q *MessageQueue) Init() {
	q.queue = make(chan pb.Email, queueSize)
}

// Enqueue : put a grpc Email in the queue
func (q *MessageQueue) Enqueue(ctx context.Context, email *pb.Email) (*pb.Empty, error) {
	q.queue <- *email
	return &pb.Empty{S: true}, nil
}

// EnqueueLocal : put a usual post.Email in the queue
func (q *MessageQueue) EnqueueLocal(email *post.Email) error {
	lEmail := (&serverapi.ParcelAdapter{}).FromEmail(email)
	q.queue <- lEmail
	return nil
}

// Dequeue : get a grpc Email from the queue
func (q *MessageQueue) Dequeue(ctx context.Context, _ *pb.Empty) (*pb.Email, error) {
	email := <-q.queue
	return &email, nil
}

// DequeueLocal : get a usual post.Email from the queue
func (q *MessageQueue) DequeueLocal() (post.Email, error) {
	lEmail := <-q.queue
	email := (&serverapi.ParcelAdapter{}).ToEmail(&lEmail)
	return email, nil
}
