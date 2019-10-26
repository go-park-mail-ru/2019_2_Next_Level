package model

import "back/config"

type UUID interface {
}

type User struct {
	Name       string
	Sirname    string
	MiddleName string
	BirthDate  string
	Sex        string
	Email      string
	Password   string
	Avatar     string
}

func (user *User) Init() {
	if user.Avatar == "" {
		user.Avatar = config.Configuration.DefaultAvatar
	}
}
func (user *User) Inflate(name, sirname, birth, sex, login, password string) {
	user.Name = name
	user.Sirname = sirname
	user.BirthDate = birth
	user.Sex = sex
	user.Email = login
	user.Password = password
}
