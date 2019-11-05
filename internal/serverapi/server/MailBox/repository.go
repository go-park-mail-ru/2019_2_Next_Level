package mailbox

import "2019_2_Next_Level/internal/model"

type MailRepository interface {
	GetEmailByCode(code string) (model.Email, error)
	GetEmailList(login string, folder string, sort string, firstNumber int, count int) ([]model.Email, error)
}