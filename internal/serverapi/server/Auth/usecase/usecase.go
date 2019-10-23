package usecase

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	"fmt"
	"net/http"
)

type AuthUsecase struct {
	repo auth.Repository
}

func (a *AuthUsecase) SetRepo(с *model.Connection) *AuthUsecase {
	return a
}

func (u *AuthUsecase) Login(login, password string) error {
	return nil
}
func (u *AuthUsecase) Logout() error {
	return nil
}
func (u *AuthUsecase) Register(login, password string) error {
	return nil
}
func (u *AuthUsecase) CheckAuthorization(cookie *http.Cookie) error {
	if cookie.Name != "ian" {
		return fmt.Errorf("Wrong token")
	}
	return nil
}
