package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	postinterface "2019_2_Next_Level/internal/postInterface"
	"2019_2_Next_Level/internal/serverapi/config"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	e "2019_2_Next_Level/pkg/Error"
	"bytes"
	"github.com/microcosm-cc/bluemonday"
	gomail "gopkg.in/gomail.v2"
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
	//email.From = email.From+"@"+config.Conf.HttpConfig.HostName
	//email.From = email.From+"@"+"nlmail.hldns.ru"
	//email.From = email.From+"@"+"mail.nl-mail.ru"
	login, host := email.Split(email.To)
	if host==""{
		email.To = login + "@"+config.Conf.HttpConfig.HostName
	}
	emailToSend, err := u.PrepareMessage(*email)
	if err != nil {
		return err
	}
	postEmail := post.Email{
		From: emailToSend.From,
		To:   emailToSend.To,
		Body: emailToSend.Body,
		Subject:emailToSend.Header.Subject,
	}
	if err := u.smtpPort.Put(postEmail); err!=nil{
		return err
	}
	return u.repo.PutSentMessage(*email)
}

func (u *MailBoxUsecase) GetMailList(login string, folder string, sort string, since int64, count int) ([]model.Email, error) {
	list, err := u.repo.GetEmailList(login, folder, sort, since, count)
	if err != nil {
		return list, e.Error{}.SetError(err).SetCode(e.ProcessError)
	}
	for i := range list {
		list[i].Sanitize()
	}
	return list, nil
}

func (u *MailBoxUsecase) GetMail(login string, mailID []models.MailID) ([]model.Email, error) {
	id := make([]string, 0, len(mailID))
	for _, elem := range mailID {
		id = append(id, strconv.Itoa(int(elem)))
	}
	email, err := u.repo.GetEmailByCode(login, id)
	for i := range email {
		email[i].Sanitize()
	}
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
func (u *MailBoxUsecase) ChangeMailFolder(login string, foldername string, mailid []models.MailID) error {
	return u.repo.ChangeMailFolder(login, foldername, mailid)
}
func (u *MailBoxUsecase) DeleteFolder(login string, foldername string) error {
	return u.repo.DeleteFolder(login, foldername)
}

func (u *MailBoxUsecase) FindMessages(login, request string) ([]int64, error) {
	return u.repo.FindMessages(login, request)
}

func (u *MailBoxUsecase) PrepareMessage(from model.Email) (*model.Email, error) {
	fromLogin :=from.From
	from.From =from.From+"@"+"mail.nl-mail.ru"
	name, avatar, err := u.repo.GetUserData(fromLogin)
	if err != nil {
		return &from, err
	}
	if avatar=="" {
		avatar = config.Conf.HttpConfig.DefaultAvatar
	}
	new := gomail.NewMessage()
	new.SetHeader("From", from.From, name)
	new.SetHeader("To", from.To)
	new.SetHeader("Subject", from.Header.Subject)
	new.SetBody("text/plain", from.Body)

	path := config.Conf.HttpConfig.RootDir + "/" + config.Conf.HttpConfig.StaticDir
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	path += config.Conf.HttpConfig.AvatarDir+ "/"
	path += avatar
	//new.Attach(path)

	var bodyWriter bytes.Buffer
	new.WriteTo(&bodyWriter)
	from.Body = bodyWriter.String()
	return &from, nil
}