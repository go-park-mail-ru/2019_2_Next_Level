package model

type UUID interface {
}

type User struct {
	Name       string `json:"firstName"`
	Sirname    string `json:"secondName"`
	MiddleName string
	BirthDate  string `json:"birthDate"`
	Login string 	  `json:"nickname"`
	Sex        string `json:"sex"`
	Email      string `json:"login"`
	Password   string `json:"password"`
	Avatar     string `json:"avatar"`
	Sault string `json:-`
}

func (user *User) Init() {
	if user.Avatar == "" {
		user.Avatar = "default.png"
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
