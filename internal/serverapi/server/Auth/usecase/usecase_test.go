package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/mock"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	"testing"

	"github.com/google/uuid"

	e "2019_2_Next_Level/internal/serverapi/server/error"

	"github.com/golang/mock/gomock"
)

func TestCheckAuth(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mock.NewMockRepository(mockCtrl)
	u := NewAuthUsecase(mockRepo)

	uuidTemp, _ := uuid.NewUUID()
	uuid := uuidTemp.String()
	funcs := []func(){
		func() {
			mockRepo.EXPECT().GetLoginBySession(uuid).Return("ian", nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().GetLoginBySession(uuid).Return("", e.Error{}.SetCode(e.NoPermission)).Times(1)
		},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.NoPermission),
	}

	for i, test := range expected {
		funcs[i]()
		got := u.CheckAuth(uuid)
		if !e.CompareErrors(got, test, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", got, test)
		}
	}

}

func TestSignOut(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mock.NewMockRepository(mockCtrl)
	u := NewAuthUsecase(mockRepo)

	uuidTemp, _ := uuid.NewUUID()
	uuid := uuidTemp.String()
	funcs := []func(){
		func() {
			mockRepo.EXPECT().DeleteSession(uuid).Return(nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().DeleteSession(uuid).Return(e.Error{}.SetCode(e.NotExists)).Times(1)
		},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.InvalidParams),
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
	mockRepo := mock.NewMockRepository(mockCtrl)
	u := NewAuthUsecase(mockRepo)

	testUser := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan", "12345", ""}

	funcs := []func(){
		func() {
			mockRepo.EXPECT().AddNewUser(&testUser).Return(nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().AddNewUser(&testUser).Return(e.Error{}.SetCode(e.AlreadyExists)).Times(1)
		},
		func() {
			mockRepo.EXPECT().AddNewUser(&testUser).Return(
				e.Error{}.SetCode(auth.ErrorWrongNickName)).
				Times(1)
		},
	}
	expected := []error{
		nil,
		e.Error{}.SetCode(e.AlreadyExists),
		e.Error{}.SetCode(e.InvalidParams).SetError(e.Error{}.SetCode(auth.ErrorWrongNickName)),
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
	mockRepo := mock.NewMockRepository(mockCtrl)
	u := NewAuthUsecase(mockRepo)

	testUser := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{"ivanov@mail.ru", "12345"}

	funcs := []func(){
		func() {
			mockRepo.EXPECT().GetUserCredentials(testUser.Login).Return([]string{testUser.Password, ""}, nil).Times(1)
			mockRepo.EXPECT().AddNewSession(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		},
		func() {
			mockRepo.EXPECT().GetUserCredentials(testUser.Login).
				Return([]string{testUser.Password, ""}, e.Error{}.SetCode(e.NotExists)).
				Times(1)
		},
		func() {
			mockRepo.EXPECT().GetUserCredentials(testUser.Login).Return([]string{testUser.Password + "rubbish", ""}, nil).Times(1)
		},
	}
	type ReturnParams struct {
		UUID  string
		Error error
	}
	expected := []ReturnParams{
		ReturnParams{"uuid", nil},
		ReturnParams{"", e.Error{}.SetCode(e.NotExists)},
		ReturnParams{"", e.Error{}.SetCode(e.InvalidParams).SetError(e.Error{}.SetCode(auth.ErrorWrongPassword))},
	}

	for i, test := range expected {
		funcs[i]()
		gotId, gotError := u.SignIn(testUser.Login, testUser.Password)
		if !e.CompareErrors(gotError, test.Error, e.CompareByCode) {
			t.Errorf("Wrong response: %s\nWanted: %s", gotError, test.Error)
		}
		if gotId == "" && test.UUID != "" {
			t.Errorf("Empty UUID returned")
		}
	}
}
