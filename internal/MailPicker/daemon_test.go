package mailpicker

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// TestNormal: test of picking messages from postQueue and adding to DB
func TestNormal(t *testing.T) {
	mockRepo := MockRepo{}
	mockRepo.Init()
	mockQueue := MockpostInterface{}
	mockQueue.Init()

	mailpicker := Secretary{}
	mailpicker.Init(&mockRepo, &mockQueue)

	input := post.Email{
		"aa@mail.ru",
		"ivan@yandex.ru",
		"Body",
	}
	expected := model.Email(input)
	timer := time.NewTimer(3 * time.Second)

	mockQueue.Put(input)
	go mailpicker.Run(&sync.WaitGroup{})

	select {
	case res := <-mockRepo.bell:
		if !cmp.Equal(res, expected) {
			t.Errorf("Wrong result got: %s", res.Stringify())
		}
		break
	case <-timer.C:
		t.Errorf("Timeout, no result")
	}
	fmt.Println("Done")
}

type MockRepo struct {
	mails []model.Email
	bell  chan model.Email
}

func (d *MockRepo) Init() {
	d.mails = make([]model.Email, 10)
	d.bell = make(chan model.Email, 10)
}
func (d *MockRepo) UserExists(username string) bool {
	return true
}

func (d *MockRepo) AddEmail(email model.Email) error {
	// fmt.Println(email.Stringify())
	d.bell <- email
	return nil
}

type MockpostInterface struct {
	queue chan post.Email
}

func (i *MockpostInterface) Init() {
	i.queue = make(chan post.Email, 100)
}
func (i *MockpostInterface) Destroy() {}
func (i *MockpostInterface) Put(email post.Email) error {
	i.queue <- email
	return nil
}
func (i *MockpostInterface) Get() (post.Email, error) {
	return <-i.queue, nil
}
