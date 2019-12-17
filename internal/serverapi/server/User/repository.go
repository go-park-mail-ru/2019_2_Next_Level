package user

import "2019_2_Next_Level/internal/model"

type UserRepository interface {
	GetUser(login string) (model.User, error)
	GetUserFolders(login string) ([]model.Folder, error)
	UpdateUserData(*model.User) error
}
