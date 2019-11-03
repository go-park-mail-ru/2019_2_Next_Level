package model

import "back/config"

type UUID interface {
}

type User struct {
	Name       string `json:"firstName"`
	Sirname    string `json:"secondName"`
	MiddleName string
	BirthDate  string `json:"birthDate"`
	Sex        string `json:"sex"`
	Email      string `json:"login"`
	Password   string `json:"password"`
	Avatar     string `json:"avatar"`
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
