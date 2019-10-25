package auth

import (
	"2019_2_Next_Level/internal/model"
)

// Usecase: usecase class for Auth module
type Usecase interface {
	CheckAuth(token string) error
	SignIn(login string, password string) (token string, result error)
	SignUp(user model.User) error
	Logout(uuid string) error
	SetRepository(Repository)
}

const (
	ErrorWrongLogin = iota
	ErrorWrongPassword
	ErrorWrongFirstName
	ErrorWrongFamilyName
	ErrorWrongNickName
	ErrorWrongBirthDate
	ErrorWrongSex
)
