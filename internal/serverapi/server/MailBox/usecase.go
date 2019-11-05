package mailbox

import "2019_2_Next_Level/internal/model"

type MailBoxUseCase interface {
	SendMail(string, string, string) error
	GetMailList(login string, folder string, sort string, from int, count int) ([]model.Email, error)
	GetMail(login string, mailID string) (model.Email, error)
}
