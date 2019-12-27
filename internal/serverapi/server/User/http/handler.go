package http

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/log"
	auth "2019_2_Next_Level/internal/serverapi/server/Auth"
	hr "2019_2_Next_Level/internal/serverapi/server/HttpError"
	user "2019_2_Next_Level/internal/serverapi/server/User"
	e "2019_2_Next_Level/pkg/Error"
	"2019_2_Next_Level/pkg/HttpTools"
	"fmt"
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
	resp := (&HttpTools.Response{}).SetError(&hr.DefaultResponse)
	handler := UserHandler{usecase: uc, resp: resp}
	return handler
}

func (h *UserHandler) InflateRouter(router *mux.Router) {
	router.HandleFunc("/get", h.GetProfile).Methods("GET")
	router.HandleFunc("/editUserInfo", h.EditUserInfo).Methods("PUT")
	router.HandleFunc("/editPassword", h.EditUserPassword).Methods("POST")
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.InflateNew(w)
	// defer resp.Send()
	login := r.Header.Get("X-Login")
	user, err := h.usecase.GetUser(login)
	if err != nil {
		resp.SetError(hr.GetError(hr.IncorrectLogin)).Send()
		return
	}
	user.Password = ""
	folders, err := h.usecase.GetUserFolders(login)
	if err != nil {
		resp.SetError(hr.GetError(hr.IncorrectLogin)).Send()
		return
	}
	ans := GetUserProfile{
		Status: "ok",
		Answer: GetUserProfileAnswer{Name:user.Name, Sirname:user.Sirname, BirthDate:user.BirthDate,
			Sex:user.Sex, Email:user.Email, Avatar:user.Avatar, Login:user.Login,
			Folders:folders},
	}
	err = HttpTools.BodyFromStruct(w, &ans)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError)).Send()
		return
	}
}

func (h *UserHandler) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.InflateNew(w)
	defer resp.Send()
	login := r.Header.Get("X-Login")

	var newProfile model.User;
	newProfile.InflateFromFormdata(r);
	newProfile.Email = login;

	avaFile, handler, _ := r.FormFile("avatar");
	if avaFile!=nil {
		newFilename, err := h.usecase.EditAvatar(login, avaFile, handler);
		if err != nil {
			log.Log().E(log.GetLogString(login, err))
			resp.SetError(hr.GetError(hr.BadParam));
			return
		}
		log.Log().L("New avatar name ", newFilename)
		newProfile.Avatar = newFilename;
	} else {
		log.Log().E("No avatar")
	}

	err := h.usecase.EditUser(&newProfile);
	if err != nil {
		log.Log().E(log.GetLogString(login, err))
		fmt.Println("FMT, err: ", err)
		//fmt.Println(err)
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
	log.Log().E(log.GetLogString(login, "ok"))

}

func (h *UserHandler) EditUserPassword(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.InflateNew(w)
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
