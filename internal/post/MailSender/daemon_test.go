package mailsender

import (
	"2019_2_Next_Level/internal/post"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	const chanSize = 100

	aChan := post.ChanPair{}.Init(chanSize)
	bChan := post.ChanPair{}.Init(chanSize)

	sender := MailSender{}
	sender.Init(aChan, bChan)

	// if &sender.queueChan == &aChan || &sender.smtpChan == &bChan {
	if !cmp.Equal(sender.queueChan, aChan) || !cmp.Equal(sender.smtpChan, bChan) {
		t.Errorf("Wrong setting chanPairs")
	}
}

// Checks if worker waits for email and ends after closing the channel
func TestRun(t *testing.T) {
	const chanSize = 100

	aChan := post.ChanPair{}.Init(chanSize)
	bChan := post.ChanPair{}.Init(chanSize)

	sender := MailSender{}
	sender.Init(aChan, bChan)
	wg := &sync.WaitGroup{}

	resChan := make(chan interface{})
	var f func()
	f = func() {
		wg.Add(1)
		sender.Run(wg)
		resChan <- struct{}{}
	}
	timer := time.NewTimer(10 * time.Millisecond)
	go f()
	select {
	case <-timer.C:
		fmt.Println("timer")
		break
	case <-resChan:
		t.Errorf("Early returning of the function")
		return
	}

	close(aChan.Out)
	timer2 := time.NewTimer(10 * time.Millisecond)
	select {
	case <-timer2.C:
		t.Errorf("Function does not exiting on closing channel")
		return
	case <-resChan:
		return
	}
}

func TestProcessEmail(t *testing.T) {
	const chanSize = 100

	aChan := post.ChanPair{}.Init(chanSize)
	bChan := post.ChanPair{}.Init(chanSize)

	sender := MailSender{}
	sender.Init(aChan, bChan)
	wg := &sync.WaitGroup{}
	go sender.Run(wg)

	inputEmail := post.Email{From: "From tester", To: "To proger", Body: "Die, please"}
	aChan.Out <- inputEmail
	timer := time.NewTimer(1 * time.Millisecond)
	var res interface{}
	select {
	case <-timer.C:
		t.Errorf("Timeout while waiting an answer")
		return
	case res = <-bChan.Out:
		break
	}

	email, ok := res.(post.Email)
	if !ok {
		t.Errorf("Wrong type returned")
	}
	expected := inputEmail
	if !cmp.Equal(expected, email) {
		t.Errorf("Wrong answer got")
	}
}
