package usecase

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/pkg/Error"
)

const (
	prefix        = "auth.usecase."
	sessionCookie = "session_id"
)

func NewAuthUsecase(repo auth.Repository, authService Auth.IAuthClient) AuthUsecase {
	usecase := AuthUsecase{repo:repo}
	usecase.authService = authService
	return usecase
}

type AuthUsecase struct {
	repo auth.Repository
	authService Auth.IAuthClient
}

func (a *AuthUsecase) SetRepository(repo auth.Repository) {
	a.repo = repo
}

func (u *AuthUsecase) CheckAuth(uuid string) (string, error) {
	login, status := u.authService.LoginBySession(uuid)
	return login, u.authService.GetError(status)
}

func (u *AuthUsecase) Logout(uuid string) error {
	return u.authService.GetError(u.authService.DestroySession(uuid))
}

func (u *AuthUsecase) SignUp(user model.User) error {
	err := u.repo.AddNewUser(&user)
	if err != nil {
		switch err.(type) {
		case e.Error:
			switch err.(e.Error).Code {
			case e.AlreadyExists:
				return err
			case auth.ErrorWrongBirthDate, auth.ErrorWrongFamilyName, auth.ErrorWrongFirstName,
				auth.ErrorWrongLogin, auth.ErrorWrongNickName, auth.ErrorWrongPassword, auth.ErrorWrongSex:
				return e.Error{}.SetCode(e.InvalidParams).SetError(err)
			default:
				return e.Error{}.SetCode(e.InvalidParams)
			}
			break
		default:
			return e.Error{}.SetCode(e.InvalidParams)
		}
	}
	_ = u.authService.RegisterUser(user.Email, user.Password)
	return nil
}

func (u *AuthUsecase) SignIn(login, password string) (string, error) {
	err := u.authService.GetError(u.authService.CheckCredentials(login, password))
	if err != nil {
		return "", e.Error{}.SetCode(e.InvalidParams)
	}
	res, status := u.authService.StartSession(login)

	return res, u.authService.GetError(status)
}
