package middleware

import (
	"2019_2_Next_Level/internal/serverapi/config"
	"2019_2_Next_Level/internal/serverapi/log"
	mokk "2019_2_Next_Level/tests/mock"
	"2019_2_Next_Level/tests/mock/Auth"
	"2019_2_Next_Level/tests/mock/mock"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func init() {
	log.SetLogger(&mokk.MockLog{})
}

func TestStaticMiddleware(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockHttpHandler := mock.NewMockHandler(mockCtrl)

	body := bytes.Reader{}
	r := httptest.NewRequest("GET", "/", &body)
	w := httptest.NewRecorder()
	expectedRequest := r
	expectedRequest.URL.Path += "index.html"
	mockHttpHandler.EXPECT().ServeHTTP(w, expectedRequest)
	StaticMiddleware()(mockHttpHandler).ServeHTTP(w,r)
}

func TestAuth(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUsecase := Auth.NewMockUsecase(mockCtrl)
	mockHttpHandler := mock.NewMockHandler(mockCtrl)

	login := "aaa"
	uuid := "12345"

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

func TestCors(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockHttpHandler := mock.NewMockHandler(mockCtrl)
	r := httptest.NewRequest("POST", "/profile/get", nil)
	w := httptest.NewRecorder()
	mockHttpHandler.EXPECT().ServeHTTP(w, r).Return().Times(1)
	hostName := "test"
	config.Conf.HttpConfig.Whitelist = make(map[string]bool)
	config.Conf.HttpConfig.Whitelist[hostName] = true

	r.Header.Add("Origin", hostName)
	(CorsMethodMiddleware()(mockHttpHandler)).ServeHTTP(w, r)

	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") != hostName ||
		headers.Get("Access-Control-Allow-Credentials") != "true" ||
		headers.Get("Access-Control-Allow-Headers") != "Content-Type" ||
		headers.Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS, PUT, DELETE" {
		t.Error("Wrong headers got")
	}
}