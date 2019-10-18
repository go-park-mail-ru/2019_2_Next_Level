package main

import (
	mailsender "2019_2_Next_Level/cmd/post/MailSender"
	"2019_2_Next_Level/cmd/post/outpq"
	"2019_2_Next_Level/cmd/post/smtpd"
	"2019_2_Next_Level/internal/post"
	"sync"
)

const (
	chanSize = 100
)

type Dispatcher struct {
	OutcomingQueue    post.ChanPair
	IncomingQueue     post.ChanPair
	MailSender_Queue  post.ChanPair
	MailSender_Server post.ChanPair
	SMTPServer_Sender post.ChanPair
	SMTPServer_Queue  post.ChanPair
}

func (d *Dispatcher) Init() {
	d.OutcomingQueue.Out = make(chan interface{}, chanSize)

	d.MailSender_Queue.In = d.OutcomingQueue.Out
	d.MailSender_Queue.Out = make(chan interface{}, chanSize)

	d.MailSender_Server.In = make(chan interface{}, chanSize)
	d.MailSender_Server.Out = make(chan interface{}, chanSize)

	d.OutcomingQueue.In = d.MailSender_Queue.Out

	d.SMTPServer_Sender.In = d.MailSender_Server.Out
	d.SMTPServer_Sender.Out = d.MailSender_Server.In

	d.SMTPServer_Queue.Out = make(chan interface{}, chanSize)
	d.IncomingQueue.In = d.SMTPServer_Queue.Out

}

var d Dispatcher

func main() {
	// Должны быть компоненты:
	// 	* Очередь исходящих
	// 	* Очередь Входящих
	// 	* Отправщик
	// 	* SMTP-сервер
	outcomingQueue := outpq.QueueDemon{}
	// incomingQueue := outpq.QueueDemon{}
	outcomingQueue.Init()
	mailsender.Init()
	smtpd.Init()
	d.Init()

	outcomingQueue.SetChanPack(d.OutcomingQueue)
	mailsender.SetChanPack(d.MailSender_Queue, d.MailSender_Server)
	smtpd.SetChanPack(d.SMTPServer_Queue, d.SMTPServer_Sender)
	// incomingQueue.SetChanPack(d.IncomingQueue)

	wg := &sync.WaitGroup{}
	wg.Add(4)
	go outcomingQueue.Run()
	go mailsender.Run()
	go smtpd.Run()
	// go incomingQueue.Run()
	wg.Wait()

}
