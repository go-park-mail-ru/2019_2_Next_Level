package http

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/log"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	httperror "2019_2_Next_Level/pkg/HttpError/Error/httpError"
	mockk "2019_2_Next_Level/tests/mock"
	"2019_2_Next_Level/tests/mock/mock"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func init() {
	log.SetLogger(&mockk.MockLog{})
}

func TestIsAuth(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	h := NewAuthHandler(mockUsecase)

	body := bytes.Reader{}

	token, _ := uuid.NewUUID()

	tests := []struct {
		param    string
		res      error
		expected string
	}{
		{token.String(), nil, `{"status":"ok"}`},
		{"123", e.Error{}.SetCode(e.InvalidParams), `{"status":"error","error":{"code":4,"msg":"User is not authorized"}}`},
	}
	for _, test := range tests {
		r := httptest.NewRequest("GET", "/auth.signin", &body)
		w := httptest.NewRecorder()
		cookie := http.Cookie{
			Name:  sessionTokenCookieName,
			Value: test.param,
		}
		http.SetCookie(w, &cookie)
		r.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
		mockUsecase.EXPECT().CheckAuth(test.param).Return("", test.res).Times(1)
		h.CheckAuthorization(w, r)
		got := w.Body.String()
		if test.expected != got {
			t.Errorf("Wrong response: %s instead %s", got, test.expected)
		}

	}
}

func TestLogout(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	h := NewAuthHandler(mockUsecase)

	body := bytes.Reader{}

	token, _ := uuid.NewUUID()
	tests := []struct {
		param    string
		res      error
		expected interface{}
	}{
		{token.String(), nil, httperror.HttpResponse{Status: "ok"}},
		{"123", e.Error{}.SetCode(e.InvalidParams), httperror.GetError(httperror.BadSession)},
	}

	for _, test := range tests {
		r := httptest.NewRequest("GET", "/auth.signin", &body)
		w := httptest.NewRecorder()
		cookie := http.Cookie{
			Name:  sessionTokenCookieName,
			Value: test.param,
		}
		http.SetCookie(w, &cookie)
		r.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
		mockUsecase.EXPECT().Logout(test.param).Return(test.res).Times(1)
		h.SignOut(w, r)
		got := w.Body.String()
		expected, _ := json.Marshal(test.expected)
		if string(expected) != got {
			t.Errorf("Wrong response: %s", got)
		}

	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	h := NewAuthHandler(mockUsecase)

	testUser := model.User{Name:"Ivan", Sirname:"Ivanov", BirthDate:"01.01.1900", Sex:"male", Email:"ivan", Password:"12345"}

	type F func()
	funcs := []F{
		func() {
			mockUsecase.EXPECT().SignUp(testUser).Return(nil).Times(1)
			mockUsecase.EXPECT().SignIn(testUser.Email, testUser.Password).Return("12345", nil).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignUp(testUser).Return(e.Error{}.SetCode(e.InvalidParams).IncludeError(e.Error{}.SetCode(auth.ErrorWrongBirthDate))).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignUp(testUser).Return(e.Error{}.SetCode(e.AlreadyExists)).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignUp(testUser).Return(nil).Times(1)
			mockUsecase.EXPECT().SignIn(testUser.Email, testUser.Password).Return("", e.Error{}.SetCode(e.InvalidParams)).Times(1)
		},
	}
	response := []interface{}{
		httperror.HttpResponse{Status: "ok"},
		httperror.GetError(httperror.IncorrectBirthDate),
		httperror.GetError(httperror.LoginAlreadyExists),
		httperror.GetError(httperror.UnknownError),
	}
	for i, test := range response {
		js, _ := json.Marshal(testUser)
		body := bytes.NewReader(js)
		r := httptest.NewRequest("GET", "/auth.signin", body)
		w := httptest.NewRecorder()

		funcs[i]()
		h.SignUp(w, r)
		got := w.Body.String()
		expected, _ := json.Marshal(test)
		if string(expected) != got {
			t.Errorf("Wrong response: %s", got)
		}
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	h := NewAuthHandler(mockUsecase)

	testUser := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{"ivanov@mail.ru", "12345"}

	type F func()
	funcs := []F{
		func() {
			mockUsecase.EXPECT().SignIn(testUser.Login, testUser.Password).Return("uuid", nil).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignIn(testUser.Login, testUser.Password).Return("", e.Error{}.SetCode(e.InvalidParams).IncludeError(e.Error{}.SetCode(auth.ErrorWrongPassword))).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignIn(testUser.Login, testUser.Password).Return("", e.Error{}.SetCode(e.InvalidParams).IncludeError(e.Error{}.SetCode(auth.ErrorWrongLogin))).Times(1)
		},
		func() {
			mockUsecase.EXPECT().SignIn(testUser.Login, testUser.Password).Return("", e.Error{}.SetCode(e.NotExists)).Times(1)
		},
	}
	response := []interface{}{
		httperror.HttpResponse{Status: "ok"},
		httperror.GetError(httperror.WrongPassword),
		httperror.GetError(httperror.WrongPassword),
		httperror.GetError(httperror.LoginNotExist),
	}
	for i, test := range response {
		js, _ := json.Marshal(testUser)
		body := bytes.NewReader(js)
		r := httptest.NewRequest("GET", "/auth.signin", body)
		w := httptest.NewRecorder()

		funcs[i]()
		h.SignIn(w, r)
		got := w.Body.String()
		expected, _ := json.Marshal(test)
		if string(expected) != got {
			t.Errorf("Wrong response: %s\nWanted: %s", got, expected)
		}
	}
}
