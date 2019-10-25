package serverapi

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/serverapi/server"
	"sync"
)

var messageQueue post.Sender

// Run : starts server tasks
func Run() {
	incomingMailHandler.Init()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go server.Run(wg)
	go incomingMailHandler.Run(wg)
	wg.Wait()
}

func SetQueue(queue post.Sender) {
	messageQueue = queue
}
