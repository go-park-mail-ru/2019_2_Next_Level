package http

import (
	"2019_2_Next_Level/internal/model"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	hr "2019_2_Next_Level/internal/serverapi/server/Error/httpError"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	e "2019_2_Next_Level/internal/serverapi/server/Error"
	"2019_2_Next_Level/pkg/HttpTools"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	sessionTokenCookieName = "session-id"
	cookieLifetime         = 10 * 60 // minutes
)

type UserHandler struct {
	usecase user.UserUsecase
	resp    *HttpTools.Response
}

func NewUserHandler(uc user.UserUsecase) UserHandler {
	resp := (&HttpTools.Response{}).SetError(hr.DefaultResponse)
	handler := UserHandler{usecase: uc, resp: resp}
	return handler
}

func (h *UserHandler) InflateRouter(router *mux.Router) {
	router.HandleFunc("/get", h.GetProfile).Methods("GET")
	router.HandleFunc("/editUserInfo", h.EditUserInfo).Methods("POST")
	router.HandleFunc("/editPassword", h.EditUserPassword).Methods("POST")
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	type Answer struct {
		Name      string `json:"firstName"`
		Sirname   string `json:"secondName"`
		BirthDate string `json:"birthDate"`
		Sex       string `json:"sex"`
		Email     string `json:"login"`
		Avatar    string `json:"avatar"`
	}
	resp := h.resp.SetWriter(w).Copy()
	// defer resp.Send()
	login := r.Header.Get("X-Login")
	user, err := h.usecase.GetUser(login)
	if err != nil {
		resp.SetError(hr.GetError(hr.IncorrectLogin)).Send()
		return
	}
	user.Password = ""
	ans := struct{
		Status string `json:"status"`
		Answer Answer `json:"userInfo"`
	}{"ok",
		Answer{user.Name, user.Sirname, user.BirthDate, user.Sex, user.Email, user.Avatar},
	}
	err = HttpTools.BodyFromStruct(w, &ans)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError)).Send()
		return
	}
}

func (h *UserHandler) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	user := model.User{}
	err := HttpTools.StructFromBody(*r, &user)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	login := r.Header.Get("X-Login")
	user.Email = login
	err = h.usecase.EditUser(&user)
	if err != nil {
		status := hr.UnknownError
		switch err.(type) {
		case e.Error:
			switch err.(e.Error).Code {
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
			case e.ProcessError:
				break
			default:
				break
			}
			break
		default:
			break
		}
		resp.SetError(hr.GetError(status))
		return
	}

}

func (h *UserHandler) EditUserPassword(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()

	passStruct := struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{}

	err := HttpTools.StructFromBody(*r, &passStruct)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	err = h.usecase.EditPassword(h.getLogin(*r), passStruct.CurrentPassword, passStruct.NewPassword)
	if err != nil {
		status := hr.UnknownError

		switch err.(type) {
		case e.Error:
			code := err.(e.Error).Code
			statusMap := map[int]int{
				e.WrongPassword: hr.WrongPassword,
				e.Wrong:         hr.SameNewPass,
				e.ProcessError:  hr.UnknownError,
			}
			temp, ok := statusMap[code]
			if !ok {
				break
			}
			status = temp
		default:
			break
		}

		resp.SetError(hr.GetError(status))
	}
}

func (h *UserHandler) getLogin(r http.Request) string {
	return r.Header.Get("X-Login")
}
