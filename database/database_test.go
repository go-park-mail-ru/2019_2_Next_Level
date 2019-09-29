package database

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestSessions(t *testing.T) {
	Init()
	email := "test@mail.ru"
	uuid, _ := uuid.NewUUID()
	RegisterNewSession(uuid.String(), email)

	newEmail, err := GetUserEmailBySession(uuid.String())
	if err != nil {
		t.Error("Cannot get email by session")
	}
	if newEmail != email {
		t.Error("Got wrong email by session")
	}
}

func TestWrongSessions(t *testing.T) {
	Init()
	email := "test@mail.ru"
	uuid, _ := uuid.NewUUID()
	RegisterNewSession(uuid.String(), email)

	_, err := GetUserEmailBySession(uuid.String() + "err")
	if err == nil {
		t.Error("Got undefined email by session")
	}
}

func TestGetuserByEmail(t *testing.T) {
	user := User{
		Name:     "Ian",
		Email:    "test@mail.ru",
		Password: "12345",
	}
	SetUser(user)

	newUser, err := GetUserByEmail(user.Email)
	if err != nil {
		t.Error("Cannot get user by email")
	}
	if !reflect.DeepEqual(newUser, user) {
		t.Error("Got corrupted user")
	}
}
