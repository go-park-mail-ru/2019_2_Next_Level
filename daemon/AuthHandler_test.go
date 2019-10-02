package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	db "back/database"

	"github.com/google/uuid"
)

//	Register()
func TestRegister(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}
	expected := db.User{
		Name:       user.Name,
		Sirname:    "",
		MiddleName: "",
		Email:      user.Email,
		Password:   user.Password,
	}

	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signup", body)
	w := httptest.NewRecorder()
	h.Register(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok", w.Code)
	}

	newUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		t.Error("Cannot get new user from db", err)
	}
	if !reflect.DeepEqual(newUser, expected) {
		t.Error("Wrong user is written to db!")
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Error("Wrong cookies count: ", len(cookies))
	}
	cookie := cookies[0]
	userById, err := db.GetUserEmailBySession(cookie.Value)
	if cookie.Name != "user-token" || err != nil || userById != newUser.Email {
		t.Error("Wrong user-token cookie value is set: ", cookie.Name)
	}

}

func TestCurruptedUserRegister(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	type CorruptedUserInput struct {
		Namae      string
		Post       string
		Pass       string
		Errorfield string
	}
	user := CorruptedUserInput{
		Namae:      "Ian",
		Post:       "aa@mail.ru",
		Pass:       "12345",
		Errorfield: "froiOFJF(*&#(jD(*EF",
	}

	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signup", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Error("status is not ok", w.Code)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Error("Wrong cookies count: ", len(cookies))
	}

}
func TestEmptyBodyRegister(t *testing.T) {
	h := AuthHandler{}
	r := httptest.NewRequest("POST", "/auth.signup", nil)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not BadRequest", w.Code)
	}
}

//	Login()
func TestLogin(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}
	db.SetUser(user.ToUser())

	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signin", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok", w.Code)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Error("Wrong cookies count: ", len(cookies))
	}
	cookie := cookies[0]
	userById, err := db.GetUserEmailBySession(cookie.Value)
	if cookie.Name != "user-token" || err != nil || userById != user.Email {
		t.Error("Wrong user-token cookie value is set: ", cookie.Name)
	}
}

func TestAlienLogin(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}

	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signin", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}

func TestEmptyBodyLogin(t *testing.T) {
	h := AuthHandler{}
	r := httptest.NewRequest("POST", "/auth.signin", nil)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}
func TestUnexistedUserLogin(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}
	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signin", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}
func TestWrongPasswordLogin(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}
	db.SetUser(user.ToUser())

	user.Password += "mistake"
	userJson, _ := json.Marshal(user)
	body := bytes.NewReader(userJson)
	r := httptest.NewRequest("POST", "/auth.signin", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}

//  Authorize
func TestAuthorization(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	user := UserInput{
		Name:     "Ian",
		Email:    "aa@mail.ru",
		Password: "12345",
	}
	db.SetUser(user.ToUser())

	uuid, _ := uuid.NewUUID()
	db.RegisterNewSession(uuid.String(), user.Email)

	w := httptest.NewRecorder()

	dbUser := user.ToUser()
	h.Authorize(w, &dbUser)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Error("Wrong cookies count: ", len(cookies))
	}
	cookie := cookies[0]
	userById, err := db.GetUserEmailBySession(cookie.Value)
	if cookie.Name != "user-token" || err != nil || userById != user.Email {
		t.Error("Wrong user-token cookie value is set: ", cookie.Name)
	}
}

// right data
func TestCheckAuthorization(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	email := "user@email.us"
	uuid, _ := uuid.NewUUID()
	db.RegisterNewSession(uuid.String(), email)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	resEmail, err := h.CheckAuthorization(r)
	if err != nil {
		t.Error("Cannot find authorization: ", err)
	}
	if resEmail != email {
		t.Error("Get wrong email: ", resEmail)
	}

}

// user is not authorized
func TestCheckNoAuthorization(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	uuid, _ := uuid.NewUUID()
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	resEmail, err := h.CheckAuthorization(r)
	if err == nil {
		t.Error("Found unexisted authorisation: ", err)
	}
	if resEmail != "" {
		t.Error("Get unempty email: ", resEmail)
	}

}

// wrong user-token ket
func TestWrongCheckAuthorization(t *testing.T) {
	h := AuthHandler{}
	db.Init()
	email := "test@mail.ru"
	uuid, _ := uuid.NewUUID()
	db.RegisterNewSession(uuid.String(), email)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String() + "error"})

	resEmail, err := h.CheckAuthorization(r)
	if err == nil {
		t.Error("Found unexisted authorisation: ", err)
	}
	if resEmail != "" {
		t.Error("Get unempty email: ", resEmail)
	}
}
