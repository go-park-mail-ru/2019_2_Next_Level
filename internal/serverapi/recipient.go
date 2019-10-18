package serverapi

import (
	"2019_2_Next_Level/internal/post"
	"fmt"
	"time"
)

var messageQueue post.Sender

func Run() {
	messageQueue.Init()
	defer messageQueue.Destroy()
	for i := 0; ; i++ {
		e := post.Email{"ivanov@mail.ru", "andrey@yandex.ru",
			fmt.Sprintf("Subject: %d\n\n Hello", i),
		}
		err := messageQueue.Put(e)
		if err != nil {
			fmt.Println("Smth went wrong: ", err)
		} else {
			fmt.Println(e.Stringify())
		}
		time.Sleep(500 * time.Millisecond)
	}

}

func SetQueue(queue post.Sender) {
	messageQueue = &QueueClient{}
}
