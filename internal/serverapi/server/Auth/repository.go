package auth

import "2019_2_Next_Level/internal/model"

type Repository interface {
	RegisterNewSession(model.User) error
	CheckSession(model.UUID) error
	DiscardSession(string) error
}
