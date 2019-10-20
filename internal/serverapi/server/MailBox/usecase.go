package mailbox

type MailBoxUseCase interface {
	SendMail(string, string, string) error
}
