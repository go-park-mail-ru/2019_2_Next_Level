package usecase

import (
	"2019_2_Next_Level/internal/post/log"
)

type MailBoxUsecase struct {
}

func (u *MailBoxUsecase) SendMail(from, to, body string) error {
	log.Log().I("Send mail:\n From: %s\n To: %s\nBody: %s\n", from, to, body)
	return nil
}
