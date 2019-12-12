package mailbox

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
)

type MailBoxUseCase interface {
	SendMail(email *model.Email) error
	GetMailList(login string, folder string, sort string, from int, count int) ([]model.Email, error)
	GetMailListPlain(login string, page int) (pageCount int, pageReal int, mails []model.Email, err error)
	GetMail(login string, mailID models.MailID) (model.Email, error)
	GetUnreadCount(login string) (int, error)
	MarkMail(login string, ids []models.MailID, mark int) error
	AddFolder(login string, foldername string) error
	DeleteFolder(login string, foldername string) error
	ChangeMailFolder(login string, foldername string, mailid int64) error
}
