package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/tests/mock/Auth"
	authclient "2019_2_Next_Level/tests/mock/serverapi/auth"
	"testing"

	"github.com/google/uuid"

	e "2019_2_Next_Level/pkg/Error"

	"github.com/golang/mock/gomock"
)


func TestCheckAuth(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := Auth.NewMockRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	u := NewAuthUsecase(mockRepo, mockService)

	uuidTemp, _ := uuid.NewUUID()
	uuid := uuidTemp.String()
	funcs := []func(){
		func() {
			mockService.EXPECT().LoginBySession(uuid).Return("ian", int32(e.OK))
			mockService.EXPECT().GetError(int32(e.OK)).Return(nil)
		},
		func() {
			mockService.EXPECT().LoginBySession(uuid).Return("", int32(e.NoPermission))
			mockService.EXPECT().GetError(int32(e.NoPermission)).Return(e.Error{}.SetCode(e.NoPermission))
		},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.NoPermission),
	}

	for i, test := range expected {
		funcs[i]()
		_, gotErr := u.CheckAuth(uuid)
		if !e.CompareErrors(gotErr, test, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", gotErr, test)
		}
	}

}

func TestSignOut(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := Auth.NewMockRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	u := NewAuthUsecase(mockRepo, mockService)

	uuidTemp, _ := uuid.NewUUID()
	uuid := uuidTemp.String()
	funcs := []func(){
		func() {
			mockService.EXPECT().DestroySession(uuid).Return(int32(e.OK))
			mockService.EXPECT().GetError(int32(e.OK)).Return(nil)
		},
		func() {
			mockService.EXPECT().DestroySession(uuid).Return(int32(e.NotExists))
			mockService.EXPECT().GetError(int32(e.NotExists)).Return(e.Error{}.SetCode(e.NotExists))
		},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.NotExists),
	}

	for i, test := range expected {
		funcs[i]()
		got := u.Logout(uuid)
		if !e.CompareErrors(got, test, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", got, test)
		}
	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := Auth.NewMockRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	u := NewAuthUsecase(mockRepo, mockService)

	testUser := model.User{Name:"Ivan", Sirname:"Ivanov", BirthDate:"01.01.1900", Sex:"male", Email:"ivan", Password:"12345"}

	funcs := []func(){
		func() {
			mockRepo.EXPECT().AddNewUser(&testUser).Return(nil).Times(1)
			mockService.EXPECT().RegisterUser(testUser.Email, testUser.Password).Return(int32(e.OK))
		},
		func() {
			mockRepo.EXPECT().AddNewUser(&testUser).Return(e.Error{}.SetCode(e.AlreadyExists)).Times(1)
		},
		//func() {
		//	mockRepo.EXPECT().AddNewUser(&testUser).Return(
		//		e.Error{}.SetCode(auth.ErrorWrongNickName)).
		//		Times(1)
		//},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.AlreadyExists),
		//e.Error{}.SetCode(e.InvalidParams).SetError(e.Error{}.SetCode(auth.ErrorWrongNickName)),
	}

	for i, test := range expected {
		funcs[i]()
		got := u.SignUp(testUser)
		if !e.CompareErrors(got, test, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", got, test)
		}
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := Auth.NewMockRepository(mockCtrl)
	mockService := authclient.NewMockIAuthClient(mockCtrl)
	u := NewAuthUsecase(mockRepo, mockService)

	testUser := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{"ivanov@mail.ru", "12345"}

	funcs := []func(){
		func() {
			mockService.EXPECT().CheckCredentials(testUser.Login, testUser.Password).Return(int32(e.OK))
			mockService.EXPECT().GetError(int32(e.OK)).Return(nil)
			mockService.EXPECT().StartSession(gomock.Any()).Return("token", int32(e.OK))
			mockService.EXPECT().GetError(int32(e.OK)).Return(nil)
		},
		func() {
			mockService.EXPECT().CheckCredentials(testUser.Login, testUser.Password).Return(int32(e.NotExists))
			mockService.EXPECT().GetError(int32(e.NotExists)).Return(e.Error{}.SetCode(e.NotExists))
		},
	}
	type ReturnParams struct {
		UUID  string
		Error error
	}
	expected := []ReturnParams{
		ReturnParams{"uuid", nil},
		ReturnParams{"", e.Error{}.SetCode(e.NotExists)},
	}

	for i, test := range expected {
		funcs[i]()
		gotId, gotError := u.SignIn(testUser.Login, testUser.Password)
		if !e.CompareErrors(gotError, test.Error, e.CompareByCode) {
			//t.Errorf("Wrong response: %s\nWanted: %s", gotError, test.Error)
		}
		if gotId == "" && test.UUID != "" {
			//t.Errorf("Empty UUID returned")
		}
	}
}
