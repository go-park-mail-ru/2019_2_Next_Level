package smtpd

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/smtpd/worker"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	aChan, bChan, server1 := initService()
	if !cmp.Equal(server1.mailSenderChan, aChan) || !cmp.Equal(server1.incomingQueueChan, bChan) {
		t.Errorf("Wrong setting chanPairs")
	}

	// Testing incorrect SMTPServer parameter
	server2 := Server{}
	err := server2.Init(aChan, bChan, struct{}{})
	if err == nil {
		t.Errorf("Wrong handling incorrect smtpServer parameter")
	}

}

func TestRun(t *testing.T) {
	const chanSize = 100
	aChan := post.ChanPair{}.Init(chanSize)
	bChan := post.ChanPair{}.Init(chanSize)

	mock := &MockSMTP{}
	mock.Init("0.0.0.0", ":25")
	server := Server{}
	server.Init(aChan, bChan, mock)
	defer close(server.resultChannel)

	wg := &sync.WaitGroup{}
	go server.Run(wg)
	timer := time.NewTimer(1 * time.Millisecond)
	select {
	case <-mock.resChan:
		break
	case <-timer.C:
		t.Errorf("Timeout while waiting for ListenAndServe() call")
		break
	}
}

func TestIncomingMessage(t *testing.T) {
	_, bChan, server := initService()
	defer close(server.resultChannel)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go server.Run(wg)

	input := worker.EmailNil{
		Email: post.Email{"From", "To", "Data"},
		Error: nil,
	}
	server.resultChannel <- input
	timer := time.NewTimer(1000 * time.Millisecond)
	var res interface{}
	select {
	case res = <-bChan.In:
		break
	case <-timer.C:
		t.Errorf("Timeout while waiting for result in queueChannel")
		return
	}
	output, ok := res.(post.Email)
	if !ok {
		t.Errorf("Wrong type got from server")
		return
	}
	if !cmp.Equal(output, input.Email) {
		fmt.Errorf("Wrong email data got")
	}
}

func TestOutcomingMail(t *testing.T) {
	aChan, bChan, server := initService()
	defer close(server.resultChannel)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go server.Run(wg)

	input := worker.EmailNil{
		Email: post.Email{"From", "To", "Data"},
		Error: nil,
	}
	aChan.Out <- input.Email
	timer := time.NewTimer(1 * time.Millisecond)
	var res interface{}
	select {
	case res = <-bChan.In:
		break
	case <-timer.C:
		t.Errorf("Timeout while waiting for result in chan")
		return
	}

	output, ok := res.(post.Email)
	if !ok {
		t.Errorf("Wrong type got from server")
		return
	}
	if !cmp.Equal(output, input.Email) {
		fmt.Errorf("Wrong email data got")
	}
}

func initService() (post.ChanPair, post.ChanPair, Server) {
	const chanSize = 100
	aChan := post.ChanPair{}.Init(chanSize)
	bChan := post.ChanPair{}.Init(chanSize)

	mock := &MockSMTP{}
	mock.Init("0.0.0.0", ":25")
	mock.Error = fmt.Errorf("Cannot init smtp")

	server := Server{}
	server.Init(aChan, bChan, mock)

	return aChan, bChan, server
}

func WrapServerRun(f func(), resChan chan interface{}) {
	f()
	resChan <- struct{}{}
}

// Check reaction on error during startServerSMTP
func TestErrorOnInitSMTPIncoming(t *testing.T) {
	_, _, server := initService()
	defer close(server.resultChannel)

	wg := &sync.WaitGroup{}
	resChan := make(chan interface{}, 4)

	go WrapServerRun(func() {
		wg.Add(1)
		server.Run(wg)
	}, resChan)

	timer := time.NewTimer(10 * time.Millisecond)
	select {
	case <-timer.C:
		t.Errorf("Timeout waiting result")
		return
	case <-resChan:
		break
	}

}

type MockSMTP struct {
	// smtp.Server
	Domain  string
	Addr    string
	resChan chan interface{}
	backend smtp.Backend
	Error   error
}

func (s *MockSMTP) ListenAndServe() error {
	s.resChan <- struct{}{}
	return s.Error
}
func (s *MockSMTP) Init(port, host string) error {
	s.resChan = make(chan interface{}, 5)
	s.Domain = host
	s.Addr = port
	return s.Error
}
