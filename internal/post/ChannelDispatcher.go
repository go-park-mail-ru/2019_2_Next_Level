package post

type Dispatcher struct {
	OutcomingQueue    ChanPair
	IncomingQueue     ChanPair
	MailSender_Queue  ChanPair
	MailSender_Server ChanPair
	SMTPServer_Sender ChanPair
	SMTPServer_Queue  ChanPair
}

func (d *Dispatcher) Init(chanSize int) {
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
