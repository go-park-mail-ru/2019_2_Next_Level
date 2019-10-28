package repository

import (
	"2019_2_Next_Level/internal/model"
	"fmt"
)

type MockRepository struct {
	sessionBook map[string]string
	usersList   map[string]model.User
}

func GetMock() *MockRepository {
	r := MockRepository{}
	return &r
}

func (r *MockRepository) AddNewSession(login, uuid string) error {
	if err := r.checkUserExist(login); err != nil {
		return err
	}
	r.sessionBook[uuid] = login
	return nil
}


func (r *MockRepository) CheckSession(id string) error {
	if _, ok := r.sessionBook[id]; !ok {
		return fmt.Errorf("Wrong")
	}
	return nil
}
func (r *MockRepository) DiscardSession(id string) error {
	if _, ok := r.sessionBook[id]; !ok {
		return fmt.Errorf("Wrong")
	}
	delete(r.sessionBook, id)
	return nil
}

func (r *MockRepository) Registrate(user *model.User) error {
	if _, ok := r.usersList[user.Email]; !ok {
		return fmt.Errorf("User exists")
	}
	r.usersList[user.Email] = *user
	return nil
}

func (r *MockRepository) GetUserCredentials(login string) (string, string, error) {
	user, ok := r.usersList[login]
	if !ok {
		return "", "", fmt.Errorf("User exists")
	}
	return user.Email, user.Password, nil
}

func (r *MockRepository) checkUserExist(login string) error {
	_, err := r.usersList[login]
	if err {
		return nil
	}
	return fmt.Errorf("No such a user")
}
