package incomingworker

import (
	"2019_2_Next_Level/internal/post"
	"time"
)

var outChan post.ChanPair

func Run() {
	for {
		email := post.Email{"ivan", "ian", "body"}
		outChan.Out <- email
		time.Sleep(500 * time.Millisecond)
	}
}
