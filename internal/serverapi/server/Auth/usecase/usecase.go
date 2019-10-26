package usecase

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"fmt"

	"github.com/google/uuid"
)

const (
	prefix        = "auth.usecase."
	sessionCookie = "session_id"
)

func NewAuthUsecase(repo auth.Repository) AuthUsecase {
	return AuthUsecase{repo: repo}
}

type AuthUsecase struct {
	repo auth.Repository
}

func (a *AuthUsecase) SetRepository(repo auth.Repository) {
	a.repo = repo
}

func (u *AuthUsecase) CheckAuth(uuid string) error {
	_, err := u.repo.GetLoginBySession(uuid)
	if err != nil {
		return e.Error{}.SetCode(e.NoPermission).SetError(err)
	}
	return nil
}

func (u *AuthUsecase) Logout(uuid string) error {
	err := u.repo.DeleteSession(uuid)
	if err != nil {
		return e.Error{}.SetCode(e.InvalidParams)
	}
	return nil
}

func (u *AuthUsecase) SignUp(user model.User) error {
	// validation
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
	return nil
}

func (u *AuthUsecase) SignIn(login, password string) (string, error) {
	credentials, err := u.repo.GetUserCredentials(login)
	if err != nil {
		return "", e.Error{}.SetCode(e.NotExists)
	}
	if len(credentials) < 2 {
		fmt.Println("Got credentials of len=", len(credentials))
		return "", e.Error{}.SetCode(e.InvalidParams)
	}
	rPass := credentials[0]
	// salt := credentials[1]

	if rPass != password {
		return "", e.Error{}.SetCode(e.InvalidParams)
	}

	uuid, _ := uuid.NewUUID()

	err = u.repo.AddNewSession(login, uuid.String())
	if err != nil {
		return "", e.Error{}.SetCode(e.NotExists)
	}

	return uuid.String(), nil
}

// Errors: WrongLogin, WrongPassword, SessionExists
// func (u *AuthUsecase) Login(user *model.User) (string, e.Error) {
// 	login := user.Email
// 	password := user.Password

// 	Err := e.Error{}.SetPlace(prefix + "Login")
// 	email, pass, err := u.repo.GetUserCredentials(login)
// 	if err != nil {
// 		return "", Err.SetCode(e.WrongLogin).SetError(err).SetPlace("auth.repository")
// 	}
// 	if pass != password {
// 		return "", Err.SetCode(e.WrongPassword)
// 	}
// 	uuid, _ := uuid.NewUUID()
// 	err = u.repo.RegisterNewSession(email, uuid.String())
// 	if err != nil {
// 		return "", Err.SetCode(e.SessionExists).SetPlace("auth.repository")
// 	}
// 	return uuid.String(), e.Ok
// }

// func (u *AuthUsecase) Logout(uuid string) e.Error {
// 	err := u.repo.DiscardSession(uuid)
// 	if err != nil {
// 		return e.Error{}.SetCode(e.InvalidParams).SetError(err)
// 	}
// 	return e.Ok
// }
// func (u *AuthUsecase) Register(data map[string]string) e.Error {
// 	Err := e.Error{}.SetPlace(prefix + ".register")
// 	login, ok1 := data["login"]
// 	password, ok2 := data["password"]
// 	name, ok3 := data["name"]
// 	if !ok1 || !ok2 || !ok3 {
// 		return Err.SetCode(e.IncorrectParams).SetString("Not enought params")
// 	}

// 	user := model.User{Email: login, Password: password, Name: name}
// 	err := u.repo.Registrate(&user)
// 	if err != nil {
// 		return Err.SetCode(e.UserExists).SetError(err)
// 	}

// 	uuid, _ := uuid.NewUUID()
// 	err = u.repo.RegisterNewSession(login, uuid.String())
// 	if err != nil {
// 		return Err.SetCode(e.NoUser).SetError(err)
// 	}

// 	return e.Ok
// }
// func (u *AuthUsecase) CheckAuth(r *http.Request) e.Error {
// 	Err := e.Error{}.SetPlace(prefix + "CheckAuth")
// 	_, err := r.Cookie(sessionCookie)
// 	if err != nil {
// 		return Err.SetCode(e.IncorrectParams).SetError(err)
// 	}
// 	return e.Ok
// }

// func (u *AuthUsecase) GetUser(data []byte) (model.User, e.Error) {
// 	user := model.User{}
// 	if err := json.Unmarshal(data, &user); err != nil {
// 		return user, e.Error{}.SetCode(e.InvalidParams).SetError(err)
// 	}
// 	return user, e.Ok
// }
