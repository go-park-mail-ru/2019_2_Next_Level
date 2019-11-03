package usecase

import (
	"2019_2_Next_Level/internal/model"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
)

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	return UserUsecase{repo: repo}
}

type UserUsecase struct {
	repo user.UserRepository
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
	currPass, _, err := u.repo.GetUserCredentials(login)
	if err != nil {
		return err
	}
	// generate pass
	if currPass != oldPass {
		return e.Error{}.SetCode(e.Wrong)
	}

	newSault := "sault"
	err = u.repo.UpdateUserPassword(login, newPass, newSault)
	return err
}
