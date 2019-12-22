package mailbox

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
)

type MailRepository interface {
	GetEmailByCode(login string, code interface{}) ([]model.Email, error)
	GetEmailList(login string, folder string, sort interface{}, since int64, count int) ([]model.Email, error)
	GetMessagesCount(login string, folder string, flag interface{}) (int, error)
	MarkMessages(login string, messagesID []models.MailID, mark interface{}) error
	PutSentMessage(email model.Email) error
	AddFolder(login string, foldername string) error
	DeleteFolder(login string, foldername string) error
	ChangeMailFolder(login string, foldername string, mailid []models.MailID) error
	FindMessages(login, request string) ([]int64, error)
}