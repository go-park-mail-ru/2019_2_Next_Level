package usecase

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
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

func (u *AuthUsecase) CheckAuth(uuid string) (string, error) {
	login, err := u.repo.GetLoginBySession(uuid)
	if err != nil {
		return "", e.Error{}.SetCode(e.NoPermission).SetError(err)
	}
	return login, nil
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
	sault := GenSault(user.Email)
	user.Password = string(PasswordPBKDF2([]byte(user.Password), sault))
	user.Sault = string(sault)
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
		return "", e.Error{}.SetCode(e.InvalidParams)
	}
	rPass := credentials[0]
	salt := credentials[1]

	if !CheckPassword([]byte(password), []byte(rPass), []byte(salt)) {
		return "", e.Error{}.SetCode(e.InvalidParams)
	}

	//if rPass != password {
	//	return "", e.Error{}.SetCode(e.InvalidParams)
	//}

	uuid, _ := uuid.NewUUID()

	err = u.repo.AddNewSession(login, uuid.String())
	if err != nil {
		return "", e.Error{}.SetCode(e.NotExists)
	}

	return uuid.String(), nil
}
