package auth

import "net/http"

type Usecase interface {
	Login(string, string) error
	Logout() error
	Register(string, string) error
	CheckAuthorization(*http.Cookie) error
}
