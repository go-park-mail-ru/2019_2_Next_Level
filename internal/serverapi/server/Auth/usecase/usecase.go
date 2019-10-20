package usecase

import (
	"fmt"
	"net/http"
)

type AuthUsecase struct {
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
