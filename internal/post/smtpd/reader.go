package smtpd

import (
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/outpq"
	"fmt"
)

func Read() {
	q := outpq.GetInstance()
	if q == nil {
		fmt.Println("Nil queue")
		return
	}
	for {
		elem := q.Dequeue()
		mail := elem.(post.Email)
		fmt.Println(mail.Stringify())
	}
}

func Run() {
	
}
