package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/log"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
)

type MailBoxUsecase struct {
	repo mailbox.MailRepository
	smtpPort postinterface.IPostInterface
}

func NewMailBoxUsecase(repo mailbox.MailRepository) *MailBoxUsecase {
	usecase := MailBoxUsecase{repo: repo}
	usecase.smtpPort = postinterface.NewQueueClient(config.Conf.HttpConfig.PostServiceHost, config.Conf.HttpConfig.PostServiceSendPort)
	usecase.smtpPort.Init()
	return &usecase
}

func (u *MailBoxUsecase) SendMail(from, to, body string) error {
	email := post.Email{From: from, To: to, Body: body}
	err := u.smtpPort.Put(email)
	log.Log().I("Send mail:\n From: %s\n To: %s\nBody: %s\n", from, to, body)
	log.Log().I(err)
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