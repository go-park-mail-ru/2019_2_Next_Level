package handlers

import (
	"2019_2_Next_Level/internal/model"
	hr "2019_2_Next_Level/internal/serverapi/server/HttpError"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"2019_2_Next_Level/pkg/TestTools"
	MailBox "2019_2_Next_Level/tests/mock/serverapi"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Params TestTools.Params

func TestMailHandler_SendMail(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "returns":nil,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":""},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "returns": hr.GetError(hr.BadSession),

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":1, "returns":fmt.Errorf(""),

			}),
	}

	mail := models.MailToSend{To:[]string{"aa@aa.aa"}, Subject:"subject", Content:"content"}
	req := struct{
		Message models.MailToSend `json:"message"`
	}{mail}
	mailP := mail.ToMain()
	mailP.SetFrom(login)

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["returns"] != nil {
			err = test.MockParams["returns"].(error)
		}
		mockUsecase.EXPECT().SendMail(&mailP).Return(err).Times(test.MockParams["times"].(int))
		js, _ := json.Marshal(req)
		body := bytes.NewReader(js)
		r := httptest.NewRequest("GET", "/auth.signin", body)
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.SendMail(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_GetMailList(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	list := []model.Email{
		model.Email{From:"from", To: "to", Body:"body"},
	}
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "list":list,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":""},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":fmt.Errorf(""), "list":[]model.Email{},

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadParam)},
			map[string]TestTools.Params{"times":1, "err":fmt.Errorf(""), "list":[]model.Email{},

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockUsecase.EXPECT().GetMailList(login, "inbox", "", 1, 0).
			Return(test.MockParams["list"].([]model.Email), err).Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("POST", "/auth.signin?folder=inbox", strings.NewReader("folder=inbox"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.GetMailList(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_GetUnreadCount(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":""},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":1, "err":fmt.Errorf(""), "returns":0,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockUsecase.EXPECT().GetUnreadCount(test.Input["login"].(string)).Return(test.MockParams["returns"].(int), err).
			Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("GET", "/auth.signin?folder=inbox", strings.NewReader("folder=inbox"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.GetUnreadCount(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_GetEmail(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	email := model.Email{From:"from", To:"to", Body:"body"}
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "returns":email,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":""},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":nil, "returns":email,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login"},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadParam)},
			map[string]TestTools.Params{"times":1, "err":fmt.Errorf(""), "returns":model.Email{},

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockUsecase.EXPECT().GetMail(test.Input["login"].(string), []models.MailID{models.MailID(12345)}).
			Return([]model.Email{test.MockParams["returns"].(model.Email)}, err).Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("GET", "/auth.signin?id=12345", strings.NewReader(""))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.GetEmail(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_MarkMailRead(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	req := struct {
		Messages []models.MailID
	}{
		[]models.MailID{models.MailID(123)},
	}
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "req":req},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "mark":models.MarkMessageRead,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"", "req":req},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":nil, "mark":models.MarkMessageRead,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		js, _ := json.Marshal(req)
		body := bytes.NewReader(js)
		mockUsecase.EXPECT().MarkMail(test.Input["login"].(string), req.Messages, test.MockParams["mark"].(int)).
			Return(err).Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("POST", "/auth.signin", body)
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.MarkMailRead(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_ChangeMailFolder(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"folder", "id":int64(12345)},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"", "folder":"folder", "id":int64(12345)},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "folder":"folder", "id":int64(12345)},
			map[string]TestTools.Params{"resp": hr.GetError(hr.UnknownError)},
			map[string]TestTools.Params{"times":1, "err":fmt.Errorf(""), "returns":12,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockUsecase.EXPECT().
			ChangeMailFolder(test.Input["login"].(string), test.Input["folder"].(string), test.Input["id"].(int64)).
			Return(err).
			Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("GET", "/messages/createFolder/folder?folder=folder", strings.NewReader("folder=inbox"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		r = mux.SetURLVars(r, map[string]string{"name":"folder", "id":"12345"})
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.ChangeMailFolder(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}

func TestMailHandler_CreateFolder(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"folder"},
			map[string]TestTools.Params{"resp": &hr.DefaultResponse},
			map[string]TestTools.Params{"times":1, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"", "folder":"folder"},
			map[string]TestTools.Params{"resp": hr.GetError(hr.BadSession)},
			map[string]TestTools.Params{"times":0, "err":nil, "returns":12,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":"login", "folder":"folder"},
			map[string]TestTools.Params{"resp": hr.GetError(hr.UnknownError)},
			map[string]TestTools.Params{"times":1, "err":fmt.Errorf(""), "returns":12,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockUsecase.EXPECT().AddFolder(test.Input["login"].(string), test.Input["folder"].(string)).
			Return(err).
			Times(test.MockParams["times"].(int))
		r := httptest.NewRequest("GET", "/messages/createFolder/folder?folder=folder", strings.NewReader("folder=inbox"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		r = mux.SetURLVars(r, map[string]string{"name":"folder"})
		w := httptest.NewRecorder()
		r.Header = http.Header{"X-Login": []string{test.Input["login"].(string)}}


		h.CreateFolder(w, r)
		var got hr.HttpResponse
		_ =json.Unmarshal([]byte(w.Body.String()), &got)
		if !cmp.Equal(got.Status, test.Expected["resp"].(*hr.HttpResponse).Status) {
			t.Error("Wrong resp")
		}

	})

}