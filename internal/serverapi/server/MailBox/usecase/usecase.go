package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/post/log"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"2019_2_Next_Level/internal/serverapi/config"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	"github.com/microcosm-cc/bluemonday"
)

type MailBoxUsecase struct {
	repo mailbox.MailRepository
	smtpPort postinterface.IPostInterface
}
var sanitizer *bluemonday.Policy

func NewMailBoxUsecase(repo mailbox.MailRepository) *MailBoxUsecase {
	sanitizer = bluemonday.UGCPolicy()
	usecase := MailBoxUsecase{repo: repo}
	usecase.smtpPort = postinterface.NewQueueClient(config.Conf.HttpConfig.PostServiceHost, config.Conf.HttpConfig.PostServiceSendPort)
	usecase.smtpPort.Init()
	return &usecase
}

func (u *MailBoxUsecase) SendMail(from, to, body string) error {
	email := post.Email{
		From: sanitizer.Sanitize(from),
		To: sanitizer.Sanitize(to),
		Body: sanitizer.Sanitize(body),
	}
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
	for i := range list {
		list[i].From = sanitizer.Sanitize(list[i].From)
		list[i].To = sanitizer.Sanitize(list[i].To)
		list[i].Body = sanitizer.Sanitize(list[i].Body)
	}
	return list, nil
}
func (u *MailBoxUsecase) GetMail(login string, mailID string) (model.Email, error) {
	return model.Email{}, nil
}