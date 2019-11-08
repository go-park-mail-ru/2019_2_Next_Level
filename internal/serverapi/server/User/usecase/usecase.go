package usecase

import (
	"2019_2_Next_Level/internal/model"
	authusecase "2019_2_Next_Level/internal/serverapi/server/Auth/usecase"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	"github.com/microcosm-cc/bluemonday"
)

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	sanitizer = bluemonday.UGCPolicy()
	return UserUsecase{repo: repo}
}

type UserUsecase struct {
	repo user.UserRepository
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
	user.Email = sanitizer.Sanitize(user.Email)
	user.BirthDate = sanitizer.Sanitize(user.BirthDate)
	user.Sex = sanitizer.Sanitize(user.Sex)
	user.Name = sanitizer.Sanitize(user.Name)
	user.Sirname = sanitizer.Sanitize(user.Sirname)
	//user.Avatar = config.Conf.HttpConfig.SelfURL + "avatar/"+user.Avatar
	user.Avatar = "/images/icons/no-avatar.svg"
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
	currPass, sault, err := u.repo.GetUserCredentials(login)
	if err != nil {
		return err
	}

	if !authusecase.CheckPassword([]byte(oldPass), []byte(currPass), []byte(sault)) {
		return e.Error{}.SetCode(e.Wrong)
	}
	newPassHash := authusecase.PasswordPBKDF2([]byte(newPass), []byte(sault))

	err = u.repo.UpdateUserPassword(login, string(newPassHash), sault)
	return err
}
