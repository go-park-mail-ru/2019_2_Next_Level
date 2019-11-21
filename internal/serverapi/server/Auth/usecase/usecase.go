package usecase

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/log"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
)

const (
	prefix        = "auth.usecase."
	sessionCookie = "session_id"
)

func NewAuthUsecase(repo auth.Repository) AuthUsecase {
	usecase := AuthUsecase{repo:repo}
	usecase.authService = &Auth.AuthClient{}
	err := usecase.authService.Init("0.0.0.0", ":6000")
	if err != nil {
		log.Log().E("Cannot init auth: ", err)
	}
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
	err := u.authService.GetError(u.authService.CheckCredentials(login, password))
	//credentials, err := u.repo.GetUserCredentials(login)
	if err != nil {
		return "", e.Error{}.SetCode(e.InvalidParams)
	}
	//if len(credentials) < 2 {
	//	return "", e.Error{}.SetCode(e.InvalidParams)
	//}
	//rPass := credentials[0]
	//salt := credentials[1]
	//
	//if !CheckPassword([]byte(password), []byte(rPass), []byte(salt)) {
	//	return "", e.Error{}.SetCode(e.InvalidParams)
	//}

	//if rPass != password {
	//	return "", e.Error{}.SetCode(e.InvalidParams)
	//}

	//uuid, _ := uuid.NewUUID()

	//err = u.repo.AddNewSession(login, uuid.String())
	//if err != nil {
	//	return "", e.Error{}.SetCode(e.NotExists)
	//}
	res, status := u.authService.StartSession(login)

	return res, u.authService.GetError(status)
}
