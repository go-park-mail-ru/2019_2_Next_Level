package usecase

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/post"
	"2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	e "2019_2_Next_Level/pkg/HttpError/Error"
	"2019_2_Next_Level/pkg/TestTools"
	"2019_2_Next_Level/tests/mock/MailBox"
	"2019_2_Next_Level/tests/mock/postinterface"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"strconv"
	"testing"
)
func init() {
	defaultConf := config.Database{DBName: "nextlevel", Port: "5432", Host: "localhost", User: "postgres", Password: "postgres"}
	config.Conf.DB = defaultConf
}
func TestMailBoxUsecase_GetMailList(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := MailBox.NewMockMailRepository(mockCtrl)
	mockService := postinterface.NewMockIPostInterface(mockCtrl)
	u := NewMailBoxUsecase(mockRepo, mockService)

	mailList := []model.Email{
		model.Email{From:"from", To:"to", Body:"body"},
	}
	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"inbox"},
			map[string]TestTools.Params{"list":mailList},
			map[string]TestTools.Params{"times":1, "err":nil,

		}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"inbox"},
			map[string]TestTools.Params{"list":mailList},
			map[string]TestTools.Params{"times":1, "err": e.Error{}.SetError(e.Error{}).SetCode(e.ProcessError),

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		login := test.Input["login"].(string)
		folder := test.Input["folder"].(string)
		sort, from, to := "", 1, 25
		mockRepo.EXPECT().GetEmailList(login, folder, sort, from, to).
			Return(mailList, err).Times(test.MockParams["times"].(int))
		list, err := u.GetMailList(login, folder, sort, from, to)
		var expected []model.Email
		if test.Expected["list"] != nil {
			expected = test.Expected["list"].([]model.Email)
			//expected[0].Sanitize()
		}
		if !cmp.Equal(list, expected) {
			t.Errorf("Wrong result")
		}
	})
}

func TestMailBoxUsecase_GetMailListPlain(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := MailBox.NewMockMailRepository(mockCtrl)
	mockService := postinterface.NewMockIPostInterface(mockCtrl)
	u := NewMailBoxUsecase(mockRepo, mockService)

	mailList := []model.Email{
		model.Email{From:"from", To:"to", Body:"body"},
	}
	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"inbox"},
			map[string]TestTools.Params{"list":mailList},
			map[string]TestTools.Params{"times":1, "times2":1,"err":nil, "err2": nil,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"inbox"},
			map[string]TestTools.Params{"list":[]model.Email{}},
			map[string]TestTools.Params{"times":1, "times2":0,"err": e.Error{}, "err2": nil,

			}),
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "folder":"inbox"},
			map[string]TestTools.Params{"list":mailList},
			map[string]TestTools.Params{"times":1, "times2":1,"err":nil, "err2": e.Error{},

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		login := test.Input["login"].(string)
		//folder := test.Input["folder"].(string)
		sort, from, to := "", 1, 25
		mockRepo.EXPECT().GetMessagesCount(login, models.InboxFolder, models.FlagMessageTotal).
			Return(10, err).Times(test.MockParams["times"].(int))
		if test.MockParams["err2"] != nil {
			err = test.MockParams["err2"].(error)
		}
		mockRepo.EXPECT().GetEmailList(login, models.InboxFolder, sort, from, to).
			Return(mailList, err).Times(test.MockParams["times2"].(int))

		_, _, list, err := u.GetMailListPlain(login, 1)
		var expected []model.Email
		if test.Expected["list"] != nil {
			expected = test.Expected["list"].([]model.Email)
			//expected[0].Sanitize()
		}
		if !cmp.Equal(list, expected) {
			t.Errorf("Wrong result")
		}
	})
}

func TestMailBoxUsecase_GetMail(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := MailBox.NewMockMailRepository(mockCtrl)
	mockService := postinterface.NewMockIPostInterface(mockCtrl)
	u := NewMailBoxUsecase(mockRepo, mockService)

	mail := model.Email{From:"from", To:"to", Body:"body"}
	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "id":models.MailID(1)},
			map[string]TestTools.Params{"res": mail},
			map[string]TestTools.Params{"times":1, "err":nil,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		login := test.Input["login"].(string)
		id := test.Input["id"].(models.MailID)
		mockRepo.EXPECT().GetEmailByCode(login, strconv.Itoa(int(id))).
			Return(test.Expected["res"].(model.Email), err).Times(test.MockParams["times"].(int))

		res, err := u.GetMail(login, id)
		var expected model.Email
		if test.Expected["res"] != nil {
			expected = test.Expected["res"].(model.Email)
			//expected[0].Sanitize()
		}
		if !cmp.Equal(res, expected) {
			t.Errorf("Wrong result")
		}
	})
}

func TestMailBoxUsecase_SendMail(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := MailBox.NewMockMailRepository(mockCtrl)
	mockService := postinterface.NewMockIPostInterface(mockCtrl)
	u := NewMailBoxUsecase(mockRepo, mockService)
	config.Conf.HttpConfig.HostName = "w"

	email := model.Email{From:"from", To:"to@dd", Body:"body"}
	postEmail := post.Email{
		From: email.From,
		To:   email.To,
		Body: email.Body,
		Subject:email.Header.Subject,
	}

	login := "login"
	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{"login":login, "id":models.MailID(1)},
			map[string]TestTools.Params{"res": nil},
			map[string]TestTools.Params{"mess":postEmail, "times":1, "err":nil,"times2":1, "err2":nil,

			}),
	}


	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		var err error
		if test.MockParams["err"] != nil {
			err = test.MockParams["err"].(error)
		}
		mockService.EXPECT().
			//Put(test.MockParams["mess"].(post.Email)).
			Put(gomock.Any()).
			Return(err).Times(test.MockParams["times"].(int))
		if test.MockParams["err2"] != nil {
			err = test.MockParams["err2"].(error)
		}
		emailT := email
		emailT.From+="@w"
		mockRepo.EXPECT().PutSentMessage(emailT).Times(test.MockParams["times2"].(int))

		err = u.SendMail(&email)
		var expected error
		if test.Expected["res"] != nil {
			expected = test.Expected["res"].(error)
		}
		if !cmp.Equal(err, expected) {
			t.Errorf("Wrong result")
		}
	})
}