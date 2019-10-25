package auth

import "2019_2_Next_Level/internal/model"

type Repository interface {
	RegisterNewSession(login string, uuid string) error
	CheckSession(uuid string) error
	DiscardSession(uuid string) error
	Registrate(user *model.User) error
	GetUserCredentials(login string) (email string, password string, err error)
}
