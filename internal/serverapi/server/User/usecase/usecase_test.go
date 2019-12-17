package usecase

import (
	"2019_2_Next_Level/internal/model"
	e "2019_2_Next_Level/pkg/Error"
	User "2019_2_Next_Level/tests/mock/serverapi"
	authclient "2019_2_Next_Level/tests/mock/serverapi/auth"
	"mime/multipart"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/golang/mock/gomock"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := User.NewMockUserRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	usecase := NewUserUsecase(mockRepo, mockService)

	login := "ivanov"
	user := model.User{Name:"Ivan", Sirname:"Ivanov", BirthDate:"01.01.1900", Sex:"male", Email:"ivan", Password:"12345"}

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
			//t.Errorf("Wrong answer got: %s instead %s\n", gotUser, user)
		}
	}
}

func TestEditUser(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := User.NewMockUserRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	usecase := NewUserUsecase(mockRepo, mockService)

	user := model.User{Name:"Ivan", Sirname:"Ivanov", BirthDate:"01.01.1900", Sex:"male", Email:"ivan", Password:"12345"}

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
	mockRepo := User.NewMockUserRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	usecase := NewUserUsecase(mockRepo, mockService)

	login := "ivanov"

	type F func()
	funcs := []F{
		func() {
			mockService.EXPECT().ChangePassword(login, "12345", "54321").Return(int32(e.OK))
			mockService.EXPECT().GetError(int32(e.OK)).Return(nil)
		},
		func() {
			mockService.EXPECT().ChangePassword(login, "12345", "12345").Return(int32(e.WrongPassword))
			mockService.EXPECT().GetError(int32(e.WrongPassword)).Return(e.Error{}.SetCode(e.WrongPassword))
		},
	}

	input := []string{
		"12345", "54321",
		"12345", "12345",
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.WrongPassword),
	}
	for i, resp := range expected {
		funcs[i]()
		err := usecase.EditPassword(login, input[2*i], input[2*i+1])
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", err, resp)
		}
	}
}

func TestUserUsecase_GetUserFolders(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := User.NewMockUserRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	usecase := NewUserUsecase(mockRepo, mockService)

	login := "ivanov"

	type F func()
	funcs := []F{
		func() {
			mockRepo.EXPECT().GetUserFolders(login).Return([]model.Folder{}, nil)
		},
		func() {
			mockRepo.EXPECT().GetUserFolders(login).Return([]model.Folder{}, e.Error{})
		},
	}

	expected := []error{
		nil,
		e.Error{},
	}
	for i, resp := range expected {
		funcs[i]()
		_, err := usecase.GetUserFolders(login)
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", err, resp)
		}
	}
}

func TestUserUsecase_EditAvatar(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := User.NewMockUserRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	usecase := NewUserUsecase(mockRepo, mockService)
	mockFile := User.NewMockFile(mockCtrl)
	//mockF := User.NewMockWriter(mockCtrl)

	login := "ivanov"
	//name := "file"

	type F func()
	//funcs := []F{
	//	func() {
	//		mockOS.EXPECT().OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666).Return(mockF, nil)
	//	},
	//	//func() {
	//	////	mockRepo.EXPECT().GetUserFolders(login).Return([]model.Folder{}, e.Error{})
	//	//},
	//}

	expected := []error{
		nil,
		e.Error{},
	}
	for _, resp := range expected {
		//funcs[i]()
		_, err := usecase.EditAvatar(login, mockFile, &multipart.FileHeader{})
		if !e.CompareErrors(err, resp, e.CompareByCode) {
			//t.Errorf("Wrong response: %v\nWanted: %v", err, resp)
		}
	}
}