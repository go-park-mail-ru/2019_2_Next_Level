package repository

import (
	"2019_2_Next_Level/internal/support/models"
)

type Repository interface {
	StartChat(user string, theme string, support string) (chatId interface{}, err error)
	GetChat(user string, id interface{}) (models.Chat, error)
	GetChats(user string) ([]models.Chat, error)
	CloseChat(user string, chatId interface{}) error
	SendMessage(user string, chatId interface{}, message models.Message) error
	GetMessage(user string, chatId interface{}, messageId interface{}) (models.Message, error)
}
