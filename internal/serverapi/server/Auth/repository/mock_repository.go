package repository

import (
	"2019_2_Next_Level/internal/model"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"fmt"
)

type MockRepository struct {
	sessionBook map[string]string
	usersList   map[string]model.User
}

func GetMock() MockRepository {
	r := MockRepository{}
	r.sessionBook = make(map[string]string)
	r.usersList = make(map[string]model.User)
	return r
}

func (r *MockRepository) AddNewSession(login, uuid string) error {
	if err := r.checkUserExist(login); err != nil {
		return err
	}
	r.sessionBook[uuid] = login
	return nil
}

func (r *MockRepository) GetLoginBySession(uuid string) (string, error) {
	login, ok := r.sessionBook[uuid]
	if !ok {
		return "", e.Error{}.SetCode(e.NoPermission)
	}
	return login, nil
}

func (r *MockRepository) DeleteSession(id string) error {
	if _, ok := r.sessionBook[id]; !ok {
		return e.Error{}.SetCode(e.NotExists)
	}
	delete(r.sessionBook, id)
	return nil
}

func (r *MockRepository) AddNewUser(user *model.User) error {
	if _, ok := r.usersList[user.Email]; ok {
		return e.Error{}.SetCode(e.AlreadyExists)
	}
	r.usersList[user.Email] = *user
	return nil
}

func (r *MockRepository) GetUserCredentials(login string) ([]string, error) {
	user, ok := r.usersList[login]
	if !ok {
		return []string{}, e.Error{}.SetCode(e.NotExists)
	}
	return []string{user.Password, ""}, nil
}

func (r *MockRepository) checkUserExist(login string) error {
	_, err := r.usersList[login]
	if err {
		return nil
	}
	return fmt.Errorf("No such a user")
}
