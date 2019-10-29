package user

import "2019_2_Next_Level/internal/model"

type UserUsecase interface {
	GetUser(login string) (model.User, error)
	EditUser(*model.User) error
	EditPassword(login string, currPass string, newPass string) error
}
