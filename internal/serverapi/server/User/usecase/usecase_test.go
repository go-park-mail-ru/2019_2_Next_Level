package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/mock"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/golang/mock/gomock"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mock.NewMockUserRepository(mockCtrl)
	usecase := NewUserUsecase(mockRepo)

	login := "ivanov"
	user := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan", "", ""}

	type F func()
	funcs := []F{
		func() {
			mockRepo.EXPECT().GetUser(login).Return(user, nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().GetUser(login).Return(user, e.Error{}.SetCode(e.InvalidParams)).Times(1)
		},
	}

	expected := []error{
		nil,
		e.Error{}.SetCode(e.InvalidParams),
	}
	for i, resp := range expected {
		funcs[i]()
		gotUser, err := usecase.GetUser(login)
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", err, resp)
		}
		if !cmp.Equal(gotUser, user) {
			t.Errorf("Wrong answer got: %s instead %s\n", gotUser, user)
		}
	}
}

func TestEditUser(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mock.NewMockUserRepository(mockCtrl)
	usecase := NewUserUsecase(mockRepo)

	user := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan@", "", ""}

	type F func()
	funcs := []F{
		func() {
			mockRepo.EXPECT().UpdateUserData(&user).Return(nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().UpdateUserData(&user).Return(e.Error{}.SetCode(e.InvalidParams)).Times(1)
		},
	}

	expected := []error{
		nil,
		e.Error{}.SetCode(e.InvalidParams),
	}
	for i, resp := range expected {
		funcs[i]()
		err := usecase.EditUser(&user)
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", err, resp)
		}
	}
}

func TestEditPassword(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mock.NewMockUserRepository(mockCtrl)
	usecase := NewUserUsecase(mockRepo)

	login := "ivanov"

	type F func()
	funcs := []F{
		func() {
			mockRepo.EXPECT().GetUserCredentials(login).Return("12345", "sault", nil)
			mockRepo.EXPECT().UpdateUserPassword(login, "54321", "sault").Return(nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().GetUserCredentials(login).Return("12345", "sault", nil)
		},
		func() {
			mockRepo.EXPECT().GetUserCredentials(login).Return("", "", e.Error{}.SetCode(e.NotExists))
		},
	}

	input := []string{
		"12345", "54321",
		"12345", "12345",
		"", "",
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.Wrong),
		e.Error{}.SetCode(e.NotExists),
	}
	for i, resp := range expected {
		funcs[i]()
		err := usecase.EditPassword(login, input[i], input[i+1])
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", err, resp)
		}
	}
}
