package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	db "back/database"

	"github.com/google/uuid"
)

type AuthHandler struct {
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	handler := &CorsHandler{}
	handler.preflightHandler(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Cannot get body")
		(&Error{ErrorBadRequest}).Send(&w)
		return
	}
	user := UserInput{}
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Error during parse profile", err)
		(&Error{ErrorBadRequest}).Send(&w)
		return
	}

	dbUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println("No such a user")
		(&Error{ErrorBadRequest}).Send(&w)
		return
	}

	if dbUser.Password != user.Password {
		fmt.Println("Wrong pasword")
		(&Error{ErrorBadRequest}).Send(&w)
		return
	}
	fmt.Println("Right user")
	a.Authorize(w, &dbUser)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	(&CorsHandler{}).preflightHandler(w, r)
	email, err := a.CheckAuthorization(r)
	if err != nil {
		fmt.Println(err)
		(&Error{ErrorBadRequest}).Send(&w)
		return
	}
	session, _ := r.Cookie("user-token")
	session.Value = "rotten"
	session.Path = "/"
	session.HttpOnly = true
	db.InvalidateSession(session.Value, email)
	session.Expires = time.Now().AddDate(-1, -1, 0)
	http.SetCookie(w, session)
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	(&CorsHandler{}).preflightHandler(w, r)

	user, err := a.parseUser(r)
	if err != nil {
		fmt.Println(err)
		(&Error{ErrorInternal}).Send(&w)
	}
	db.SetUser(user)
	a.Authorize(w, &user)
}

func (a *AuthHandler) parseUser(r *http.Request) (db.User, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Cannot get body")
		return db.User{}, err
	}

	userInput := UserInput{}
	if err := json.Unmarshal(body, &userInput); err != nil {
		fmt.Println("Error during parse profile", err)
		return db.User{}, err
	}

	user := userInput.ToUser()
	return user, nil
}

func (a *AuthHandler) CheckAuthorization(r *http.Request) (string, error) {
	session, err := r.Cookie("user-token")
	if err != nil {
		return "", errors.New("No cookie")
	}
	fmt.Printf("Cookie: %s\n", session.Value)
	email, err := db.GetUserEmailBySession(session.Value)
	if err != nil {
		return "", errors.New("Wrong session key")
	}
	return email, nil
}

func (a *AuthHandler) Authorize(w http.ResponseWriter, user *db.User) {
	out, _ := uuid.NewUUID()
	db.RegisterNewSession(out.String(), user.Email)
	cookie := http.Cookie{
		Name:     "user-token",
		Value:    out.String(),
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
