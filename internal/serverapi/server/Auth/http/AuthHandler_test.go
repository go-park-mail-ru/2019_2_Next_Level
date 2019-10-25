package http

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/mock"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestIsAuth(t *testing.T) {
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
		{token.String(), nil, `{"status":"OK"}`},
		{"123", e.Error{}.SetCode(e.InvalidParams), `{"status":"4"}`},
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
		mockUsecase.EXPECT().CheckAuth(test.param).Return(test.res).Times(1)
		h.CheckAuthorization(w, r)
		got := w.Body.String()
		if test.expected != got {
			t.Errorf("Wrong response: %s instead %s", got, test.expected)
		}

	}
}

func TestLogout(t *testing.T) {
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
		{token.String(), nil, `{"status":"OK"}`},
		{"123", e.Error{}.SetCode(e.InvalidParams), `{"status":"4"}`},
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
		if test.expected != got {
			t.Errorf("Wrong response: %s", got)
		}

	}
}

func TestSignUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	h := NewAuthHandler(mockUsecase)

	testUser := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan", "12345", ""}

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
	tests := []string{
		`{"status":"OK"}`,
		`{"status":"15"}`,
		`{"status":"19"}`,
		`{"status":"1"}`,
	}
	for i, test := range tests {
		js, _ := json.Marshal(testUser)
		body := bytes.NewReader(js)
		r := httptest.NewRequest("GET", "/auth.signin", body)
		w := httptest.NewRecorder()

		funcs[i]()
		h.SignUp(w, r)
		got := w.Body.String()
		if test != got {
			t.Errorf("Wrong response: %s", got)
		}
	}
}
