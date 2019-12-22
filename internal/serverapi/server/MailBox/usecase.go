package mailbox

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
)

type MailBoxUseCase interface {
	SendMail(email *model.Email) error
	GetMailList(login string, folder string, sort string, since int64, count int) ([]model.Email, error)
	GetMail(login string, mailID []models.MailID) ([]model.Email, error)
	GetUnreadCount(login string) (int, error)
	MarkMail(login string, ids []models.MailID, mark int) error
	AddFolder(login string, foldername string) error
	DeleteFolder(login string, foldername string) error
	ChangeMailFolder(login string, foldername string, mailid []models.MailID) error
	FindMessages(login, request string) ([]int64, error)
}
