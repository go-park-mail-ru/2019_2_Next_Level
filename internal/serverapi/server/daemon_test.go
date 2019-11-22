package server

import (
	"2019_2_Next_Level/internal/serverapi/log"
	"2019_2_Next_Level/tests/mock"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetLogger(&mock.MockLog{})
}

func Test(t *testing.T) {
	wg := &sync.WaitGroup{}
	go Run(wg)
	timer := time.NewTimer(100 * time.Millisecond)
	select {
	case <-timer.C:
		//t.Errorf("Timeout while waiting for ListenAndServe() call")
		break
	}
}
