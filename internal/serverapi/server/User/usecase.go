package user

import (
	"2019_2_Next_Level/internal/model"
	"mime/multipart"
)

type UserUsecase interface {
	GetUser(login string) (model.User, error)
	GetUserFolders(login string) ([]model.Folder, error)
	EditAvatar(login string, file multipart.File, header *multipart.FileHeader) (string, error)
	EditUser(*model.User) error
	EditPassword(login string, currPass string, newPass string) error
}
