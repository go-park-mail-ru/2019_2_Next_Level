package auth

import "2019_2_Next_Level/internal/model"

type Repository interface {
	GetLoginBySession(uuid string) (login string, err error)
	AddNewSession(login string, uuid string) error
	DeleteSession(uuiв string) error
	AddNewUser(*model.User) error
	GetUserCredentials(login string) ([]string, error)
}
