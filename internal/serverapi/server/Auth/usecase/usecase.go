package usecase

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	prefix        = "auth.usecase."
	sessionCookie = "session_id"
)

func GetUsecase() AuthUsecase {
	return AuthUsecase{}
}

type AuthUsecase struct {
	repo auth.Repository
}

func (a *AuthUsecase) SetRepo(repo auth.Repository) {
	a.repo = repo
}

// Errors: WrongLogin, WrongPassword, SessionExists
func (u *AuthUsecase) Login(user *model.User) (string, e.Error) {
	login := user.Email
	password := user.Password

	Err := e.Error{}.SetPlace(prefix + "Login")
	email, pass, err := u.repo.GetUserCredentials(login)
	if err != nil {
		return "", Err.SetCode(e.WrongLogin).SetError(err).SetPlace("auth.repository")
	}
	if pass != password {
		return "", Err.SetCode(e.WrongPassword)
	}
	uuid, _ := uuid.NewUUID()
	err = u.repo.RegisterNewSession(email, uuid.String())
	if err != nil {
		return "", Err.SetCode(e.SessionExists).SetPlace("auth.repository")
	}
	return uuid.String(), e.Ok
}


func (u *AuthUsecase) Logout(uuid string) e.Error {
	err := u.repo.DiscardSession(uuid)
	if err != nil {
		return e.Error{}.SetCode(e.InvalidParams).SetError(err)
	}
	return e.Ok
}
func (u *AuthUsecase) Register(data map[string]string) e.Error {
	Err := e.Error{}.SetPlace(prefix + ".register")
	login, ok1 := data["login"]
	password, ok2 := data["password"]
	name, ok3 := data["name"]
	if !ok1 || !ok2 || !ok3 {
		return Err.SetCode(e.IncorrectParams).SetString("Not enought params")
	}

	user := model.User{Email: login, Password: password, Name: name}
	err := u.repo.Registrate(&user)
	if err != nil {
		return Err.SetCode(e.UserExists).SetError(err)
	}

	uuid, _ := uuid.NewUUID()
	err = u.repo.RegisterNewSession(login, uuid.String())
	if err != nil {
		return Err.SetCode(e.NoUser).SetError(err)
	}

	return e.Ok
}
func (u *AuthUsecase) CheckAuthorization(r *http.Request) e.Error {
	Err := e.Error{}.SetPlace(prefix + "CheckAuth")
	_, err := r.Cookie(sessionCookie)
	if err != nil {
		return Err.SetCode(e.IncorrectParams).SetError(err)
	}
	return e.Ok
}

func (u *AuthUsecase) GetUser(data []byte) (model.User, e.Error) {
	user := model.User{}
	if err := json.Unmarshal(data, &user); err != nil {
		return user, e.Error{}.SetCode(e.InvalidParams).SetError(err)
	}
	return user, e.Ok
}
