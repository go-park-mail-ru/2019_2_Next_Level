package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post/log"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
)

type MailBoxUsecase struct {
	repo mailbox.MailRepository
}

func NewMailBoxUsecase(repo mailbox.MailRepository) *MailBoxUsecase {
	return &MailBoxUsecase{repo: repo}
}

func (u *MailBoxUsecase) SendMail(from, to, body string) error {
	log.Log().I("Send mail:\n From: %s\n To: %s\nBody: %s\n", from, to, body)
	return nil
}

func (u *MailBoxUsecase) GetMailList(login string, folder string, sort string, from int, count int) ([]model.Email, error) {
	list, err := u.repo.GetEmailList(login, folder, sort, from, count)
	if err != nil {
		return list, e.Error{}.SetError(err).SetCode(e.ProcessError)
	}
	return list, nil
}
func (u *MailBoxUsecase) GetMail(login string, mailID string) (model.Email, error) {
	return model.Email{}, nil
}