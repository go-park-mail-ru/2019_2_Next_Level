package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	db "../database"
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
		http.Error(w, err.Error(), 500)
		return
	}
	user := UserInput{}
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Error during parse profile", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println("No such a user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dbUser.Password != user.Password {
		fmt.Println("Wrong pasword")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Right user")
	a.Authorize(&w, &dbUser)
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	handler := &CorsHandler{}
	handler.preflightHandler(w, r)

	user, err := a.parseUser(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	db.SetUser(user)
	a.Authorize(&w, &user)
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
	session, err := r.Cookie("user-id")
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

func (a *AuthHandler) Authorize(w *http.ResponseWriter, user *db.User) {
	out, _ := uuid.NewUUID()
	db.RegisterNewSession(out.String(), user.Email)
	cookie := http.Cookie{
		Name:    "user-id",
		Value:   out.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(*w, &cookie)
}
