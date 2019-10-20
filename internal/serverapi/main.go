package serverapi

import (
	"2019_2_Next_Level/internal/post"
	incommail "2019_2_Next_Level/internal/serverapi/IncomingMailSecretary"
	"2019_2_Next_Level/internal/serverapi/server"
	"sync"
)

var messageQueue post.Sender

func Run() {
	incomingMailHandler := incommail.Secretary{}
	incomingMailHandler.Init()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go server.Run(wg)
	go incomingMailHandler.Run(wg)

	// messageQueue.Init()
	// defer messageQueue.Destroy()
	// for i := 0; ; i++ {
	// 	e := post.Email{"ivanov@mail.ru", "andrey@yandex.ru",
	// 		fmt.Sprintf("Subject: %d\n\n Hello", i),
	// 	}
	// 	err := messageQueue.Put(e)
	// 	if err != nil {
	// 		fmt.Println("Smth went wrong: ", err)
	// 	} else {
	// 		fmt.Println(e.Stringify())
	// 	}
	// 	time.Sleep(500 * time.Millisecond)
	// }
	wg.Wait()

}

func SetQueue(queue post.Sender) {
	messageQueue = queue
}
