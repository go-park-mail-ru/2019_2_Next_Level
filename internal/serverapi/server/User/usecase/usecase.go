package usecase

import (
	"2019_2_Next_Level/internal/Auth"
	"2019_2_Next_Level/internal/model"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	"github.com/microcosm-cc/bluemonday"
)

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	sanitizer = bluemonday.UGCPolicy()
	usecase := UserUsecase{repo:repo}
	usecase.auth = &Auth.AuthClient{}
	usecase.auth.Init("0.0.0.0", ":6000")
	return usecase
}

type UserUsecase struct {
	repo user.UserRepository
	auth Auth.IAuthClient
}
var sanitizer *bluemonday.Policy
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
	user.Login = sanitizer.Sanitize(user.Login)
	user.Email = sanitizer.Sanitize(user.Email)
	user.BirthDate = sanitizer.Sanitize(user.BirthDate)
	user.Sex = sanitizer.Sanitize(user.Sex)
	user.Name = sanitizer.Sanitize(user.Name)
	user.Sirname = sanitizer.Sanitize(user.Sirname)
	//user.Avatar = config.Conf.HttpConfig.SelfURL + "avatar/"+user.Avatar
	user.Avatar = "/static/images/icon/no-avatar.svg"
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
