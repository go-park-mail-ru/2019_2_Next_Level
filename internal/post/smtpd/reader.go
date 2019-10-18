package smtpd

import (
	"fmt"
	"testBackend/internal/post"
	"testBackend/internal/post/outpq"
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
