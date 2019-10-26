package http

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	e "2019_2_Next_Level/internal/serverapi/server/error"
	"2019_2_Next_Level/pkg/HttpTools"
	"net/http"
	"time"
)

const (
	sessionTokenCookieName = "session-id"
	cookieLifetime         = 10 * 60 // minutes
)

type AuthHandler struct {
	usecase auth.Usecase
	resp    *HttpTools.Response
}

func NewAuthHandler(uc auth.Usecase) AuthHandler {
	resp := (&HttpTools.Response{}).SetError(hr.DefaultResponse)
	return AuthHandler{usecase: uc, resp: resp}
}

func (a *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	resp := a.resp.SetWriter(w).Copy()
	defer resp.Send()
	tokenCookie, err := r.Cookie(sessionTokenCookieName)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	err = a.usecase.Logout(tokenCookie.Value)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	a.setCookie(w, sessionTokenCookieName, "delete", time.Now().AddDate(0, 0, -1))
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	resp := a.resp.SetWriter(w).Copy()
	defer resp.Send()
	userInput := model.User{}
	err := HttpTools.StructFromBody(*r, &userInput)
	if err != nil {
		// a.resp.SetStatus(hr.BadParam)
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	err = a.usecase.SignUp(userInput)
	if err != nil {
		status := hr.UnknownError
		switch err.(type) {
		case e.Error:
			switch err.(e.Error).Code {
			case e.AlreadyExists:
				status = hr.LoginAlreadyExists
				break
			case e.InvalidParams:
				internalErr, ok := err.(e.Error).Unwrap().(e.Error)
				if !ok {
					break
				}
				statusMap := map[int]int{
					auth.ErrorWrongLogin:      hr.IncorrectLogin,
					auth.ErrorWrongPassword:   hr.IncorrectPassword,
					auth.ErrorWrongSex:        hr.IncorrectSex,
					auth.ErrorWrongNickName:   hr.IncorrectNickname,
					auth.ErrorWrongFirstName:  hr.IncorrectName,
					auth.ErrorWrongFamilyName: hr.IncorrectFamilyname,
					auth.ErrorWrongBirthDate:  hr.IncorrectBirthDate,
				}
				temp, ok := statusMap[internalErr.Code]
				if !ok {
					break
				}
				status = temp
				break
			default:
				break
			}
			break
		default:
			break
		}
		// a.resp.SetStatus(status)
		resp.SetError(hr.GetError(status))
		return
	}
	token, err := a.usecase.SignIn(userInput.Email, userInput.Password)
	if err != nil {
		status := hr.UnknownError
		switch err.(type) {
		case e.Error:
			status = hr.UnknownError
			break
		default:
			break // в будущем, возможно, появятся другие обработчики
		}
		// a.resp.SetStatus(status)
		resp.SetError(hr.GetError(status))
		return
	}
	a.setCookie(w, sessionTokenCookieName, token, time.Now().Add(cookieLifetime*time.Minute))
}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	resp := a.resp.SetWriter(w).Copy()
	defer resp.Send()
	credentials := struct {
		Login    string
		Password string
	}{}
	err := HttpTools.StructFromBody(*r, &credentials)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}

	uuid, err := a.usecase.SignIn(credentials.Login, credentials.Password)
	if err != nil {
		status := hr.UnknownError
		switch err.(type) {
		case e.Error:
			switch err.(e.Error).Code {
			case e.NotExists:
				status = hr.LoginNotExist
				break
			case e.InvalidParams:
				status = hr.WrongPassword
				break
			default:
			}
			break
		default:
		}
		resp.SetError(hr.GetError(status))
		return
	}
	a.setCookie(w, sessionTokenCookieName, uuid, time.Now().Add(cookieLifetime*time.Minute))
}

func (a *AuthHandler) CheckAuthorization(w http.ResponseWriter, r *http.Request) {
	resp := a.resp.SetWriter(w).Copy()
	defer resp.Send()
	tokenCookie, err := r.Cookie(sessionTokenCookieName)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	err = a.usecase.CheckAuth(tokenCookie.Value)
	if err != nil {
		switch err.(type) {
		case e.Error:
			switch err.(e.Error).Code {
			case e.InvalidParams, e.NoPermission:
				resp.SetError(hr.GetError(hr.BadSession))
			default:
				resp.SetError(hr.GetError(hr.UnknownError))
			}
		default:
			resp.SetError(hr.GetError(hr.UnknownError))
		}
		return
	}
}

func (a *AuthHandler) setCookie(w http.ResponseWriter, name, value string, lifetime time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  lifetime,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
