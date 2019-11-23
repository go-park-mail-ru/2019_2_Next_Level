package usecase

import (
	"2019_2_Next_Level/internal/support/models"
	"2019_2_Next_Level/internal/support/repository"
)

type SupportUsecase struct {
	repo repository.Repository
}

func NewSupportUsecase(repo repository.Repository) *SupportUsecase {
	return &SupportUsecase{repo: repo}
}

func (u *SupportUsecase) StartChat(user string, theme string) (chatId interface{}, supporterName string, err error) {
	supporter, _ := u.getSupport()
	res, err := u.repo.StartChat(user, theme, supporter)
	return res, supporter, err
}
func (u *SupportUsecase) GetChat(user string, id interface{}) (models.Chat, error) {
	return u.repo.GetChat(user, id)
}
func (u *SupportUsecase) GetChats(user string) ([]models.Chat, error) {
	return u.repo.GetChats(user)
}
func (u *SupportUsecase) CloseChat(user string, chatId interface{}) error {
	return u.repo.CloseChat(user, chatId)
}
func (u *SupportUsecase) SendMessage(user string, chatId interface{}, message models.Message) error {
	return u.repo.SendMessage(user, chatId, message)
}
func (u *SupportUsecase) GetMessage(user string, chatId interface{}, messageId interface{}) (models.Message, error) {
	return u.repo.GetMessage(user, chatId, messageId)
}

func (u *SupportUsecase) getSupport() (string, error) {
	return "support", nil
}