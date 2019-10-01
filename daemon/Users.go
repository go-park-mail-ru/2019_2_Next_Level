package daemon

import db "back/database"

type UserInput struct {
	Name     string
	Email    string
	Password string
}
type UserOutput struct {
	Name       string `json:"name"`
	Sirname    string
	MiddleName string
	Email      string `json:"email"`
	AvaUrl     string
}

func (u *UserOutput) FromUser(dbuser db.User) UserOutput {
	user := UserOutput{
		Name:       dbuser.Name,
		Sirname:    dbuser.Sirname,
		MiddleName: dbuser.MiddleName,
		Email:      dbuser.Email,
	}
	return user
}

func (input *UserInput) ToUser() db.User {
	return db.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
}
