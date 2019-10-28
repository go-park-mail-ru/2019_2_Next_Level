package auth

import (
	"2019_2_Next_Level/internal/model"
)

// Usecase: usecase class for Auth module
type Usecase interface {
	CheckAuth(token string) (login string, err error)
	SignIn(login string, password string) (token string, result error)
	SignUp(user model.User) error
	Logout(uuid string) error
	SetRepository(Repository)
}

const (
	ErrorWrongLogin = iota - 1000
	ErrorWrongPassword
	ErrorWrongFirstName
	ErrorWrongFamilyName
	ErrorWrongNickName
	ErrorWrongBirthDate
	ErrorWrongSex
)
