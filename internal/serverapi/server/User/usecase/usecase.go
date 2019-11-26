package usecase

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/model"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	e "2019_2_Next_Level/pkg/HttpError/Error"
)

func NewUserUsecase(repo user.UserRepository, auth Auth.IAuthClient) UserUsecase {
	usecase := UserUsecase{repo:repo}
	usecase.auth = auth
	return usecase
}

type UserUsecase struct {
	repo user.UserRepository
	auth Auth.IAuthClient
}
func (u *UserUsecase) GetUser(login string) (model.User, error) {
	user, err := u.repo.GetUser(login)
	if err != nil {
		switch err.(type) {
		case e.Error:
			return user, err
		default:
			break
		}
		return user, e.Error{}.SetCode(e.ProcessError)
	}
	user.Login = user.Email
	//user.Avatar = config.Conf.HttpConfig.SelfURL + "avatar/"+user.Avatar
	user.Avatar = "/static/images/icon/no-avatar.svg"
	user.Sanitize()
	return user, nil
}
func (u *UserUsecase) EditUser(user *model.User) error {
	user.Password = ""
	err := u.repo.UpdateUserData(user)
	if err != nil {
		switch err.(type) {
		case e.Error:
			return err
		default:
			break
		}
		return e.Error{}.SetCode(e.InvalidParams)
	}

	return nil
}
func (u *UserUsecase) EditPassword(login string, oldPass string, newPass string) error {
	return u.auth.GetError(u.auth.ChangePassword(login, oldPass, newPass))
}
