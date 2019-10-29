package http

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/mock"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"2019_2_Next_Level/pkg/HttpTools"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGet(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUserUsecase(mockCtrl)
	h := NewUserHandler(mockUsecase)
	login := "Ian"
	type Answer struct {
		Status    string `json:"status"`
		Name      string `json:"firstName"`
		Sirname   string `json:"secondName"`
		BirthDate string `json:"birthDate"`
		Sex       string `json:"sex"`
		Email     string `json:"login"`
		Avatar    string `json:"avatar"`
	}
	user := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan", "12345", ""}
	userResponse, err := json.Marshal(Answer{
		"ok", "Ivan", "Ivanov", "01.01.1900", "male", "ivan", "",
	})
	if err != nil {
		t.Errorf("Cannot get json answer")
		return
	}
	type F func()
	funcs := []F{
		func() {
			mockUsecase.EXPECT().GetUser(login).Return(user, nil).Times(1)
		},
	}
	response := []string{
		string(userResponse),
	}

	for i, resp := range response {
		body := &bytes.Reader{}
		r := httptest.NewRequest("GET", "/user/get", body)
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{login}}
		funcs[i]()
		h.GetProfile(w, r)
		got := w.Body.String()
		if got != resp {
			t.Errorf("Wrong answer got: %s instead %s\n", got, resp)
		}
	}
}

func TestEditUser(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUserUsecase(mockCtrl)
	h := NewUserHandler(mockUsecase)

	user := model.User{"Ivan", "Ivanov", "", "01.01.1900", "male", "ivan", "12345", ""}
	// type Data struct {
	// 	Name      string `json:"firstName"`
	// 	Sirname   string `json:"secondName"`
	// 	BirthDate string `json:"birthDate"`
	// 	NickName  string `json:"nickName"`
	// 	Sex       string `json:"sex"`
	// 	Avatar    []byte `json:"avatar"`
	// }

	type F func()
	funcs := []F{
		func() {
			mockUsecase.EXPECT().EditUser(&user).Return(nil).Times(1)
		},
		func() {
			mockUsecase.EXPECT().EditUser(&user).Return(
				e.Error{}.SetCode(e.InvalidParams).SetError(e.Error{}.SetCode(auth.ErrorWrongPassword))).Times(1)
		},
	}
	response := []string{
		`{"status":"ok"}`,
		(&HttpTools.Response{}).SetError(hr.GetError(hr.IncorrectPassword)).String(),
	}
	for i, resp := range response {
		js, _ := json.Marshal(user)
		body := bytes.NewReader(js)
		r := httptest.NewRequest("GET", "/user/get", body)
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{user.Email}}
		funcs[i]()
		h.EditUserInfo(w, r)
		got := w.Body.String()
		if got != resp {
			t.Errorf("Wrong answer got: %s instead %s\n", got, resp)
		}
	}

}

func TestEditPassword(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUserUsecase(mockCtrl)
	h := NewUserHandler(mockUsecase)
	login := "ivanov"

	type F func()
	funcs := []F{
		func() {
			mockUsecase.EXPECT().EditPassword(login, "12345", "54321").Return(nil).Times(1)
		},
		func() {
			mockUsecase.EXPECT().EditPassword(login, "12345", "54321").Return(
				e.Error{}.SetCode(e.Wrong)).
				Times(1)
		},
		func() {},
	}
	input := []string{
		`{"currentPassword":"12345","newPassword":"54321"}`,
		`{"currentPassword":"12345","newPassword":"54321"}`,
		`{invalidJSON}`,
	}
	response := []string{
		`{"status":"ok"}`,
		(&HttpTools.Response{}).SetError(hr.GetError(hr.SameNewPass)).String(),
		(&HttpTools.Response{}).SetError(hr.GetError(hr.BadParam)).String(),
	}
	for i, resp := range response {
		js := input[i]
		body := bytes.NewReader([]byte(js))
		r := httptest.NewRequest("GET", "/user/get", body)
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{login}}
		funcs[i]()
		h.EditUserPassword(w, r)
		got := w.Body.String()
		if got != resp {
			t.Errorf("Wrong answer got: %s instead %s\n", got, resp)
		}
	}
}
