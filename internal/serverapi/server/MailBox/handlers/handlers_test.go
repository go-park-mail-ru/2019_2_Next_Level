package handlers

import (
	"2019_2_Next_Level/pkg/TestTools"
	"2019_2_Next_Level/tests/mock/MailBox"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"net/http/httptest"
	"testing"
)

type Params TestTools.Params

func TestMailHandler_SendMail(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := MailBox.NewMockMailBoxUseCase(mockCtrl)
	h := NewMailHandler(mockUsecase)

	tests := []TestTools.TestStructMap{
		*TestTools.NewTestStructMap(
			map[string]TestTools.Params{},
			map[string]TestTools.Params{},
			map[string]TestTools.Params{

			}),
	}

	TestTools.RunTestingMapped(tests, func(map[string]TestTools.Params){}, func(test TestTools.TestStructMap) {
		body := bytes.Reader{}
		r := httptest.NewRequest("GET", "/auth.signin", &body)
		w := httptest.NewRecorder()
	})

}
