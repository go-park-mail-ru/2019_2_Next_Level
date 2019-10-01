package database

import "errors"

type User struct {
	Name       string
	Sirname    string
	MiddleName string
	Email      string
	Password   string `json:"-"`
	Avatar     string
}

var sessionBook map[string]string

var usersList map[string]User

func GetUserEmailBySession(sessionId string) (string, error) {
	email, err := sessionBook[sessionId]
	if !err {
		return email, errors.New("No element")
	}
	return email, nil
}

func RegisterNewSession(sessionId, email string) {
	sessionBook[sessionId] = email
}

func InvalidateSession(sessionId, email string) {
	if sessionBook[sessionId] == email {
		delete(sessionBook, sessionId)
	}
}

func GetUserByEmail(email string) (User, error) {
	user, err := usersList[email]
	if !err {
		return User{}, errors.New("No such a user")
	}
	return user, nil
}

func SetUser(u User) {
	usersList[u.Email] = u
	u.Email = ""
}

func UpdateUser(u User) {
	SetUser(u)
}

func GetAvaFilename(u User) string {
	return u.Email + ".png"
}
func Init() {
	sessionBook = map[string]string{
		"12345": "a@mail.ru",
	}
	usersList = map[string]User{
		"aa@mail.ru": User{
			Name:       "Ian",
			Sirname:    "Ivanov",
			MiddleName: "tamerlanchik",
			Email:      "aa@mail.ru",
			Password:   "pass",
			Avatar:     "626d57a727d65725a21a891ed278810f.jpg",
		},
	}
}
