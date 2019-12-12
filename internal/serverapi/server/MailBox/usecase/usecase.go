package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"2019_2_Next_Level/internal/serverapi/config"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"github.com/microcosm-cc/bluemonday"
	"strconv"
)

type MailBoxUsecase struct {
	repo mailbox.MailRepository
	smtpPort postinterface.IPostInterface
}
var sanitizer *bluemonday.Policy

func NewMailBoxUsecase(repo mailbox.MailRepository, smtp postinterface.IPostInterface) *MailBoxUsecase {
	sanitizer = bluemonday.UGCPolicy()
	usecase := MailBoxUsecase{repo: repo}
	//usecase.smtpPort = postinterface.NewQueueClient(config.Conf.HttpConfig.PostServiceHost, config.Conf.HttpConfig.PostServiceSendPort)
	//usecase.smtpPort.Init()
	usecase.smtpPort = smtp
	return &usecase
}
func (u *MailBoxUsecase) SendMail(email *model.Email) error 	{
	email.From = email.From+"@"+config.Conf.HttpConfig.HostName
	login, host := email.Split(email.To)
	if host==""{
		email.To = login + "@"+config.Conf.HttpConfig.HostName
	}
	postEmail := post.Email{
		From: email.From,
		To:   email.To,
		Body: email.Body,
		Subject:email.Header.Subject,
	}
	if err := u.smtpPort.Put(postEmail); err!=nil{
		return err
	}
	return u.repo.PutSentMessage(*email)
}

func (u *MailBoxUsecase) GetMailList(login string, folder string, sort string, from int, count int) ([]model.Email, error) {
	list, err := u.repo.GetEmailList(login, folder, sort, from, count)
	if err != nil {
		return list, e.Error{}.SetError(err).SetCode(e.ProcessError)
	}
	for i := range list {
		list[i].Sanitize()
	}
	return list, nil
}

func (u *MailBoxUsecase) GetMailListPlain(login string, page int) (int, int, []model.Email, error) {
	mailsPerPage := 25
	count, err := u.repo.GetMessagesCount(login, models.InboxFolder, models.FlagMessageTotal)
	if err != nil {
		return 0, 0, []model.Email{}, err
	}
	from := mailsPerPage*(page-1)+1
	list, err := u.repo.GetEmailList(login, models.InboxFolder, "", from, mailsPerPage)
	for i, _ := range list {
		list[i].Sanitize()
	}
	if err != nil {
		return 0, 0, list, e.Error{}.SetError(err).SetCode(e.ProcessError)
	}
	return count/mailsPerPage+1, page, list, nil
}

func (u *MailBoxUsecase) GetMail(login string, mailID models.MailID) (model.Email, error) {
	id := strconv.Itoa(int(mailID))
	email, err := u.repo.GetEmailByCode(login, id)
	email.Sanitize()
	return email, err
}

func (u *MailBoxUsecase) GetUnreadCount(login string) (int, error) {
	return u.repo.GetMessagesCount(login, models.InboxFolder, models.FlagMessageTotal)
}

func (u *MailBoxUsecase) MarkMail(login string, ids []models.MailID, mark int) error {
	return u.repo.MarkMessages(login, ids, mark)
}

func (u *MailBoxUsecase) AddFolder(login string, foldername string) error {
	return u.repo.AddFolder(login, foldername)
}
func (u *MailBoxUsecase) ChangeMailFolder(login string, foldername string, mailid int64) error {
	return u.repo.ChangeMailFolder(login, foldername, mailid)
}
func (u *MailBoxUsecase) DeleteFolder(login string, foldername string) error {
	return u.repo.DeleteFolder(login, foldername)
}