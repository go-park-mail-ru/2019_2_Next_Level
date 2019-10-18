package mailsender

import (
	"2019_2_Next_Level/internal/post"
	"log"
)

var inputChan <-chan post.Email
var outputChan chan<- post.Email
var chansQueue, chansServer post.ChanPair

func Init() {
	log.SetPrefix("MailSender: ")

}
func SetChanPack(chsQueue, chsServer post.ChanPair) {
	chansQueue = chsQueue
	chansServer = chsServer
}

func GetOutputChan() chan<- post.Email {
	return outputChan
}

func Run() {
	ProcessEmail()
}

func ProcessEmail() {
	i := 0
	for pack := range chansQueue.In {
		email := pack.(post.Email)
		log.Println(email.Body)
		// chansServer.Out <- email
		i++

	}
}
