package messagequeue

import (
	pb "2019_2_Next_Level/generated/post/MessageQueue/service"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/log"
	"context"
)

const (
	queueSize = 100
)

// MessageQueueCore : class of the Queue storing emails
type MessageQueueCore struct {
	queue chan pb.Email
	Test  int
}

// Init : initialize the queue
func (q *MessageQueueCore) Init() {
	q.queue = make(chan pb.Email, queueSize)
}

// Enqueue : put a grpc Email in the queue
func (q *MessageQueueCore) Enqueue(ctx context.Context, email *pb.Email) (*pb.Empty, error) {
	q.queue <- *email
	return &pb.Empty{S: true}, nil
}

// EnqueueLocal : put a usual post.Email in the queue
func (q *MessageQueueCore) EnqueueLocal(email *post.Email) error {
	lEmail := (&model.ParcelAdapter{}).FromEmail(email)
	q.queue <- lEmail
	return nil
}

// Dequeue : get a grpc Email from the queue
func (q *MessageQueueCore) Dequeue(ctx context.Context, e *pb.Empty) (*pb.Email, error) {
	email := <-q.queue
	log.Log().L("dequeued")
	return &email, nil
}

// DequeueLocal : get a usual post.Email from the queue
func (q *MessageQueueCore) DequeueLocal() (post.Email, error) {
	lEmail := <-q.queue
	email := (&model.ParcelAdapter{}).ToEmail(&lEmail)
	return email, nil
}
