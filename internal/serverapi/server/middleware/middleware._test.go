package middleware

import (
	"2019_2_Next_Level/internal/serverapi/mock"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := mock.NewMockUsecase(mockCtrl)
	mockHttpHandler := mock.NewMockHandler(mockCtrl)
	// h := authHttp.NewAuthHandler(mockUsecase)

	// body := bytes.Reader{}
	// r := httptest.NewRequest("GET", "/auth.signin", &body)
	// w := httptest.NewRecorder()

	login := "aaa"
	uuid := "12345"

	// cookie := http.Cookie{
	// 	Name:  "session-id",
	// 	Value: uuid,
	// }
	// http.SetCookie(w, &cookie)
	// r.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	// mockUsecase.EXPECT().CheckAuth(uuid).Return(login, nil).Times(1)
	// r.Header.Set("login", login)
	// mockHttpHandler.EXPECT().ServeHTTP(w, r).Return().Times(1)

	type F func(*httptest.ResponseRecorder, *http.Request)
	funcs := []F{
		func(w *httptest.ResponseRecorder, r *http.Request) {
			cookie := http.Cookie{
				Name:  "session-id",
				Value: uuid,
			}
			http.SetCookie(w, &cookie)
			r.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
			mockUsecase.EXPECT().CheckAuth(uuid).Return(login, nil).Times(1)
			r.Header.Set("login", login)
			mockHttpHandler.EXPECT().ServeHTTP(w, r).Return().Times(1)
		},
	}

	for _, f := range funcs {
		body := bytes.Reader{}
		r := httptest.NewRequest("GET", "/auth.signin", &body)
		w := httptest.NewRecorder()
		f(w, r)
		(AuthentificationMiddleware(mockUsecase)(mockHttpHandler)).ServeHTTP(w, r)
		got := w.Body.String()
		expected := `{"status":"ok"}`
		if expected != got && got != "" {
			t.Errorf("Wrong response: %s instead %s", got, expected)
		}

	}

}
