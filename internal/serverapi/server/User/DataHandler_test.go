package daemon

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"

	// db "../database"
	db "back/database"
)

func TestGetProfile(t *testing.T) {
	h := DataHandler{}
	db.Init()
	user := db.User{
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
	db.SetUser(user)

	uuid, _ := uuid.NewUUID()
	db.RegisterNewSession(uuid.String(), user.Email)
	r := httptest.NewRequest("GET", "/auth.signup", nil)
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	w := httptest.NewRecorder()
	h.GetProfile(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok", w.Code)
	}
	userGot, err := db.GetUserByEmail(user.Email)
	if err != nil {
		t.Error("Cannot get user from db")
	}
	if !reflect.DeepEqual(userGot, expected) {
		t.Error("Wrond user is written to db!")
	}
}

func TestGetProfileWrongSession(t *testing.T) {
	h := DataHandler{}
	db.Init()

	r := httptest.NewRequest("GET", "/auth.signup", nil)
	uuid, _ := uuid.NewUUID()
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	w := httptest.NewRecorder()
	h.GetProfile(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}

func TestGetProfileEmptySession(t *testing.T) {
	h := DataHandler{}
	db.Init()

	r := httptest.NewRequest("GET", "/auth.signup", nil)
	uuid, _ := uuid.NewUUID()
	db.RegisterNewSession(uuid.String(), "email@yandex.ru")
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	w := httptest.NewRecorder()
	h.GetProfile(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}

func TestGetProfileRottenSession(t *testing.T) {
	h := DataHandler{}
	db.Init()

	r := httptest.NewRequest("GET", "/auth.signup", nil)
	uuid, _ := uuid.NewUUID()
	r.AddCookie(&http.Cookie{Name: "user-token", Value: uuid.String()})

	w := httptest.NewRecorder()
	h.GetProfile(w, r)

	if w.Code != http.StatusBadRequest {
		t.Error("status is not ok", w.Code)
	}
}
