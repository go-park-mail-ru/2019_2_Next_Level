package handlers

import (
	"2019_2_Next_Level/internal/model"
	"2019_2_Next_Level/internal/serverapi/log"
	hr "2019_2_Next_Level/internal/serverapi/server/HttpError"
	mailbox "2019_2_Next_Level/internal/serverapi/server/MailBox"
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"2019_2_Next_Level/pkg/HttpTools"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MailHandler struct {
	usecase mailbox.MailBoxUseCase
	resp    *HttpTools.Response
}

// NewMailHandler : sets handlers for specified routes (prefix = "/mail")
func NewMailHandler(usecase mailbox.MailBoxUseCase) *MailHandler{
	handler := MailHandler{usecase: usecase}
	handler.resp = (&HttpTools.Response{}).SetError(&hr.DefaultResponse)
	return &handler
}

func (h *MailHandler) InflateRouter(router *mux.Router) {
	router.HandleFunc("/send", h.SendMail).Methods("POST")
	router.HandleFunc("/getByPage", h.GetMailList).Methods("GET")
	router.HandleFunc("/get", h.GetEmail).Methods("GET")
	router.HandleFunc("/getById", h.GetEmailsById).Methods("POST")
	router.HandleFunc("/getUnreadCount", h.GetUnreadCount).Methods("GET")
	router.HandleFunc("/read", h.MarkMailRead).Methods("POST")
	router.HandleFunc("/unread", h.MarkMailUnRead).Methods("POST")
	router.HandleFunc("/remove", h.DeleteEmail).Methods("POST")
	router.HandleFunc("/addFolder/{name}", h.CreateFolder).Methods("POST")
	router.HandleFunc("/deleteFolder/{name}", h.DeleteFolder).Methods("POST")
	router.HandleFunc("/changeFolder/{id}/{name}", h.ChangeMailFolder).Methods("POST")
	router.HandleFunc("/search/{request}", h.FindMessages).Methods("GET")
}

func (h *MailHandler) FindMessages(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}

	args := mux.Vars(r)
	request, _ := args["request"]
	list, err := h.usecase.FindMessages(login, request)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}
	respList := GetMessagesList{
		Status: "ok",
		Messages: list,
	}
	resp.SetContent(&respList)

}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	mail := models.MailToSend{}
	req := struct{
		Message models.MailToSend `json:"message"`
	}{mail}
	err := HttpTools.StructFromBody(*r, &req)
	mail = req.Message
	if err !=nil {
		log.Log().E(err)
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	email := mail.ToMain()
	email.SetFrom(login)
	err = h.usecase.SendMail(&email)
	if err != nil {
		log.Log().E("Cannot send email")
		resp.SetError(hr.GetError(hr.UnknownError))
	}
}

func (h *MailHandler) GetMailList(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	pageTemp := r.FormValue("page")
	page, err := strconv.ParseInt(pageTemp, 10, 64)
	perPage, err := strconv.ParseInt(r.FormValue("perPage"), 10, 64)
	folder := r.FormValue("folder")

	startLetter := perPage*(page-1)+1
	list, err := h.usecase.GetMailList(login, folder, "", int(startLetter), int(perPage))
	if err != nil {
		log.Log().E("Error after getMailList", err)
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	resp.SetContent(&GetFolderMessagesResponse{
		Status:"ok",
		PagesCount:10,
		Page: int(page),
		Messages: func()[]models.MailToGet{
			localList := make([]models.MailToGet, 0, len(list))
			for _, i := range list {
				localList = append(localList, models.MailToGet{}.FromMain(&i))
			}
			return localList
		}(),
	})
}

func (h *MailHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	count, err := h.usecase.GetUnreadCount(login)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		log.Log().E(err)
		return
	}
	resp.SetContent(&GetMessagesCountResponse{Status:"ok", Count:count})
}

func (h *MailHandler) GetEmail(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	idTemp := r.FormValue("id")
	id, err := strconv.ParseInt(idTemp, 10, 64)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	mails, err := h.usecase.GetMail(login, []models.MailID{models.MailID(id)})
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	mail := mails[0]
	answer := GetMessageResponse{
		Status: "ok",
		Message: models.MailToGet{}.FromMain(&mail),
	}
	resp.SetContent(&answer)
}

func (h *MailHandler) GetEmailsById (w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	req := struct{
		Id []int64 `json:"ids"`
	}{}
	_ = HttpTools.StructFromBody(*r, &req)
	ids := make([]models.MailID, 0, len(req.Id))
	for _, elem := range req.Id {
		ids = append(ids, models.MailID(elem))
	}

	mails, err := h.usecase.GetMail(login, ids)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	resp.SetContent(h.prepareList(mails))
	return

}

func (h *MailHandler) MarkMailRead(w http.ResponseWriter, r *http.Request) {
	h.markMail(w, r, models.MarkMessageRead)
}

func (h *MailHandler) MarkMailUnRead(w http.ResponseWriter, r *http.Request) {
	h.markMail(w, r, models.MarkMessageUnread)
}

func (h *MailHandler) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	h.markMail(w, r, models.MarkMessageDeleted)
}

func (h *MailHandler) markMail(w http.ResponseWriter, r *http.Request, mark int) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login=="" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}

	req := struct {
		Messages []models.MailID
	}{}
	err := HttpTools.StructFromBody(*r, &req)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	err = h.usecase.MarkMail(login, req.Messages, mark)
	if err != nil {
		resp.SetError(hr.GetError(hr.UnknownError))
		return
	}
}

func (h *MailHandler) ChangeMailFolder(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()

	login := h.getLogin(r)
	if login == "" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	args := mux.Vars(r)
	folderName, ok := args["name"]
	if !ok {
		log.Log().E("No such a param: ", "slug")
		return
	}
	mailIdTemp, ok := args["id"]
	if !ok {
		log.Log().E("No such a param: ", "slug")
		return
	}
	mailId, err := strconv.ParseInt(mailIdTemp, 10, 64)
	if err != nil {
		return
	}
	err = h.usecase.ChangeMailFolder(login, folderName, mailId)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	resp.SetContent(&hr.DefaultResponse)

}

func (h *MailHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login == "" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	args := mux.Vars(r)
	folderName, ok := args["name"]
	if !ok {
		log.Log().E("No such a param: ", "slug")
		return
	}

	err := h.usecase.AddFolder(login, folderName)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	resp.SetContent(&hr.DefaultResponse)
}

func (h *MailHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	resp := h.resp.SetWriter(w).Copy()
	defer resp.Send()
	login := h.getLogin(r)
	if login == "" {
		resp.SetError(hr.GetError(hr.BadSession))
		return
	}
	args := mux.Vars(r)
	folderName, ok := args["name"]
	if !ok {
		log.Log().E("No such a param: ", "slug")
		return
	}

	err := h.usecase.DeleteFolder(login, folderName)
	if err != nil {
		resp.SetError(hr.GetError(hr.BadParam))
		return
	}
	resp.SetContent(&hr.DefaultResponse)
}

func (h *MailHandler) getLogin(r *http.Request) string {
	return r.Header.Get("X-Login")
}

func (h *MailHandler) prepareList(list []model.Email) *GetMessagesListResponse {
	answer := GetMessagesListResponse{
		Status: "ok",
		Messages: func()[]models.MailToGet{
			localList := make([]models.MailToGet, 0, len(list))
			for _, elem := range list {
				localList = append(localList, models.MailToGet{}.FromMain(&elem))
			}
			return localList
		}(),
	}
	return &answer
}