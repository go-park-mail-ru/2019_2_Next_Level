package repository

import (
	"2019_2_Next_Level/internal/db"
	"2019_2_Next_Level/internal/model"
)

type IRepository interface {
	SetDB(*db.IDB) error
	GetUser(string) (model.User, error)
	PutUser(model.User) error
	RegisterNewSession(model.User) error
	CheckSession(model.UUID) error
	AddMail(model.Email, model.User) error
	GetMail(model.UUID) (model.Email, error)
	SetMailStatus(model.Email, string) error
}
